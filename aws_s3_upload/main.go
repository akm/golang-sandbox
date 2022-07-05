package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s SOURCE_FILE_PATH DEST_S3_URL\n", os.Args[0])
	}
	if err := process(os.Args[1], os.Args[2]); err != nil {
		panic(err)
	}
}

func process(srcPath, destUrl string) error {
	log.Printf("Uploading %q to %q\n", srcPath, destUrl)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	cfg, err := newConfig(ctx)
	if err != nil {
		return err
	}

	log.Printf("cfg: %#v\n", *cfg)

	srcFile := NewLocalFile(srcPath)
	defer srcFile.TearDown()

	contentType, err := srcFile.ContentType()
	if err != nil {
		return err
	}

	log.Printf("contentType: %q\n", contentType)

	dest, err := url.Parse(destUrl)
	if err != nil {
		return errors.Wrapf(err, "failed to parse dest url: %q", destUrl)
	}

	log.Printf("dest: %#v\n", *dest)

	reader, err := srcFile.Reader()
	if err != nil {
		return err
	}

	{
		cli := s3.NewFromConfig(*cfg)
		log.Printf("cli: %#v\n", *cli)

		_, err := cli.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(dest.Host),
			Key:         aws.String(dest.Path),
			ContentType: aws.String(contentType),
			Body:        reader,
		})
		if err != nil {
			return errors.Wrapf(err, "failed to create multipart upload")
		}
	}

	log.Printf("Uploading done\n")

	return nil
}

func newConfig(ctx context.Context) (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load AWS config")
	}
	if cfg.Region == "" {
		if envRegion := os.Getenv("AWS_REGION"); envRegion != "" {
			cfg.Region = envRegion
		} else {
			cfg.Region = "ap-northeast-1"
		}
	}
	return &cfg, nil
}

type LocalFile struct {
	Path string
	file *os.File
}

func NewLocalFile(path string) *LocalFile {
	return &LocalFile{Path: path}
}

func (m *LocalFile) getFile() (*os.File, error) {
	if m.file == nil {
		f, err := os.Open(m.Path)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to open file: %q", m.Path)
		}
		m.file = f
	} else {
		if _, err := m.file.Seek(0, 0); err != nil {
			return nil, errors.Wrapf(err, "failed to seek file: %q", m.Path)
		}
	}
	return m.file, nil
}

func (m *LocalFile) TearDown() error {
	if m.file == nil {
		return nil
	}
	return m.file.Close()
}

func (m *LocalFile) ContentType() (string, error) {
	f, err := m.getFile()
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read file: %q", m.Path)
	}
	return http.DetectContentType(b), nil
}

func (m *LocalFile) Reader() (io.Reader, error) {
	return m.getFile()
}

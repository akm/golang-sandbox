package main

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
)

func main() {
	if err := process(); err != nil {
		panic(err)
	}
}

func process() error {
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to load AWS config")
	}

	// Regionが空文字列の場合、EKSのDescribeClusterを実行すると以下のようなエラーが発生する
	// operation error EKS: DescribeCluster, failed to resolve service endpoint, an AWS region is required, but was not found
	if cfg.Region == "" {
		if envRegion := os.Getenv("AWS_REGION"); envRegion != "" {
			cfg.Region = envRegion
		} else {
			cfg.Region = "ap-northeast-1"
		}
	}

	s3.NewFromConfig(cfg)
	return nil
}

package pkgerrors

import (
	"fmt"

	stderrors "errors"                // https://golang.org/pkg/errors/#pkg-index
	pkgerrors "github.com/pkg/errors" // https://pkg.go.dev/github.com/pkg/errors
	// "golang.org/x/xerrors" // https://pkg.go.dev/golang.org/x/xerrors
)

func Foo(f func() error) error {
	return Bar(f)
}

func Bar(f func() error) error {
	return Baz(f)
}

func Baz(f func() error) error {
	return f()
}

func StdNewError(s string) error {
	return stderrors.New(fmt.Sprintf("%s with stderrors.New", s))
}

func StdWrapError(err error, s string) error {
	return fmt.Errorf("%s with std errors because of %w", s, err)
}

func PkgNewError(s string) error {
	return pkgerrors.New(fmt.Sprintf("%s with pkgerrors.New", s))
}

func PkgWrapError(err error, s string) error {
	return pkgerrors.Wrapf(err, "%s with pkgerrors.Wrapf", s)
}

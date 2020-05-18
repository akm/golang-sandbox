package main

import (
	"errors"

	"fmt"
	"io"

	"github.com/akm/golang-sandbox/pkg_errors"
)

type CustomError struct {
	Msg string
}

func NewCustomError(s string) *CustomError {
	return &CustomError{s}
}

func (err *CustomError) Error() string {
	return err.Msg
}

func main() {
	errorSlice := []error{
		pkgerrors.Foo(func() error {
			return pkgerrors.StdNewError("stdErr1")
		}),
		pkgerrors.Foo(func() error {
			return pkgerrors.StdWrapError(io.EOF, "stdErr2")
		}),
		pkgerrors.Foo(func() error {
			return pkgerrors.StdWrapError(NewCustomError("CustomError1"), "stdErr3")
		}),

		pkgerrors.Foo(func() error {
			return pkgerrors.PkgNewError("pkgErr1")
		}),
		pkgerrors.Foo(func() error {
			return pkgerrors.PkgWrapError(io.EOF, "pkgErr2")
		}),
		pkgerrors.Foo(func() error {
			return pkgerrors.PkgWrapError(NewCustomError("CustomError2"), "pkgErr3")
		}),
		pkgerrors.Foo(func() error {
			return pkgerrors.PkgWrapError(
				pkgerrors.Foo(func() error {
					return io.EOF
				}), "pkgErr3")
		}),
		pkgerrors.Foo(func() error {
			return pkgerrors.PkgWrapError(
				pkgerrors.Foo(func() error {
					return pkgerrors.PkgNewError("pkgErr4")
				}), "PkgErr4 MEMO")
		}),
	}

	for idx, err := range errorSlice {
		fmt.Printf("%d: v: %v\n", idx, err)

		if errors.Is(err, io.EOF) {
			fmt.Printf("%d: is EOF\n", idx)
		} else {
			fmt.Printf("%d: is NOT EOF\n", idx)
		}

		var cerr *CustomError
		if errors.As(err, &cerr) {
			fmt.Printf("%d: is CustomError: %v\n", idx, cerr)
		} else {
			fmt.Printf("%d: is NOT CustomError\n", idx)
		}

		{
			cause := errors.Unwrap(err)
			fmt.Printf("%d: errors.Unwrap: %T %v\n", idx, cause, cause)
		}

		fmt.Printf("%d: +v: %+v\n", idx, err)
	}
}

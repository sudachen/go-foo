package fu

import (
	"bytes"
	"fmt"
	"golang.org/x/xerrors"
	"strings"
)

func Etrace(err error) error {
	if _, ok := err.(xerrors.Formatter); ok {
		return err
	}
	return error_{err, xerrors.Caller(1)}
}

func Errorf(f string, a ...interface{}) error {
	return error_{fmt.Errorf(f, a...), xerrors.Caller(1)}
}

type error_ struct {
	error
	frame xerrors.Frame
}

func (e error_) FormatError(p xerrors.Printer) error {
	p.Print(e.error.Error() + " at ")
	e.frame.Format(p)
	return nil
}

func Panic(err error) interface{} {
	return panic_{err}
}

type panic_ struct{ err error }

func stringifyError(err error) (string, error) {
	ep := &errorPrinter{details: true}
	if f, ok := err.(xerrors.Formatter); ok {
		err = f.FormatError(ep)
	} else {
		ep.Print(err.Error())
		err = nil
	}
	return strings.Join(strings.Fields(ep.String()), " "), err
}

func (x panic_) stringify(indepth bool) string {
	s, e := stringifyError(x.err)
	ns := []string{s}
	for e != nil && indepth {
		s, e = stringifyError(e)
		ns = append(ns, s)
	}
	return strings.Join(ns, "\n")
}

func (x panic_) Error() string {
	return x.stringify(false)
}

func (x panic_) String() string {
	return x.stringify(true)
}

func (x panic_) Unwrap() error {
	if w, ok := x.err.(xerrors.Wrapper); ok {
		return w.Unwrap()
	}
	return x.err
}

type errorPrinter struct {
	bytes.Buffer
	details bool
}

func (ep *errorPrinter) Print(args ...interface{}) {
	ep.Buffer.WriteString(fmt.Sprint(args...))
}

func (ep *errorPrinter) Printf(format string, args ...interface{}) {
	ep.Buffer.WriteString(fmt.Sprintf(format, args...))
}

func (ep errorPrinter) Detail() bool {
	return ep.details
}

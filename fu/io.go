package fu

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

type Input interface {
	Open() (io.ReadCloser, error)
}

type Whole interface {
	io.Writer
	Commit() error
	End()
}

type Output interface {
	Create() (Whole, error)
}

type Inout interface {
	Input
	Output
}

type Sizeable interface {
	Size() int64
}

func FileSize(rd io.Reader) int64 {
	if i, ok := rd.(Sizeable); ok {
		return i.Size()
	}
	if f, ok := rd.(*os.File); ok {
		st, err := f.Stat()
		if err == nil {
			return st.Size()
		}
	}
	return 0
}

type Resettable interface {
	Reset() error
}

func ResetFile(rd io.Reader) error {
	if i, ok := rd.(Resettable); ok {
		return i.Reset()
	}
	return nil
}

type regularXf struct {
	*os.File
}

func (tf regularXf) Reset() error {
	_, err := tf.File.Seek(0, 0)
	return err
}

func (tf regularXf) Size() int64 {
	st, _ := tf.File.Stat()
	return st.Size()
}

type temporalXf struct {
	regularXf
	deleted bool
}

func (tf temporalXf) Close() error {
	_ = tf.File.Close()
	if !tf.deleted {
		_ = os.Remove(tf.File.Name())
		tf.deleted = true
	}
	return nil
}

func Tempfile(pattern string) (_ io.ReadWriteCloser, err error) {
	var f *os.File
	if f, err = ioutil.TempFile("", pattern); err != nil {
		return
	}
	return &temporalXf{regularXf{f}, false}, nil
}

type StringIO string

func (s StringIO) Open() (io.ReadCloser, error) {
	return Reader_{bytes.NewBufferString(string(s)), nil},
		nil
}

type CloserChain [2]io.Closer

func (c CloserChain) Close() error {
	for _, x := range c {
		if x != nil {
			_ = x.Close()
		}
	}
	return nil
}

type File string

func (f File) Open() (io.ReadCloser, error) {
	return os.Open(string(f))
}

type FFile_ struct{ *os.File }

func (ff FFile_) Fail() {
	fname := ff.File.Name()
	_ = ff.File.Truncate(0)
	_ = ff.File.Close()
	_ = os.Remove(fname)
}

func (f File) Create() (Whole, error) {
	x, err := os.Create(string(f))
	if err != nil {
		return nil, err
	}
	return &Whole_{FFile_{x}}, nil
}

type Writer_ struct {
	io.Writer
	close []func(bool) error
}

func Writer(wr io.Writer, close ...func(bool) error) Writer_ {
	return Writer_{wr, close}
}

func (w Writer_) Create() (Whole, error) {
	return &w, nil
}

func (w *Writer_) End() {
	for _, f := range w.close {
		_ = f(false)
	}
	return
}

func (w *Writer_) Commit() (err error) {
	for _, f := range w.close {
		if e := f(true); e != nil {
			err = e
		}
	}
	w.close = nil
	return
}

type Reader_ struct {
	io.Reader
	close []func() error
}

func (r Reader_) Close() (err error) {
	for _, f := range r.close {
		if e := f(); e != nil {
			err = e
		}
	}
	return
}

func (r Reader_) Open() (io.ReadCloser, error) {
	return r, nil
}

func Reader(rd io.Reader, close ...func() error) Reader_ {
	return Reader_{rd, close}
}

type Whole_ struct{ io.Writer }
type Fallible interface {
	Fail()
}

func (t *Whole_) End() {
	if t.Writer != nil {
		if f, ok := t.Writer.(Fallible); ok {
			f.Fail()
		} else if c, ok := t.Writer.(io.Closer); ok {
			_ = c.Close()
		}
		t.Writer = nil
	}
}

func (t *Whole_) Commit() (err error) {
	if c, ok := t.Writer.(io.Closer); ok {
		err = c.Close()
	}
	t.Writer = nil
	return
}

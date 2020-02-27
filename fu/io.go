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

type Sizeable interface {
	Size() int64
}

func FileSize(rd io.Reader) int64 {
	if i, ok := rd.(Sizeable); ok {
		return i.Size()
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

type wrapcloseXf struct {
	io.Reader
	close func() error
}

func (w wrapcloseXf) Close() error {
	if w.close != nil {
		err := w.close()
		w.close = nil
		return err
	}
	return nil
}

func WrapClose(rd io.Reader, close func() error) io.ReadCloser {
	return &wrapcloseXf{rd, close}
}

type StringIO string

func (s StringIO) Open() (io.ReadCloser, error) {
	return wrapcloseXf{bytes.NewBufferString(string(s)), nil},
		nil
}

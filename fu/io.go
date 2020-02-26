package fu

import (
	"bytes"
	"github.com/sudachen/go-fp/internal"
	"golang.org/x/xerrors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Close io.ReadCloser

type Input interface {
	Open() (io.ReadCloser,error)
}

type Sizeble interface {
	Size() int64
}

func FileSize(rd io.Reader) int64 {
	if i, ok := rd.(Sizeble); ok {
		return i.Size()
	}
	return 0
}

type Resetable interface {
	Reset() error
}

func ResetFile(rd io.Reader) error {
	if i, ok := rd.(Resetable); ok {
		return i.Reset()
	}
	return nil
}

type regularXf struct {
	*os.File
}

func (tf regularXf) Reset() error {
	_, err := tf.File.Seek(0,0)
	return err
}

func (tf regularXf) Size() int64 {
	st, _ := tf.File.Stat()
	return st.Size()
}


type temporalXf struct {
	regularXf
	path string
}

func (tf temporalXf) Close() error {
	_ = tf.File.Close()
	if tf.path != "" {
		_ = os.Remove(tf.path)
		tf.path = ""
	}
	return nil
}

type wrapcloseXf struct {
	io.Reader
	close func()error
}

func (w wrapcloseXf) Close() error {
	if w.close != nil { return w.close() }
	return nil
}

func WrapClose(rd io.Reader, close func()error) io.ReadCloser {
	return &wrapcloseXf{rd, close}
}

var tempfileRand = internal.NaiveRandom{}

func Tempfile(pattern string) (_ io.ReadWriteCloser, err error) {
	dir := os.TempDir()
	var prefix, suffix string
	if pos := strings.LastIndex(pattern, "*"); pos != -1 {
		prefix, suffix = pattern[:pos], pattern[pos+1:]
	} else {
		prefix = pattern
	}

	var f *os.File
	var name string
	nconflict := 0
	for i := 0; i < 10000; i++ {
		r := strconv.Itoa(int(1e9 + tempfileRand.Next()%1e9))[1:]
		name = filepath.Join(dir, prefix + r + suffix)
		f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if os.IsExist(err) {
			if nconflict++; nconflict > 10 {
				tempfileRand.Reseed()
			}
			continue
		}
		break
	}
	if err != nil {
		return nil, xerrors.Errorf("failed to create unique temporal file")
	}
	return &temporalXf{regularXf{f}, name}, nil
}

type StringIO string
func (s StringIO) Open() (io.ReadCloser, error) {
	return wrapcloseXf{
		bytes.NewBufferString(string(s)),
		func()error{ return nil }},
		nil
}

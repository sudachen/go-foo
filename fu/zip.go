package fu

import (
	"archive/zip"
	"golang.org/x/xerrors"
	"io"
	"os"
)

type ZipFile_ struct {
	Arch     interface{}
	FileName string
}

func ZipFile(fileName string, arch interface{}) ZipFile_ {
	return ZipFile_{arch, fileName}
}

func (q ZipFile_) Open() (f io.ReadCloser, err error) {
	var xf io.ReadCloser
	if e, ok := q.Arch.(Input); ok {
		xf, err = e.Open()
	} else {
		xf, err = os.Open(q.Arch.(string))
	}
	if err != nil {
		return
	}
	defer func() {
		if xf != nil {
			_ = xf.Close()
		}
	}()
	var r *zip.Reader
	if r, err = zip.NewReader(xf.(io.ReaderAt), FileSize(xf)); err != nil {
		return
	}
	for _, n := range r.File {
		if n.Name == q.FileName {
			zf, err := n.Open()
			if err != nil {
				return nil, err
			}
			xxf := xf
			xf = nil
			return Reader(zf, func() error {
				_ = zf.Close()
				return xxf.Close()
			}), nil
		}
	}
	return nil, xerrors.Errorf("zip archive does not contain file " + q.FileName)
}

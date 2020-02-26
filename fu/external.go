package fu

import (
	"golang.org/x/xerrors"
	"io"
	"net/http"
	"os"
)

type Cached string
func (c Cached) Remove() (err error) {
	s := CacheFile(string(c))
	_, err = os.Stat(s)
	if err == nil {
		return os.Remove(s)
	}
	return nil
}

type External_ struct {
	url string
	cache string
}

func External(url string, opts ...interface{}) External_ {
	return External_{url, StrOption(Cached(""),opts)}
}

func (e External_) Open() (io.ReadCloser, error) {
	return CachedDownload(e.url,e.cache)
}

func CachedDownload(url string, cached string) (_ io.ReadCloser, err error) {
	var f io.ReadWriteCloser
	if cached != "" {
		cached = CacheFile(cached)
		if _, err = os.Stat(cached); err == nil {
			if f, err = os.Open(cached); err != nil {
				return nil, xerrors.Errorf("filed to open cached file: %w", err)
			}
			return
		}
		if f, err = os.Create(cached); err != nil {
			return
		}
	} else {
		if f, err = Tempfile("external-noncached"); err != nil {
			return nil, xerrors.Errorf("could not create temporal file: %w", err)
		}
	}
	err = download(url, f.(io.Writer))
	if err != nil {
		_ = f.Close()
		return nil, xerrors.Errorf("download error: %w", err)
	}
	_ = ResetFile(f)
	return f, nil
}

func download(url string, writer io.Writer) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(writer, resp.Body)
	return err
}


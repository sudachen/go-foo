package fu

import (
	"os"
	"path"
	"path/filepath"
)

const cacheGoFp = ".cache/gofper"

var FullCacheDir string

func init() {
	if u, ok := os.LookupEnv("HOME"); ok {
		FullCacheDir, _ = filepath.Abs(filepath.Join(u, cacheGoFp))
	} else {
		FullCacheDir, _ = filepath.Abs(cacheGoFp)
	}
}

func CacheDir(d string) string {
	r := path.Join(FullCacheDir, d)
	_ = os.MkdirAll(r, 0777)
	return r
}

func CacheFile(f string) string {
	r := path.Join(FullCacheDir, f)
	_ = os.MkdirAll(path.Dir(r), 0777)
	return r
}


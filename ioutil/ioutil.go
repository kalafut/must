package ioutil

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/kalafut/must"
)

func NopCloser(r io.Reader) io.ReadCloser {
	return ioutil.NopCloser(r)
}

func ReadAll(r io.Reader) []byte {
	data, err := ioutil.ReadAll(r)
	must.PanicIfErr(err)

	return data
}

func ReadDir(dirname string) []os.FileInfo {
	files, err := ioutil.ReadDir(dirname)
	must.PanicIfErr(err)

	return files
}

func ReadFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	must.PanicIfErr(err)

	return data
}

func TempDir(dir, pattern string) string {
	name, err := ioutil.TempDir(dir, pattern)
	must.PanicIfErr(err)

	return name

}

func TempFile(dir, pattern string) *os.File {
	f, err := ioutil.TempFile(dir, pattern)
	must.PanicIfErr(err)

	return f
}

func WriteFile(filename string, data []byte, perm os.FileMode) {
	err := ioutil.WriteFile(filename, data, perm)
	must.PanicIfErr(err)
}

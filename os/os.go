package os

import (
	"os"

	"github.com/kalafut/must"
)

type File os.File
type FileInfo os.FileInfo

func Open(name string) *File {
	ret, err := os.Open(name)
	must.PanicErr(err)

	return (*File)(ret)
}

func (f *File) Stat() FileInfo {
	ret, err := (*os.File)(f).Stat()
	must.PanicErr(err)

	return ret
}

func (f *File) Readdirnames(n int) []string {
	ret, err := (*os.File)(f).Readdirnames(n)
	must.PanicErr(err)

	return ret
}

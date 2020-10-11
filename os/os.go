package os

import (
	"os"
	"time"

	"github.com/kalafut/must"
)

type File os.File
type FileInfo os.FileInfo
type FileMode os.FileMode

func Chdir(dir string) {
	must.PanicErr(os.Chdir(dir))
}

func Chmod(name string, mode FileMode) {
	must.PanicErr(os.Chmod(name, os.FileMode(mode)))
}

func Chown(name string, uid, gid int) error {
	must.PanicErr(os.Chown(name, uid, gid))
}

func Chtimes(name string, atime time.Time, mtime time.Time) error {
	must.PanicErr(os.Chtimes(name, atime, mtime))
}

func Clearenv() {
	os.Clearenv()
}

func Environ() []string {
	return os.Environ()
}

func Executable() string {
	ret, err := os.Executable()
	must.PanicErr(err)
	return ret
}

func Exit(code int) {
	os.Exit(code)
}

func Expand(s string, mapping func(string) string) string {
	return os.Expand(s, mapping)
}

func ExpandEnv(s string) string {
	return os.ExpandEnv(s)
}

func Getegid() int {
	return os.Getegid()
}

func Getenv(key string) string {
	return os.Getenv(key)
}

func Geteuid() int {
	return os.Geteuid()
}

func Getgid() int {
	return os.Getgid()
}

func Getgroups() []int {
	ret, err := os.Getgroups()
	must.PanicErr(err)
	return ret
}

func Getpagesize() int {
	return os.Getpagesize()
}

func Getpid() int {
	return os.Getpid()
}

func Getppid() int {
	return os.Getppid()
}

func Getuid() int {
	return os.Getuid()
}

func Getwd() (dir string) {
	ret, err := os.Getwd()
	must.PanicErr(err)
	return ret
}

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

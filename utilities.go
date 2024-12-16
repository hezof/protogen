package protogen

import (
	"fmt"
	"os"
	"path/filepath"
)

func Exists(path string) bool {
	sta, err := os.Stat(path)
	return sta != nil || os.IsExist(err)
}

func Create(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	if !Exists(dir) {
		os.MkdirAll(dir, 0755)
	}
	return os.Create(path)
}

func Env(key string, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return val
	}
	return val
}

func PrintError(format string, args ...interface{}) {
	//fmt.Fprintf(os.Stdout, "%s [E] - %s\n", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, args...))
	fmt.Fprintln(os.Stdout, "[protogen]", fmt.Sprintf(format, args...))
}

func PrintInfo(format string, args ...interface{}) {
	//fmt.Fprintf(os.Stdout, "%s [I] - %s\n", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, args...))
	fmt.Fprintln(os.Stdout, "[protogen]", fmt.Sprintf(format, args...))
}

func GoGet(goproxy, goprivate string, module, version string, path string) error {

}

func GoInstall(goproxy, goprivate string, module, version string, path string) error {

}

func HttpGet(furl string, path string) error {

}

func GoGetFile(goproxy, goprivate string, module, version string, srcPath, dstPath string) error {

}

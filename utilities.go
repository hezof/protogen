package protogen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func PrintExit(format string, args ...any) {
	fmt.Fprintln(os.Stdout, `protogen [E]`, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func PrintInfo(format string, args ...any) {
	fmt.Fprintln(os.Stdout, `protogen [I]`, fmt.Sprintf(format, args...))
}

func PrintWarn(format string, args ...any) {
	fmt.Fprintln(os.Stdout, `protogen [W]`, fmt.Sprintf(format, args...))
}

func Env(key string, def string) string {
	val := os.Getenv(key)
	if val != `` {
		return val
	}
	return def
}

func Exists(path string) bool {
	sta, err := os.Stat(path)
	return sta != nil || os.IsExist(err)
}

func Lookup(cmd string) string {
	path, err := exec.LookPath(cmd)
	if err != nil {
		path = filepath.Join(filepath.Dir(os.Args[0]), cmd)
	}
	path, _ = filepath.Abs(path)
	return path
}

func RealPath(tmp string, module string, steps ...string) string {

	/*
		#https://golang.google.cn/ref/mod:
		To avoid ambiguity when serving from case-insensitive file systems, the $module and $version elements
		are case-encoded by replacing every uppercase letter with an exclamation mark followed by the corresponding
		lower-case letter. This allows modules example.com/M and example.com/m to both be stored on disk,
		since the former is encoded as example.com/!m.
	*/

	{
		const diff = 'a' - 'A'
		sb := new(strings.Builder)
		for _, c := range module {
			if c >= 'A' && c <= 'Z' {
				sb.WriteByte('!')
				sb.WriteRune(c + diff)
			} else {
				sb.WriteRune(c)
			}
		}
		module = sb.String()
	}

	parent := filepath.Join(tmp, filepath.Dir(module))
	if !Exists(parent) {
		return ""
	}
	prefix := filepath.Base(module)
	if at := strings.IndexByte(prefix, '@'); at > 0 {
		prefix = prefix[:at+1]
	}

	list, err := os.ReadDir(parent)
	if err != nil {
		return ""
	}
	for _, item := range list {
		if strings.HasPrefix(item.Name(), prefix) {
			root := filepath.Join(parent, item.Name())
			for _, step := range steps {
				root = filepath.Join(root, step)
				if !Exists(root) {
					return ""
				}
			}
			return root
		}
	}
	return ""
}

func EnvironExclude(excludes ...string) (result []string) {
__NEXT__:
	for _, env := range os.Environ() {
		for _, exclude := range excludes {
			if strings.HasPrefix(env, exclude) {
				continue __NEXT__
			}
		}
		result = append(result, env)
	}
	return
}

func keys(m map[string]any) []string {
	s := make([]string, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	return s
}

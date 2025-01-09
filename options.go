package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// CustomOptions 用户选项
type CustomOptions struct {
	Help      bool   // 打印帮助
	Debug     bool   // 打印调试
	Update    bool   // 更新插件
	Config    string // 配置变量, 例如: "VERSION=0.5.1;GOPROXY=https://goproxy.cn;GOPRIVATE=*.net,*.cn"
	GoOut     string // GO输出目录
	ProtoBase string // PB基准目录
	ProtoPath string // PB查找目录,多值逗号分隔
	GrpcV2    bool   // 生成GRPCv2代码[require_unimplemented_servers=true]
}

// SystemOptions 系统选项
type SystemOptions struct {
	HomeDir    string // .protogen目录
	TempDir    string // .protogen/tmp目录
	IncludeDir string // .protogen/include目录
	GoModFile  string // .protogen/go.mod文件
	GoSumFile  string // .protogen/go.sum文件

	GO111MODULE string
	GOSUMDB     string
	GOEXE       string
	GOBIN       string
	GOMODCACHE  string
	GOCACHE     string
	GOTMPDIR    string
}

func initCustomOptions(ops *Context) {
	ops.flagset.BoolVar(&ops.Help, `help`, false, `打印帮助`)
	ops.flagset.BoolVar(&ops.Debug, `debug`, false, `打印调试`)
	ops.flagset.BoolVar(&ops.Update, `update`, false, `更新插件`)
	ops.flagset.StringVar(&ops.Config, `config`, ``, fmt.Sprintf(`配置变量.默认"VERSION=%v;GOPROXY=%v;GOPRIVATE=%v;MAVEN_CENTRAL=%v"`, Version, `https://goproxy.cn`, `*.net,*.cn`, `https://maven.aliyun.com/repository/central`))
	ops.flagset.StringVar(&ops.GoOut, `go_out`, work(), `GO输出目录,默认当前目录`)
	ops.flagset.StringVar(&ops.ProtoBase, `proto_base`, work(), `PB基准目录,默认当前目录`)
	ops.flagset.StringVar(&ops.ProtoPath, `proto_path`, ``, `PB查找目录[逗号分隔]`)
	ops.flagset.BoolVar(&ops.GrpcV2, `grpc_v2`, false, `生成GRPC代码[require_unimplemented_servers=true]`)
}

func initSystemOptions(ops *Context) {
	ops.HomeDir = home()
	ops.TempDir = filepath.Join(ops.HomeDir, `tmp`)
	ops.IncludeDir = filepath.Join(ops.HomeDir, `include`)
	ops.GoModFile = filepath.Join(ops.HomeDir, `go.mod`)
	ops.GoSumFile = filepath.Join(ops.HomeDir, `go.sum`)

	ops.GO111MODULE = `on`
	ops.GOSUMDB = `off`
	ops.GOEXE = goexe()
	ops.GOBIN = ops.TempDir
	ops.GOMODCACHE = ops.TempDir
	ops.GOCACHE = ops.TempDir
	ops.GOTMPDIR = ops.TempDir

}

type Config struct {
	VERSION       string // protogen版本, 默认: Version
	GOPROXY       string // go代理仓库, 默认: https://goproxy.cn
	GOPRIVATE     string // go私有代理, 默认: *.net,*.cn
	MAVEN_CENTRAL string // maven中央仓库, 默认: https://maven.aliyun.com/repository/central
}

func parseConfig(s string) *Config {
	c := new(Config)
	// 默认值
	c.VERSION = Version
	c.GOPROXY = `https://goproxy.cn`
	c.GOPRIVATE = `*.net,*.cn`
	c.MAVEN_CENTRAL = `https://maven.aliyun.com/repository/central`
	// 参数值
	for _, env := range strings.Split(s, ";") {
		kvs := strings.SplitN(strings.TrimSpace(env), "=", 3)
		if len(kvs) > 1 {
			k := strings.TrimSpace(kvs[0])
			v := strings.TrimSpace(kvs[1])
			switch k {
			case `VERSION`:
				c.VERSION = v
			case `GOPROXY`:
				c.GOPROXY = v
			case `GOPRIVATE`:
				c.GOPRIVATE = v
			case `MAVEN_CENTRAL`:
				c.MAVEN_CENTRAL = v
			}
		}
	}
	return c
}

func work() string {
	cwd, _ := os.Getwd()
	if cwd == "" {
		cwd = "./"
	}
	cwd, _ = filepath.Abs(cwd)
	return cwd
}

func home() string {
	loc, err := exec.LookPath(program)
	if err != nil {
		loc = program
	}
	abs, err := filepath.Abs(loc)
	if err != nil {
		abs = loc
	}
	return filepath.Join(filepath.Dir(abs), `.protogen`)
}

func goexe() string {
	switch runtime.GOOS {
	case `windows`:
		return `.exe`
	default:
		return ``
	}
}

func root(path string) string {
	if path == `` {
		if cwd, _ := os.Getwd(); cwd != `` {
			path = cwd
		} else {
			path = `./`
		}
	}
	ret, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return ret
}

// 获取程序名称,避免有人窜改os.args
var program = os.Args[0]

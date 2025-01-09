package protogen

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// CustomOptions 用户选项
type CustomOptions struct {
	Version      bool
	Update       bool
	Debug        bool
	GoOut        string // GO输出目录
	ProtoPath    string // PD查找路径,多值逗号分隔
	GrpcV2       bool
	GoProxy      string
	GoPrivate    string
	MavenCentral string
}

// SystemOptions 系统选项
type SystemOptions struct {
	WorkDir    string // 当前工作目录
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
	ops.flagset.BoolVar(&ops.Version, `version`, false, `打印版本`)
	ops.flagset.BoolVar(&ops.Update, `update`, false, `更新插件`)
	ops.flagset.BoolVar(&ops.Debug, `debug`, false, `打印调试`)
	ops.flagset.StringVar(&ops.GoOut, `go_out`, work(), `Go输出路径`)
	ops.flagset.StringVar(&ops.ProtoPath, `proto_path`, work(), `PB查找路径[逗号分隔]`)
	ops.flagset.BoolVar(&ops.GrpcV2, `grpc_v2`, false, `生成GRPC代码[require_unimplemented_servers=true]`)
	ops.flagset.StringVar(&ops.GoProxy, `goproxy`, Env(`GOPROXY`, `https://goproxy.cn`), `GOPROXY代理仓库`)
	ops.flagset.StringVar(&ops.GoPrivate, `goprivate`, Env(`GOPRIVATE`, `*.net,*.cn`), `GOPRIVATE私有模块`)
	ops.flagset.StringVar(&ops.MavenCentral, `maven_central`, Env(`MAVEN_CENTRAL`, `https://maven.aliyun.com/repository/central`), `MAVEN中央仓库`)
}

func initSystemOptions(ops *Context) {
	ops.WorkDir = work()
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

package protogen

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// CustomOptions 用户选项
type CustomOptions struct {
	Help      bool
	Version   bool
	Debug     bool
	Update    bool
	RootPath  string // PD根路径
	ProtoPath string // PD查找路径,多值逗号分隔
	All       bool
	Docs      bool
	Json      bool
	Bson      bool
	Sqlx      bool
	Grpc      bool
	GrpcV2    bool
	Http      bool

	GO        string
	GOPROXY   string
	GOPRIVATE string
	CENTRAL   string
	PROFILE   string
}

// SystemOptions 系统选项
type SystemOptions struct {
	HOME        string
	TEMP        string
	GO111MODULE string
	GOSUMDB     string
	GOEXE       string
	GOBIN       string
	GOMODCACHE  string
	GOCACHE     string
	GOTMPDIR    string

	GoModFile string
	GoSumFile string
}

func initCustomOptions(ops *Context) {
	ops.flagset.BoolVar(&ops.Help, `h`, false, `打印帮助`)
	ops.flagset.BoolVar(&ops.Help, `help`, false, `打印帮助`)
	ops.flagset.BoolVar(&ops.Version, `version`, false, `打印版本`)
	ops.flagset.BoolVar(&ops.Debug, `debug`, false, `打印调试`)
	ops.flagset.BoolVar(&ops.Update, `update`, false, `更新插件`)
	ops.flagset.StringVar(&ops.RootPath, `root_path`, ``, `PB根路径`)
	ops.flagset.StringVar(&ops.ProtoPath, `proto_path`, ``, `PB查找路径[逗号分隔]`)
	ops.flagset.BoolVar(&ops.All, `all`, false, `执行所有插件`)
	ops.flagset.BoolVar(&ops.Docs, `docs`, false, `生成文档片段[openapi]`)
	ops.flagset.BoolVar(&ops.Json, `json`, false, `生成JSON代码`)
	//ops.flagset.BoolVar(&ops.Json, `bson`, false, `生成BSON代码`)
	//ops.flagset.BoolVar(&ops.Sqlx, `sqlx`, false, `生成SQLX代码`)
	ops.flagset.BoolVar(&ops.Grpc, `grpc`, false, `生成GRPC代码`)
	ops.flagset.BoolVar(&ops.GrpcV2, `grpc_v2`, false, `生成GRPC代码[require_unimplemented_servers=true]`)
	ops.flagset.BoolVar(&ops.Http, `http`, false, `生成HTTP代码[restful,websocket,sse]`)

	ops.flagset.StringVar(&ops.GO, `go`, Env(`GO`, `go`), `GO命令路径`)
	ops.flagset.StringVar(&ops.GOPROXY, `goproxy`, Env(`GOPROXY`, `https://goproxy.cn`), `GOPROXY代理仓库`)
	ops.flagset.StringVar(&ops.GOPRIVATE, `goprivate`, Env(`GOPRIVATE`, `*.net,*.cn`), `GOPRIVATE私有模块`)
	ops.flagset.StringVar(&ops.CENTRAL, `central`, Env(`MAVEN_CENTRAL`, `https://maven.aliyun.com/repository/central`), `MAVEN中央仓库`)
	ops.flagset.StringVar(&ops.PROFILE, `profile`, Env(`PROXY_PROFILE`, Profile), `PROXY配置模块`)
}

func initSystemOptions(ops *Context) {
	ops.HOME = home()
	ops.TEMP = filepath.Join(ops.HOME, `tmp`)
	ops.GO111MODULE = `on`
	ops.GOSUMDB = `off`
	ops.GOEXE = goexe()
	ops.GOBIN = ops.TEMP
	ops.GOMODCACHE = ops.TEMP
	ops.GOCACHE = ops.TEMP
	ops.GOTMPDIR = ops.TEMP

	ops.GoModFile = filepath.Join(ops.HOME, `go.mod`)
	ops.GoSumFile = filepath.Join(ops.HOME, `go.sum`)
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

var program = os.Args[0] // 批复程序名称,避免有人窜改os.args

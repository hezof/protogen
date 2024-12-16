package protogen

import (
	"flag"
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
	ProtoPath string // 多值逗号分隔
	GoOut     string
	All       bool
	Docs      bool
	Json      bool
	Bson      bool
	Sqlx      bool
	Grpc      bool
	GrpcV2    bool
	Http      bool
	Args      []string

	GO        string
	GOPROXY   string
	GOPRIVATE string
	MAVEN     string
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

func getCustomOptions(args []string) CustomOptions {
	var ops CustomOptions
	flagset.BoolVar(&ops.Help, `h`, false, `打印帮助`)
	flagset.BoolVar(&ops.Help, `help`, false, `打印帮助`)
	flagset.BoolVar(&ops.Version, `version`, false, `打印版本`)
	flagset.BoolVar(&ops.Debug, `debug`, false, `打印调试`)
	flagset.BoolVar(&ops.Update, `update`, false, `更新插件`)
	flagset.StringVar(&ops.ProtoPath, `proto_path`, ``, `PB查找路径(逗号分隔)`)
	flagset.StringVar(&ops.GoOut, `go_out`, ``, `GO输出路径`)
	flagset.BoolVar(&ops.All, `all`, false, `执行所有插件`)
	flagset.BoolVar(&ops.Docs, `docs`, false, `生成文档片段(openapi)`)
	flagset.BoolVar(&ops.Json, `json`, false, `生成JSON代码`)
	flagset.BoolVar(&ops.Sqlx, `sqlx`, false, `生成SQLX代码`)
	flagset.BoolVar(&ops.Grpc, `grpc`, false, `生成GRPC代码`)
	flagset.BoolVar(&ops.GrpcV2, `grpc`, false, `生成GRPC代码(require_unimplemented_servers=true)`)
	flagset.BoolVar(&ops.Http, `http`, false, `生成HTTP代码(restful,websocket,server-send-events)`)

	flagset.StringVar(&ops.GO, `go`, Env(`GO`, `go`), `GO命令路径`)
	flagset.StringVar(&ops.GOPROXY, `goproxy`, Env(`GOPROXY`, `https://goproxy.cn`), `$GOPROXY代理仓库`)
	flagset.StringVar(&ops.GOPRIVATE, `goprivate`, Env(`GOPRIVATE`, `*.net,*.cn`), `$GOPRIVATE私有仓库`)
	flagset.StringVar(&ops.MAVEN, `maven`, Env(`MAVEN`, `https://maven.aliyun.com/repository/central`), `$MAVEN代理仓库`)

	err := flagset.Parse(args)
	if err != nil {
		PrintError("parse argument error: %v", err)
		os.Exit(1)
	}

	ops.Args = flagset.Args()

	return ops
}

func getSystemOptions() SystemOptions {
	var ops SystemOptions
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
	return ops
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

var program = os.Args[0]                                    // 批复程序名称,避免有人窜改os.Args
var flagset = flag.NewFlagSet("protogen", flag.ExitOnError) // 命令行解析工具

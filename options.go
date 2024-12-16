package protogen

import (
	"os"
	"os/exec"
	"path/filepath"
)

type Options struct {
	Help      bool   // 打印帮助信息
	Version   bool   // 打印版本, 包括各类组件的版本
	Debug     bool   // 打印调试信息
	GoCmd     string // GO程序位置
	GoExe     string // GO程序后缀
	GoProxy   string // GO代理服务,默认https://goproxy.cn
	GoPrivate string // GO私有仓库.
	Update    bool   // 升级protogen版本
	ProtoPath string // PB查找目录,多值用逗号分隔
	All       bool   // 执行所有插件
	Base      bool   // protoc-gen-go
	Json      bool   // protoc-gen-go-json
	Sqlx      bool   // protc-gen-go-sqlx
	Bson      bool   // protc-gen-go-bson
	Grpc      bool   // protoc-gen-go-grpc(require_unimplemented_servers=false)
	GrpcV2    bool   // protoc-gen-go-grpc(require_unimplemented_servers=true)
	Http      bool   // protoc-gen-go-http
	Docs      bool   // protoc-gen-go-docs
}

func options(args []string) (ops *Options, code int) {

	ops = new(Options)

	flags.BoolVar(&ops.Help, "h", false, "打印帮助信息")
	flags.BoolVar(&ops.Help, "help", false, "打印帮助信息")
	flags.BoolVar(&ops.Version, "v", false, "打印版本信息")
	flags.BoolVar(&ops.Version, "version", false, "打印版本信息")
	flags.BoolVar(&ops.Debug, "d", false, "打印调试信息")
	flags.BoolVar(&ops.Debug, "debug", false, "打印调试信息")

	flags.StringVar(&ops.GoProxy, "goproxy", Env("GOPROXY", "https://goproxy.cn"), "GO代理服务,默认$GOPROXY(https://goproxy.cn)")
	flags.StringVar(&ops.GoPrivate, "goprivate", Env("GOPRIVATE", ""), "GO私有仓库,默认$GOPRIVATE")
	flags.BoolVar(&ops.Update, "update", false, "更新依赖插件")

	flags.StringVar(&ops.ProtoPath, "proto_path", "", "PB查找目录(多值用逗号分隔)")

	flags.BoolVar(&ops.All, "all", false, "执行全部插件")
	flags.BoolVar(&ops.Base, "cwd", false, "执行PB插件[protoc-gen-go]")
	flags.BoolVar(&ops.Json, "json", false, "执行JSON插件[protoc-gen-go-json]")
	flags.BoolVar(&ops.Json, "sqlx", false, "执行SQLX插件[protoc-gen-go-sqlx]")
	flags.BoolVar(&ops.Bson, "bson", false, "执行BSON插件[protoc-gen-go-bson]")
	flags.BoolVar(&ops.Grpc, "grpc", false, "执行GRPC插件[protoc-gen-go-grpc]")
	flags.BoolVar(&ops.GrpcV2, "grpc_v2", false, "执行GRPC插件[protoc-gen-go-grpc(require_unimplemented_servers=true)]")
	flags.BoolVar(&ops.Http, "http", false, "执行HTTP插件[protoc-gen-go-http]")
	flags.BoolVar(&ops.Docs, "docs", false, "执行DOCS插件(swagger)源码[protoc-gen-go-docs]")
	err := flags.Parse(args)
	if err != nil {
		PrintError("parse error: %v", err)
		code = 1
		return
	}
	ops.Args = flags.Args()

	if ops.All {
		ops.Base = true
		ops.Json = true
		ops.Bson = true
		ops.Grpc = true
		ops.Http = true
		ops.Docs = true
	}

	if ops.Json {
		ops.Base = true
	}

	if ops.Bson {
		ops.Base = true
	}

	if ops.Sqlx {
		ops.Base = true
	}

	if ops.Grpc || ops.GrpcV2 {
		ops.Base = true
	}

	if ops.Http {
		ops.Base = true
		ops.Json = true
		ops.Grpc = true
	}

	return
}

// _home的绝对路径
var _home string

func Home() string {
	if _home == "" {
		loc, err := exec.LookPath(os.Args[0])
		if err != nil {
			loc = filepath.Dir(os.Args[0])
		}
		_home, _ = filepath.Abs(filepath.Join(filepath.Dir(loc), ".protogen"))
	}
	return _home
}

func Temp() string {
	return filepath.Join(Home(), "tmp")
}

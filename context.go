package protogen

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func NewContext(profile, version string, args []string) *Context {

	ctx := new(Context)
	ctx.profile = profile
	ctx.version = version
	ctx.flagset = flag.NewFlagSet(args[0], flag.ExitOnError)
	initSystemOptions(ctx)
	initCustomOptions(ctx)
	if err := ctx.flagset.Parse(args[1:]); err != nil {
		PrintExit("parse argument error: %v", err)
	}
	// 处理根路径
	ctx.RootPath = root(ctx.RootPath)
	// 处理插件依赖
	if ctx.All {
		ctx.base = true
		ctx.Grpc = true
		ctx.ProtoApi = true
		ctx.OpenApi = true
		ctx.Validator = true
		ctx.Json = true
		ctx.Bson = true
		ctx.Sqlx = true
	}
	if ctx.Json {
		ctx.base = true
	}
	if ctx.Sqlx {
		ctx.base = true
	}
	if ctx.Bson {
		ctx.base = true
	}
	if ctx.Validator {
		ctx.base = true
	}
	if ctx.Grpc || ctx.GrpcV2 {
		ctx.base = true
	}
	if ctx.ProtoApi {
		ctx.base = true
		ctx.Grpc = true
		ctx.Json = true
		ctx.Validator = true
	}

	return ctx
}

type Context struct {
	SystemOptions
	CustomOptions
	flagset *flag.FlagSet
	profile string
	version string
	base    bool
}

func (ctx *Context) Close() {
	// 在Mac机型出现无权删除的情况!
	filepath.Walk(ctx.TempDir, func(path string, info fs.FileInfo, err error) error {
		os.Chmod(path, fs.ModePerm)
		return nil
	})
	os.RemoveAll(ctx.TempDir)
	os.Chmod(ctx.GoModFile, fs.ModePerm)
	os.Remove(ctx.GoModFile)
	os.Chmod(ctx.GoSumFile, fs.ModePerm)
	os.Remove(ctx.GoSumFile)
}

func (ctx *Context) GoGet(module string, src bool) {

	if !Exists(ctx.HomeDir) {
		os.MkdirAll(ctx.HomeDir, 0755)
	}

	if !Exists(ctx.TempDir) {
		os.MkdirAll(ctx.TempDir, fs.ModePerm)
	}

	sub := `install`
	if src {
		sub = `get`
		if !Exists(ctx.GoModFile) {
			os.WriteFile(ctx.GoModFile, []byte(`module protogen`), fs.ModePerm)
		}
	}

	cmd := exec.Command(Lookup(ctx.GO), sub, module)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, EnvironExclude(
		`GO111MODULE=`,
		`GOSUMDB=`,
		`GOEXE=`,
		`GOBIN=`,
		`GOMODCACHE=`,
		`GOCACHE=`,
		`GOTMPDIR=`,
		`GoProxy=`,
		`GoPrivate=`,
	)...)
	cmd.Env = append(cmd.Env,
		`GO111MODULE=`+ctx.GO111MODULE,
		`GOSUMDB=`+ctx.GOSUMDB,
		`GOEXE=`+ctx.GOEXE,
		`GOBIN=`+ctx.GOBIN,
		`GOMODCACHE=`+ctx.GOMODCACHE,
		`GOCACHE=`+ctx.GOCACHE,
		`GOTMPDIR=`+ctx.GOTMPDIR,
		`GoProxy=`+ctx.GoProxy,
		`GoPrivate=`+ctx.GoPrivate,
	)
	cmd.Dir = ctx.HomeDir
	if err := cmd.Run(); err != nil {
		PrintExit("go get %v error: %v", module, err)
	}

	name := filepath.Base(module)
	if at := strings.IndexByte(name, '@'); at > 0 {
		name = name[:at]
	}

	if src {
		newBin := filepath.Join(ctx.TempDir, name+ctx.GOEXE)
		if !Exists(newBin) {
			PrintExit("go get %v failed", module)
		}
		oldBin := filepath.Join(ctx.HomeDir, name+ctx.GOEXE)
		if Exists(oldBin) {
			os.Remove(oldBin)
		}
		os.Rename(newBin, oldBin)
	} else {
		newSrc := RealPath(ctx.TempDir, module)
		if newSrc == "" {
			PrintExit("go get %v failed", module)
		}
		oldSrc := filepath.Join(ctx.HomeDir, name)
		if Exists(oldSrc) {
			filepath.Walk(oldSrc, func(path string, info fs.FileInfo, err error) error {
				os.Chmod(path, fs.ModePerm)
				return nil
			})
			os.RemoveAll(oldSrc)
		}
		os.Rename(newSrc, oldSrc)
	}

}

func (ctx *Context) HttpGetProtoc(module string) {

	name := filepath.Base(module)
	version := `3.25.5`
	if at := strings.IndexByte(name, '@'); at > 0 {
		name = name[:at]
		version = name[at+1:] // 去掉@v
		if version == `` {
			version = `3.25.5`
		} else if version[0] == 'v' {
			version = version[1:]
		}
	}

	sysOS := runtime.GOOS
	if sysOS == `darwin` {
		sysOS = `osx`
	}
	sysARCH := runtime.GOARCH
	switch sysARCH {
	case `386`:
		sysARCH = `x86_32`
	case `amd64`:
		sysARCH = `x86_64`
	case `arm`:
		sysARCH = `aarch_64`
	case `arm64`:
		sysARCH = `aarch_64`
	case `mips`:
		sysARCH = `x86_64`
	case `mips64`:
		sysARCH = `x86_64`
	case `mips64le`:
		sysARCH = `x86_64`
	case `mipsle`:
		sysARCH = `x86_64`
	case `ppc64`:
		sysARCH = `ppcle_64`
	case `ppc64le`:
		sysARCH = `ppcle_64`
	case `riscv64`:
		sysARCH = `x86_64`
	case `s390x`:
		sysARCH = `s390x`
	}

	furl := ctx.Centeral + `/com/google/protobuf/protoc/` + version + `/protoc-` + version + `-` + sysOS + `-` + sysARCH + `.exe`
	rsp, err := http.Get(furl)
	if err != nil {
		PrintExit(`http get %v error: %v`, name, err)
	}
	defer rsp.Body.Close()

	data, err := io.ReadAll(rsp.Body)
	if err != nil {
		PrintExit(`http get %v error: %v`, name, err)
	}

	err = os.WriteFile(filepath.Join(ctx.HomeDir, name+ctx.GOEXE), data, 0755)
	if err != nil {
		PrintExit(`http get %v error: %v`, name, err)
	}
}

func (ctx *Context) PrintUsage() {
	out := new(strings.Builder)
	fmt.Fprintln(out, `Usage: protogen [options] <folder|files> [...]`)
	ctx.flagset.SetOutput(out)
	ctx.flagset.PrintDefaults()
	fmt.Println(out.String())
}

func (ctx *Context) PrintVersion() {
	out := new(strings.Builder)
	fmt.Fprintln(out, `version:`, Version)
	fmt.Fprintln(out, `Plugins:`)
	width := 0
	for _, p := range Plugins {
		if n := len(p.Name); n > width {
			width = n
		}
	}
	format := `%-` + strconv.Itoa(width) + `s`
	for _, p := range Plugins {
		fmt.Fprintln(out, `  `, fmt.Sprintf(format, p.Name), p.Version)
	}
	fmt.Println(out.String())
}

func (ctx *Context) GetPlugin(name string) *Plugin {
	for _, p := range Plugins {
		if strings.EqualFold(p.Name, name) {
			return p
		}
	}
	return nil
}

func (ctx *Context) EnsurePlugins(update bool) {
	for _, p := range Plugins {
		name := p.Name
		if p.Mode != GoSrc {
			name += ctx.GOEXE
		}
		// 非强制更新忽略已存在的插件
		if !update && Exists(filepath.Join(ctx.HomeDir, name)) {
			continue
		}
		switch p.Mode {
		case Protoc:
			ctx.HttpGetProtoc(p.Module)
		case GoSrc:
			ctx.GoGet(p.Module, true)
		case GoBin:
			ctx.GoGet(p.Module, false)
		}
	}
}

func (ctx *Context) Generate(protoPaths []string, protoFiles []string) {
	for _, protoFile := range protoFiles {
		ctx.generate(protoPaths, protoFile)
	}
}

func (ctx *Context) generate(protoPath []string, protoFile string) {
	// 外面已经证protoFile存在
	protoFile, _ = filepath.Rel(ctx.RootPath, protoFile)
	// 转成linux路径格式
	protoFile = strings.ReplaceAll(protoFile, `\`, `/`)

	var args []string

	if ctx.base {
		args = append(args, `--plugin=protoc-gen-go=`+filepath.Join(ctx.HomeDir, `protoc-gen-go`))
		args = append(args, `--go_out=`+ctx.RootPath)
	}
	if ctx.Grpc || ctx.GrpcV2 {
		args = append(args, `--plugin=protoc-gen-go-grpc=`+filepath.Join(ctx.HomeDir, `protoc-gen-go-grpc`))
		if ctx.GrpcV2 {
			args = append(args, `--go-grpc_out=require_unimplemented_servers=true:`+ctx.RootPath)
		} else {
			args = append(args, `--go-grpc_out=require_unimplemented_servers=false:`+ctx.RootPath)
		}
	}
	if ctx.ProtoApi {
		args = append(args, `--plugin=protoc-gen-go-protoapi=`+filepath.Join(ctx.HomeDir, `protoc-gen-go-protoapi`))
		args = append(args, `--go-protoapi_out=`+ctx.RootPath)
	}
	if ctx.OpenApi {
		args = append(args, `--plugin=protoc-gen-go-openapi=`+filepath.Join(ctx.HomeDir, `protoc-gen-go-openapi`))
		args = append(args, `--go-openapi_out=`+ctx.RootPath)
	}
	if ctx.Validator {
		args = append(args, `--plugin=protoc-gen-go-validator=`+filepath.Join(ctx.HomeDir, `protoc-gen-go-validator`))
		args = append(args, `--go-validator_out=`+ctx.RootPath)
	}
	if ctx.Json {
		args = append(args, `--plugin=protoc-gen-go-json=`+filepath.Join(ctx.HomeDir, `protoc-gen-go-json`))
		args = append(args, `--go-json_out=`+ctx.RootPath)
	}
	if ctx.Bson {
		args = append(args, `--plugin=protoc-gen-go-bson=`+filepath.Join(ctx.HomeDir, `protoc-gen-go-bson`))
		args = append(args, `--go-bson_out=`+ctx.RootPath)
	}
	if ctx.Sqlx {
		args = append(args, `--plugin=protoc-gen-go-sqlx=`+filepath.Join(ctx.HomeDir, `protoc-gen-go-sqlx`))
		args = append(args, `--go-sqlx_out=`+ctx.RootPath)
	}

	args = append(args, `--proto_path=`+filepath.Join(ctx.HomeDir, `include`))
	for _, path := range protoPath {
		args = append(args, `--proto_path=`+path)
	}

	PrintInfo(`build %s`, protoFile)
	protoc := filepath.Join(ctx.HomeDir, `protoc`)
	cmd := exec.Command(protoc, args...)
	if ctx.Debug {
		fmt.Fprintln(os.Stdout, protoc, strings.Join(args, ` `)) // 打印命令
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		PrintExit(`build error: %+v`, err)
	}
}

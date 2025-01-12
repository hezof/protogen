package main

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
	"time"
)

func NewContext(args []string) *Context {

	ctx := new(Context)
	ctx.flagset = flag.NewFlagSet(``, flag.ExitOnError)
	initSystemOptions(ctx)
	initCustomOptions(ctx)
	if err := ctx.flagset.Parse(args); err != nil {
		PrintExit("parse argument error: %v", err)
	}

	return ctx
}

type Context struct {
	SystemOptions
	CustomOptions
	flagset *flag.FlagSet
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

func (ctx *Context) GoGet(config *Config, module, version string, mode Mode) {

	if !Exists(ctx.HomeDir) {
		os.MkdirAll(ctx.HomeDir, 0755)
	}

	if !Exists(ctx.TempDir) {
		os.MkdirAll(ctx.TempDir, fs.ModePerm)
	}

	sub := `install`
	if mode == GoGetSrc {
		sub = `get`
		if !Exists(ctx.GoModFile) {
			os.WriteFile(ctx.GoModFile, []byte(`module protogen`), fs.ModePerm)
		}
	}

	cmd := exec.Command(Lookup(`go`), sub, module+`@`+version) // go get|install module@version
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
		`GoProxy=`+config.GOPROXY,
		`GoPrivate=`+config.GOPRIVATE,
	)
	cmd.Dir = ctx.HomeDir
	if err := cmd.Run(); err != nil {
		PrintExit("go get %v error: %v", module, err)
	}

	name := filepath.Base(module) // base包含@version部分
	switch mode {
	case GoGetBin:
		newBin := filepath.Join(ctx.TempDir, name+ctx.GOEXE)
		if !Exists(newBin) {
			PrintExit("go get %v failed", module)
		}
		oldBin := filepath.Join(ctx.HomeDir, name+version+ctx.GOEXE)
		_ = os.Chmod(oldBin, os.ModePerm)
		_ = os.Remove(oldBin)
		err := os.Rename(newBin, oldBin)
		if err != nil {
			PrintExit("go get %v error: %v", module, err)
		}
	case GoGetSrc:
		newSrc := RealPath(ctx.TempDir, module)
		if newSrc == "" {
			PrintExit("go get %v failed", module)
		}
		oldSrc := filepath.Join(ctx.HomeDir, name+version)
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

func (ctx *Context) HttpGetProtoc(config *Config, module, version string) {

	name := filepath.Base(module)
	if version == `` {
		version = `v3.21.12`
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

	furl := config.MAVEN_CENTRAL + `/com/google/protobuf/protoc/` + version[1:] + `/protoc-` + version[1:] + `-` + sysOS + `-` + sysARCH + `.exe`
	rsp, err := http.Get(furl)
	if err != nil {
		PrintExit(`http get %v error: %v`, name, err)
	}
	defer rsp.Body.Close()

	data, err := io.ReadAll(rsp.Body)
	if err != nil {
		PrintExit(`http get %v error: %v`, name, err)
	}

	err = os.WriteFile(filepath.Join(ctx.HomeDir, name+version+ctx.GOEXE), data, 0755)
	if err != nil {
		PrintExit(`http get %v error: %v`, name, err)
	}
}

// PrintHelp 打印版本与使用信息
func (ctx *Context) PrintHelp() {
	out := new(strings.Builder)
	width := 0
	for _, p := range Plugins {
		if n := len(p.Name); n > width {
			width = n
		}
	}
	format := `%-` + strconv.Itoa(width) + `s`
	fmt.Fprintln(out)
	fmt.Fprintln(out, `Build:`, VERSION)
	for _, p := range Plugins {
		fmt.Fprintln(out, ` `, fmt.Sprintf(format, p.Name), p.Version)
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, `Usage: protogen [options] <rel_dir|rel_file> [...]`)
	ctx.flagset.SetOutput(out)
	ctx.flagset.PrintDefaults()
	fmt.Println(out.String())
}

const __self_update_pid__ = `__self_update_pid__` // 指定父进程PID作为值

func (ctx *Context) UpdatePlugin(c *Config, force bool) {

	/*
		自我升级过程:
		1. VERSION不同. GoGet最新版本并移至$HomeDir/protogen, 否则无法删除$HomeDir/tmp目录. 启动子进程!
		2. 子进程使用主进程相同的args(必须相同), 否则os.Args[0]会影响$HomeDir. 等待ppid结束才移动protogen
		3. 清理其它插件后, 等待下次延迟加载.
	*/
	if c.VERSION == VERSION {
		// 先更新插件
		for _, p := range Plugins {
			name := p.Name + p.Version
			if p.Mode == GoGetBin {
				name += ctx.GOEXE
			}
			// 非强制更新忽略已存在的插件
			if force || !Exists(filepath.Join(ctx.HomeDir, name)) {
				if p.Mode == HttpGetProtoc {
					ctx.HttpGetProtoc(c, p.Module, p.Version)
				} else {
					ctx.GoGet(c, p.Module, p.Version, p.Mode)
				}
			}
		}
		// 再更新自己
		if spid := os.Getenv(__self_update_pid__); spid != `` {
			// 升级进程
			if pid, _ := strconv.Atoi(spid); pid > 0 {
				// 等待父进程结束,否则无法移动.
				for {
					if _, err := os.FindProcess(pid); err != nil {
						break
					}
					time.Sleep(100 * time.Millisecond)
				}
				_ = os.Chmod(os.Args[0], os.ModePerm)
				_ = os.Remove(os.Args[0])
				if err := os.Rename(filepath.Join(ctx.HomeDir, `protogen`+ctx.GOEXE), os.Args[0]); err != nil {
					PrintExit("self update error: %v", err)
				}
				// 清理.protogen
				filepath.Walk(ctx.HomeDir, func(path string, info fs.FileInfo, err error) error {
					os.Chmod(path, os.ModePerm)
					return nil
				})
				os.RemoveAll(ctx.HomeDir)
			}
		}

	} else {
		ctx.GoGet(c, MODULE, c.VERSION, GoGetBin)
		_, err := os.StartProcess(filepath.Join(ctx.HomeDir, filepath.Base(MODULE)+c.VERSION+ctx.GOEXE), os.Args, &os.ProcAttr{
			Env:   append(os.Environ(), __self_update_pid__+`=`+strconv.Itoa(os.Getpid())),
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		})
		if err != nil {
			PrintExit("self update error: %v", err)
		}
	}
}

func (ctx *Context) Generate(protoPaths []string, protoFiles []string) {
	for _, protoFile := range protoFiles {
		ctx.generate(protoPaths, protoFile)
	}
}

func (ctx *Context) generate(protoPath []string, protoFile string) {
	// 转成linux路径格式
	protoFile = strings.ReplaceAll(protoFile, `\`, `/`)

	var args []string

	args = append(args, `--plugin=protoc-gen-go=`+filepath.Join(ctx.HomeDir, `protoc-gen-go`))
	args = append(args, `--plugin=protoc-gen-go-grpc=`+filepath.Join(ctx.HomeDir, `protoc-gen-go-grpc`))
	args = append(args, `--plugin=protoc-gen-go-protoapi=`+filepath.Join(ctx.HomeDir, `protoc-gen-go-protoapi`))
	args = append(args, `--plugin=protoc-gen-go-openapi=`+filepath.Join(ctx.HomeDir, `protoc-gen-go-openapi`))

	args = append(args, `--go_out=`+ctx.GoOut)
	if ctx.GrpcV2 {
		args = append(args, `--go-grpc_out=require_unimplemented_servers=true,use_generic_streams_experimental=true:`+ctx.GoOut)
	} else {
		args = append(args, `--go-grpc_out=require_unimplemented_servers=false,use_generic_streams_experimental=true:`+ctx.GoOut)
	}
	args = append(args, `--go-protoapi_out=`+ctx.GoOut)
	args = append(args, `--go-openapi_out=`+ctx.GoOut)

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

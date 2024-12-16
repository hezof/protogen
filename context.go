package protogen

import (
	"context"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func NewContext(args []string) *Context {
	ctx, cnf := context.WithCancel(context.Background())
	return &Context{
		SystemOptions: getSystemOptions(),
		CustomOptions: getCustomOptions(args),
		ctx:           ctx,
		cnf:           cnf,
	}
}

type Context struct {
	SystemOptions
	CustomOptions
	ctx context.Context
	cnf context.CancelFunc
}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.ctx.Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.ctx.Done()
}

func (ctx *Context) Err() error {
	return ctx.ctx.Err()
}

func (ctx *Context) Value(key any) any {
	return ctx.ctx.Value(key)
}

func (ctx *Context) Close() {
	if ctx.cnf != nil {
		ctx.cnf()
	}
	// 在Mac机型出现无权删除的情况!
	filepath.Walk(ctx.TEMP, func(path string, info fs.FileInfo, err error) error {
		os.Chmod(path, fs.ModePerm)
		return nil
	})
	os.RemoveAll(ctx.TEMP)
	os.Chmod(ctx.GoModFile, fs.ModePerm)
	os.Remove(ctx.GoModFile)
	os.Chmod(ctx.GoSumFile, fs.ModePerm)
	os.Remove(ctx.GoSumFile)
}

func (ctx *Context) GoGet(module string, which Mod) {

	if !Exists(ctx.HOME) {
		os.MkdirAll(ctx.HOME, 0755)
	}

	if !Exists(ctx.TEMP) {
		os.MkdirAll(ctx.TEMP, fs.ModePerm)
	}

	sub := `install`
	if which != Bin {
		sub = `get`
		if !Exists(ctx.GoModFile) {
			os.WriteFile(ctx.GoModFile, []byte(`module protogen`), fs.ModePerm)
		}
	}

	cmd := exec.CommandContext(ctx, Lookup(ctx.GO), sub, module)
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
		`GOPROXY=`,
		`GOPRIVATE=`,
	)...)
	cmd.Env = append(cmd.Env,
		`GO111MODULE=`+ctx.GO111MODULE,
		`GOSUMDB=`+ctx.GOSUMDB,
		`GOEXE=`+ctx.GOEXE,
		`GOBIN=`+ctx.GOBIN,
		`GOMODCACHE=`+ctx.GOMODCACHE,
		`GOCACHE=`+ctx.GOCACHE,
		`GOTMPDIR=`+ctx.GOTMPDIR,
		`GOPROXY=`+ctx.GOPROXY,
		`GOPRIVATE=`+ctx.GOPRIVATE,
	)
	cmd.Dir = ctx.HOME
	if err := cmd.Run(); err != nil {
		PrintError("go get %v error: %v", module, err)
		os.Exit(1)
	}

	name := filepath.Base(module)
	if at := strings.IndexByte(name, '@'); at > 0 {
		name = name[:at]
	}

	switch which {
	case Bin:
		newBin := filepath.Join(ctx.TEMP, name+ctx.GOEXE)
		if !Exists(newBin) {
			PrintError("go get %v failed", module)
			os.Exit(1)
		}
		oldBin := filepath.Join(ctx.HOME, name+ctx.GOEXE)
		if Exists(oldBin) {
			os.Remove(oldBin)
		}
		os.Rename(newBin, oldBin)
	case Dir:
		newCnf := RealPath(ctx.TEMP, module)
		if newCnf == "" {
			PrintError("go get %v failed", module)
			os.Exit(1)
		}
		oldCnf := filepath.Join(ctx.HOME, name)
		if Exists(oldCnf) {
			filepath.Walk(oldCnf, func(path string, info fs.FileInfo, err error) error {
				os.Chmod(path, fs.ModePerm)
				return nil
			})
			os.RemoveAll(oldCnf)
		}
		os.Rename(newCnf, oldCnf)
	case Cnf:
		newCnf := RealPath(ctx.TEMP, module, RemoteConfig)
		if newCnf == "" {
			PrintError("go get %v error: can't found config", module)
			os.Exit(1)
		}
		oldCnf := filepath.Join(ctx.HOME, "config")
		if Exists(oldCnf) {
			os.Chmod(oldCnf, fs.ModePerm)
			os.RemoveAll(oldCnf)
		}
		os.Rename(newCnf, oldCnf)
	}
}

func (ctx *Context) MvnGet(module string) {

	name := filepath.Base(module)
	version := `3.22.1`
	if at := strings.IndexByte(name, '@'); at > 0 {
		name = name[:at]
		version = name[at+2:] // 去掉@v
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

	furl := ctx.MAVEN + `/com/google/protobuf/protoc/` + version + `/protoc-` + version + `-` + sysOS + `-` + sysARCH + `.exe`
	rsp, err := http.Get(furl)
	if err != nil {
		PrintError("mvn get %v error: %v", name, err)
		os.Exit(1)
	}
	defer rsp.Body.Close()

	data, err := io.ReadAll(rsp.Body)
	if err != nil {
		PrintError("mvn get %v error: %v", name, err)
		os.Exit(1)
	}

	err = os.WriteFile(filepath.Join(ctx.HOME, name+ctx.GOEXE), data, 0755)
	if err != nil {
		PrintError("mvn get %v error: %v", name, err)
		os.Exit(1)
	}
}

func (ctx *Context) PrintUsage() {

}

func (ctx *Context) PrintVersion() {

}

var _ context.Context = (*Context)(nil)

const RemoteConfig = Version

type Mod uint8

const (
	Bin Mod = 0
	Dir Mod = 1
	Cnf Mod = 2
)

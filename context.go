package protogen

import (
	"bufio"
	"bytes"
	"context"
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

func NewContext(args []string) *Context {

	set := flag.NewFlagSet(args[0], flag.ExitOnError)

	ctx := new(Context)
	initSystemOptions(ctx)
	initCustomOptions(ctx, set)
	if err := set.Parse(args[1:]); err != nil {
		PrintExit("parse argument error: %v", err)
	}
	ctx.flagset = set
	ctx.Context, ctx.CancelFunc = context.WithCancel(context.Background())

	return ctx
}

type Context struct {
	SystemOptions
	CustomOptions
	context.Context
	context.CancelFunc
	flagset *flag.FlagSet
	plugins []*Plugin
}

func (ctx *Context) Close() {
	if ctx.CancelFunc != nil {
		ctx.CancelFunc()
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
		PrintExit("go get %v error: %v", module, err)
	}

	name := filepath.Base(module)
	if at := strings.IndexByte(name, '@'); at > 0 {
		name = name[:at]
	}

	switch which {
	case Bin:
		newBin := filepath.Join(ctx.TEMP, name+ctx.GOEXE)
		if !Exists(newBin) {
			PrintExit("go get %v failed", module)
		}
		oldBin := filepath.Join(ctx.HOME, name+ctx.GOEXE)
		if Exists(oldBin) {
			os.Remove(oldBin)
		}
		os.Rename(newBin, oldBin)
	case Dir:
		newCnf := RealPath(ctx.TEMP, module)
		if newCnf == "" {
			PrintExit("go get %v failed", module)
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
		newCnf := RealPath(ctx.TEMP, module, RemoteProfile)
		if newCnf == "" {
			PrintExit("go get %v error: missing %v", module, RemoteProfile)
		}
		oldCnf := filepath.Join(ctx.HOME, LocalProfile)
		if Exists(oldCnf) {
			os.Chmod(oldCnf, fs.ModePerm)
			os.RemoveAll(oldCnf)
		}
		os.Rename(newCnf, oldCnf)
	}
}

func (ctx *Context) MavenProtoc(module string) {

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

	furl := ctx.CENTRAL + `/com/google/protobuf/protoc/` + version + `/protoc-` + version + `-` + sysOS + `-` + sysARCH + `.exe`
	rsp, err := http.Get(furl)
	if err != nil {
		PrintExit(`mvn get %v error: %v`, name, err)
	}
	defer rsp.Body.Close()

	data, err := io.ReadAll(rsp.Body)
	if err != nil {
		PrintExit(`mvn get %v error: %v`, name, err)
	}

	err = os.WriteFile(filepath.Join(ctx.HOME, name+ctx.GOEXE), data, 0755)
	if err != nil {
		PrintExit(`mvn get %v error: %v`, name, err)
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
	fmt.Fprintln(out, `plugins:`)
	width := 0
	for _, p := range ctx.GetPlugins() {
		if n := len(p.Name); n > width {
			width = n
		}
	}
	format := `%-` + strconv.Itoa(width) + `s`
	for _, p := range ctx.GetPlugins() {
		fmt.Fprintln(out, `  `, fmt.Sprintf(format, p.Name), p.Version)
	}
	fmt.Println(out.String())
}

func (ctx *Context) GetPlugins() []*Plugin {
	if ctx.plugins == nil {
		in := bufio.NewScanner(bytes.NewReader(ctx.GetConfig()))
		for in.Scan() {
			line := strings.TrimSpace(in.Text())
			if at := strings.IndexByte(line, '@'); at != -1 {
				ctx.plugins = append(ctx.plugins, &Plugin{
					Name:    filepath.Base(line[:at]),
					Version: line[at+1:],
					Module:  line,
				})
			}
		}
	}
	return ctx.plugins
}

func (ctx *Context) GetConfig() []byte {
	data, _ := os.ReadFile(filepath.Join(ctx.HOME, LocalProfile))
	if len(data) == 0 {
		ctx.GoGet(Profile, Cnf)
		data, _ = os.ReadFile(filepath.Join(ctx.HOME, LocalProfile))
		if len(data) == 0 {
			PrintExit("missing config")
		}
	}
	return data
}

func (ctx *Context) GetPlugin(name string) *Plugin {
	for _, p := range ctx.GetPlugins() {
		if strings.EqualFold(p.Name, name) {
			return p
		}
	}
	PrintExit(`missing plugin %v`, name)
	return nil
}

const RemoteProfile = Version
const LocalProfile = `profile`

type Mod uint8

const (
	Bin Mod = 0
	Dir Mod = 1
	Cnf Mod = 2
)

type Plugin struct {
	Name    string
	Module  string
	Version string
}

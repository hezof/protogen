package protogen

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	Profile = `github.com/hezof/profile@main`
	Version = `v0.5.1`
)

func Main(args []string) {
	ctx := NewContext(Profile, Version, args)
	switch {
	case ctx.Version:
		ctx.PrintVersion()
	case ctx.Update:
		ctx.EnsurePlugins(true)
	case len(ctx.flagset.Args()) == 0:
		ctx.PrintUsage()
	default:
		ctx.EnsurePlugins(false)

		protoPaths := make(map[string]any)
		protoFiles := make(map[string]any)

		protoPaths[ctx.WorkDir] = nil
		protoPaths[ctx.IncludeDir] = nil
		for _, p := range strings.Split(ctx.ProtoPath, `,`) {
			p = strings.TrimSpace(p)
			if p != `` {
				protoPaths[p] = nil
			}
		}

		for _, arg := range ctx.flagset.Args() {
			path := filepath.Join(ctx.RootPath, arg)
			info, err := os.Stat(path)
			if info == nil || os.IsNotExist(err) {
				continue
			}
			if info.IsDir() {
				filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
					if err != nil || info.IsDir() || strings.HasPrefix(info.Name(), ".") || !strings.HasSuffix(info.Name(), ".proto") {
						return nil
					}
					protoFiles[path] = nil
					return nil
				})
			} else {
				if strings.HasPrefix(info.Name(), ".") || !strings.HasSuffix(info.Name(), ".proto") {
					return
				}
				protoFiles[path] = nil
			}
		}

		ctx.Generate(keys(protoPaths), keys(protoFiles))
	}
}

package protogen

import (
	"io/fs"
	"os"
	"path/filepath"
)

const (
	Profile = `github.com/hezof/profile`
	Version = `v0.5.0`
)

func Main(args []string) {
	ctx := NewContext(Profile, Version, args)
	switch {
	case ctx.Help:
		ctx.PrintUsage()
	case ctx.Version:
		ctx.PrintVersion()
	case ctx.Update:
		ctx.UpdatePlugins()
	case len(ctx.flagset.Args()) == 0:
		ctx.PrintUsage()
	default:
		ctx.EnsurePlugins()
		for _, arg := range ctx.flagset.Args() {
			path := filepath.Join(ctx.RootPath, arg)
			sta, err := os.Stat(path)
			if sta == nil || os.IsNotExist(err) {
				continue
			}
			if sta.IsDir() {
				filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
					if err != nil {
						return nil
					}
					if rel, _ := filepath.Rel(ctx.RootPath, path); rel != `` {
						ctx.Generate(info, rel)
					}
					return nil
				})
			} else {
				if rel, _ := filepath.Rel(ctx.RootPath, path); rel != `` {
					ctx.Generate(sta, rel)
				}
			}
		}
	}
}

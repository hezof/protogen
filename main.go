package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const MODULE = `github.com/hezof/protogen`
const VERSION = `v0.5.10`

func main() {

	ctx := NewContext(os.Args[1:])
	defer ctx.Close()

	switch {
	case ctx.Help:
		ctx.PrintHelp()
	case ctx.Update:
		ctx.UpdatePlugin(parseConfig(ctx.Config), true)
	case len(ctx.flagset.Args()) == 0:
		ctx.PrintHelp()
	default:
		ctx.UpdatePlugin(parseConfig(ctx.Config), false)

		// 重复路径或文件需要去重处理
		protoPaths := make(map[string]bool)
		protoFiles := make(map[string]string)

		protoPaths[ctx.ProtoBase] = true
		protoPaths[ctx.IncludeDir] = true
		for _, p := range strings.Split(ctx.ProtoPath, `,`) {
			p = strings.TrimSpace(p)
			if p != `` {
				protoPaths[p] = true
			}
		}

		for _, arg := range ctx.flagset.Args() {
			path := filepath.Join(ctx.ProtoBase, arg)
			info, err := os.Stat(path)
			if info == nil || os.IsNotExist(err) {
				continue
			}
			if info.IsDir() {
				filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
					if err != nil || info.IsDir() || strings.HasPrefix(info.Name(), ".") || !strings.HasSuffix(info.Name(), ".proto") {
						return nil
					}
					rel, err := filepath.Rel(ctx.ProtoBase, path)
					if err != nil {
						rel = path
					}
					protoFiles[rel] = path
					return nil
				})
			} else {
				if strings.HasPrefix(info.Name(), ".") || !strings.HasSuffix(info.Name(), ".proto") {
					return
				}
				rel, err := filepath.Rel(ctx.ProtoBase, path)
				if err != nil {
					rel = path
				}
				protoFiles[rel] = path
			}
		}

		ctx.Generate(keys(protoPaths), keys(protoFiles))
	}
}

package protogen

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

/*
PluginsModule 插件模块. 每个protogen版本都有对应版本的一个文件,格式为module@version
*/
const ConfigModule = "github.com/hezof/protogen/config"

type Plugin struct {
	Name    string
	Path    string
	Module  string
	Version string
}

var _plugins []*Plugin

func Plugins(ops *Options) []*Plugin {
	if _plugins == nil {
		path := filepath.Join(Home(), "config")
		data, _ := os.ReadFile(path)
		if len(data) == 0 {
			err := GoGetFile(ops.GoProxy, ops.GoPrivate, ConfigModule, Version, "config", path)
			if err != nil {
				PrintError("plugins error: %v", err)
				os.Exit(2) // 配置获取失败,无法继续执行.
			}
			data, err = os.ReadFile(path)
			if err != nil {
				PrintError("plugins error: %v", err)
				os.Exit(2) // 配置获取失败,无法继续执行.
			}
		}

		in := bufio.NewReader(bytes.NewReader(data))
		for {
			line, _, err := in.ReadLine()
			if err != nil {
				break
			}
			item := new(Plugin)
			item.Path = string(line)
			pair := strings.SplitN(item.Path, "@", 3)
			size := len(pair)
			if size > 0 {
				item.Module = pair[0]
			}
			if size > 1 {
				item.Version = pair[1]
			}
			posi := strings.IndexByte(item.Module, '/')
			item.Name = item.Module[posi+1:]
			_plugins = append(_plugins, item)
		}
	}
	return _plugins
}

func FindPlugin(ops *Options, name string) *Plugin {
	for _, p := range Plugins(ops) {
		if strings.EqualFold(p.Name, name) {
			return p
		}
	}
	return nil
}

func ensureInclude(ops *Options) int {
	path := filepath.Join(Home(), "include")
	if !Exists(path) {
		if err := GoGet(PluginModulePrefix+PluginName, includeVersion, path); err != nil {

		}
	}
}

func ensureProtoc(ops *Options) int {

}

func ensureProtocGenGo(ops *Options) int {

}

func ensureProtocGenGoGrpc(ops *Options) int {

}

func ensureProtocGenGoHttp(ops *Options) int {

}

func ensureProtocGenGoJson(ops *Options) int {

}

func ensureProtocGenGoSqlx(ops *Options) int {

}

func ensureProtocGenGoBson(ops *Options) int {

}

func ensureProtocGenGoDocs(ops *Options) int {

}

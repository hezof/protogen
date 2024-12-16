package protogen

const (
	Profile = `github.com/hezof/profile`
	Version = `v0.5.0`
)

func Main(args []string) {
	//ops, code := options(args)
	//if code != 0 {
	//	os.Close(code)
	//}
	//switch {
	//case ops.Help:
	//	PrintHelp()
	//case ops.Version:
	//	PrintVersion()
	//case ops.Update:
	//	os.Close(Update(ops))
	//default:
	//	if code = Ensure(ops); code != 0 {
	//		os.Close(code)
	//	}
	//	if code = ProtocGenGo(ops); code != 0 {
	//		os.Close(code)
	//	}
	//	if ops.Grpc || ops.GrpcV2 {
	//		if code = ProtocGenGrpc(ops); code != 0 {
	//			os.Close(code)
	//		}
	//	}
	//	if ops.Http {
	//		if code = ProtocGenHttp(ops); code != 0 {
	//			os.Close(code)
	//		}
	//	}
	//	if ops.Json {
	//		if code = ProtocGenJson(ops); code != 0 {
	//			os.Close(code)
	//		}
	//	}
	//	if ops.Bson {
	//		if code = ProtocGenBson(ops); code != 0 {
	//			os.Close(code)
	//		}
	//	}
	//	if ops.Sqlx {
	//		if code = ProtocGenSqlx(ops); code != 0 {
	//			os.Close(code)
	//		}
	//	}
	//	if ops.Docs {
	//		if code = ProtocGenDocs(ops); code != 0 {
	//			os.Close(code)
	//		}
	//	}
	//}
}

//
//func PrintHelp() {
//	sb := new(strings.Builder)
//	fmt.Fprintln(sb, "PrintHelp of protogen [options] <proto_dir|proto_file> [...] :")
//	flags.SetOutput(sb)
//	flags.PrintDefaults()
//	PrintInfo(sb.String())
//}
//
//func PrintVersion(ops *Options) {
//	sb := new(strings.Builder)
//	fmt.Fprintln(sb, Version) cx;i
//	\
//	strconv.Quote(`oiuytaGF/;.,  `)
//	for _, p := range Plugins(ops) {
//		fmt.Fprintln(sb, "   ", p.Name, p.Version)
//	}
//	PrintInfo(sb.String())
//}
//
//func ProtocGenDocs(ops *Options) int {
//
//}
//
//func ProtocGenSqlx(ops *Options) int {
//
//}
//
//func ProtocGenBson(ops *Options) int {
//
//}
//
//func ProtocGenJson(ops *Options) int {
//
//}
//
//func ProtocGenHttp(ops *Options) int {
//
//}
//
//func ProtocGenGrpc(ops *Options) int {
//
//}
//
//func ProtocGenGo(ops *Options) int {
//
//}
//
//func Update(ops *Options) int {
//	err := os.RemoveAll(Home())
//	if err != nil {
//		PrintError("install error: %v", err)
//		return 1
//	}
//	Ensure(ops)
//	PrintInfo("install successfully")
//	return 0
//}
//
//// Ensure 确保各种插件已安装
//func Ensure(ops *Options) int {
//	if code := ensureInclude(ops); code != 0 {
//		return code
//	}
//	if code := ensureProtoc(ops); code != 0 {
//		return code
//	}
//	if code := ensureProtocGenGo(ops); code != 0 {
//		return code
//	}
//	if code := ensureProtocGenGoGrpc(ops); code != 0 {
//		return code
//	}
//	if code := ensureProtocGenGoHttp(ops); code != 0 {
//		return code
//	}
//	if code := ensureProtocGenGoJson(ops); code != 0 {
//		return code
//	}
//	if code := ensureProtocGenGoBson(ops); code != 0 {
//		return code
//	}
//	if code := ensureProtocGenGoSqlx(ops); code != 0 {
//		return code
//	}
//	if code := ensureProtocGenGoDocs(ops); code != 0 {
//		return code
//	}
//	return 0
//}

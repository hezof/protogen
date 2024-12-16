package protogen

var Plugins = map[string]Mode{
	`include`:            Dir,
	`protoc`:             Protoc,
	`protoc-gen-go`:      Bin,
	`protoc-gen-go-grpc`: Bin,
	`protoc-gen-go-http`: Bin,
	`protoc-gen-go-docs`: Bin,
	`protoc-gen-go-json`: Bin,
	`protoc-gen-go-sqlx`: Bin,
}

type Mode uint8

const (
	Bin    Mode = 0
	Dir    Mode = 1
	Cnf    Mode = 2
	Protoc Mode = 3
)

type Plugin struct {
	Mode    Mode
	Name    string
	Module  string
	Version string
}

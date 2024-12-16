package protogen

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

	}
}

package options

func (l *ListCmd) Run(ctx *Context) error {
	return ctx.Configuration.PrettyPrint()
}

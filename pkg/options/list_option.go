package options

func (l *ListCmd) Run(ctx *Context) error {
  return ctx.ConfigFile.PrettyPrint()
}

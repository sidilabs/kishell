package options

// Runs the list option.
// Lists the whole config file in a pretty printed style.
func (l *ListCmd) Run(ctx *Context) error {
	return ctx.Configuration.PrettyPrint()
}

package options

func (u *UseCmd) Run(ctx *Context) error {
  if len(u.Server) > 0 {
    ctx.ConfigFile.CurrentServer = u.Server
    err := ctx.ConfigFile.Save()
    if err != nil {
      return err
    }
  }
  if len(u.Role) > 0 {
    ctx.ConfigFile.CurrentRole = u.Role
    err := ctx.ConfigFile.Save()
    if err != nil {
      return err
    }
  }
  return nil
}

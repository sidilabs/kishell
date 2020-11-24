package options

import (
  "fmt"
)

func (u *UseCmd) Run(ctx *Context) error {
  if len(u.Server) > 0 {
    _, ok := ctx.ConfigFile.Servers[u.Server]
    if !ok {
      return fmt.Errorf("server '%s' is not a valid option", u.Server)
    }
    ctx.ConfigFile.CurrentServer = u.Server
    err := ctx.ConfigFile.Save()
    if err != nil {
      return err
    }
  }
  if len(u.Role) > 0 {
    _, ok := ctx.ConfigFile.Roles[u.Role]
    if !ok {
      return fmt.Errorf("role '%s' is not a valid option", u.Role)
    }
    ctx.ConfigFile.CurrentRole = u.Role
    err := ctx.ConfigFile.Save()
    if err != nil {
      return err
    }
  }
  return nil
}

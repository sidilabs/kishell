package options

import (
	"errors"
	"fmt"
)

func (u *UseCmd) Run(ctx *Context) error {
	if len(u.Server) > 0 {
		_, ok := ctx.Configuration.FindServer(u.Server)
		if !ok {
			return fmt.Errorf("server '%s' is not a valid option", u.Server)
		}
		ctx.Configuration.SetServer(u.Server)
		return ctx.Configuration.Save()
	} else if len(u.Role) > 0 {
		_, ok := ctx.Configuration.FindRole(u.Role)
		if !ok {
			return fmt.Errorf("role '%s' is not a valid option", u.Role)
		}
		ctx.Configuration.SetRole(u.Role)
		return ctx.Configuration.Save()
	}
	return errors.New("missing parameter. One of the following is expected: --server | --role")
}

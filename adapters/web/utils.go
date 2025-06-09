package web

import "github.com/labstack/echo/v4"

func GetPayload(ctx echo.Context, result interface{}) error {
	if err := ctx.Bind(result); err != nil {
		return err
	}
	return nil
}

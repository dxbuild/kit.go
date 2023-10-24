package kit

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HttpBind[IN, OUT any](e *echo.Echo, path string, fn func(in IN, out *OUT) error) {
	e.POST(path, func(ctx echo.Context) error {
		return HttpHandle(ctx, fn)
	})
}

func HttpHandle[IN, OUT any](ctx echo.Context, fn func(in IN, out *OUT) error) error {
	var in IN
	if err := ctx.Bind(&in); err != nil {
		return err
	}
	var out OUT
	err := fn(in, &out)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, out)
}

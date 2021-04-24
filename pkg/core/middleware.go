package core

import echo "github.com/labstack/echo/v4"

type Middleware interface {
	Process(echo.Context) error
}

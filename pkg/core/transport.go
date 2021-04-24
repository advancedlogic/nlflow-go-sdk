package core

import (
	echo "github.com/labstack/echo/v4"
)

type TransportIn interface {
	Process(echo.Context) (Model, error)
}

type TransportOut interface {
	Process(echo.Context, Model) error
}

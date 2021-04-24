package core

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

type defaultTransportIn struct {
}

func (*defaultTransportIn) Process(c echo.Context) (Model, error) {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return nil, err
	}
	defer c.Request().Body.Close()
	var input Model
	err = json.Unmarshal(body, &input)
	if err != nil {
		return nil, err
	}
	return input, nil
}

type defaultTransportOut struct {
}

func (*defaultTransportOut) Process(c echo.Context, outputModel Model) error {
	// convert business logic output to json and send back
	c.JSON(http.StatusOK, outputModel)
	return nil
}

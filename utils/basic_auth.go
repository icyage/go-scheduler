package utils

import "github.com/labstack/echo/v4"

func InternalUser(username, password string, c echo.Context) (bool, error) {
	if username == "internal" && password == "9R33k2xKlu1u9DPDJbR3dt3IFp1JnYpwBJRRG2axUNE9C3Wkgkk0mNrph853YC8H" {
		return true, nil
	}
	return false, nil
}

package handlers

import (
	"github.com/obsurvive/voyager/log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/obsurvive/voyager/utils"
)

type Status struct {
	DB    bool `json:"db"`
	Queue bool `json:"queue"`
	Cache bool `json:"cache"`
}

func Healthz(c echo.Context) error {

	status := &Status{}
	status.DB = true
	status.Cache = true
	responseCode := http.StatusOK

	err, db := utils.GetDB()
	if err != nil {
		status.DB = false
		responseCode = http.StatusInternalServerError
		log.Errorf("DB Connection Error ", err)
	}

	err = db.Ping()
	if err != nil {
		status.DB = false
		responseCode = http.StatusInternalServerError
		log.Errorf("Ping DB Error ", err)
	}

	cache := utils.GetCache()
	_, err = cache.Ping().Result()
	if err != nil {
		status.Cache = false
		responseCode = http.StatusInternalServerError
		log.Errorf("Ping Cache Error ", err)
	}

	return c.JSON(responseCode, status)
}

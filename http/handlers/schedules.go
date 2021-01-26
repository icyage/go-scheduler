package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/obsurvive/voyager/scheduler"
	"net/http"
	"time"
)

type ScheduleCreateRequest struct {
	UUID     string `json:"uuid" form:"uuid" query:"uuid"`
	Interval uint32 `json:"interval" form:"interval" query:"interval"`
}

func ScheduleCreate(c echo.Context) (err error) {

	scheduleRequest := new(ScheduleCreateRequest)
	if err = c.Bind(scheduleRequest); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	interval := time.Duration(scheduleRequest.Interval) * time.Second

	jm := scheduler.GetJobManager()
	job, err := jm.Schedule(scheduleRequest.UUID, &interval)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, job)
}

func ScheduleGet(c echo.Context) (err error) {

	id := c.Param("id")

	jm := scheduler.GetJobManager()
	job, err := jm.Get(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, job)
}

func ScheduleTrigger(c echo.Context) (err error) {

	id := c.Param("id")

	jm := scheduler.GetJobManager()

	job, err := jm.Trigger(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, job)

}

func ScheduleDelete(c echo.Context) (err error) {

	id := c.Param("id")

	jm := scheduler.GetJobManager()
	err = jm.Cancel(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "True")
}

func ScheduleStatus(c echo.Context) (err error) {

	jm := scheduler.GetJobManager()

	count, err := jm.Total()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	m := map[string]interface{}{
		"running_jobs": jm.Running(),
		"total_jobs":   count,
	}
	return c.JSON(http.StatusOK, m)
}

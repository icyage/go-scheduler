package scheduler

import (
	"github.com/obsurvive/voyager/config"
)

var (
	cfg config.Provider
)

func init() {
	cfg = config.Config()
}

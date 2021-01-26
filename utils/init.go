package utils

import (
	"github.com/cenkalti/backoff"
	"github.com/obsurvive/voyager/config"
)

var (
	cfg config.Provider
	bo  backoff.BackOff
)

func init() {
	cfg = config.Config()
	bo = backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 4)
}
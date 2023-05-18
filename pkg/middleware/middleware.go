package middleware

import (
	"fmt"
	"tacy/pkg/botlogger"
	"time"
)

func ConnectWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	logger := botlogger.GetLogger()
	for i := 0; i < attempts; i++ {
		err = fn()
		if err != nil {
			logger.Warn().Err(err).Msg(fmt.Sprintf("Try num: %d, connection failed", i))
			continue
		}
	}
	return err
}

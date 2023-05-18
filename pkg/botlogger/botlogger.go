package botlogger

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var once sync.Once
var log zerolog.Logger

func GetLogger() zerolog.Logger {
	once.Do(func() {
		c := color.New(color.FgRed)
		log = zerolog.New(
			zerolog.ConsoleWriter{
				Out:        os.Stderr,
				NoColor:    false,
				TimeFormat: time.RFC822,
				FormatCaller: func(i interface{}) string {
					return "|" + filepath.Base(fmt.Sprintf("%s|", i))
				},
				FormatErrFieldName: func(i interface{}) string {
					return c.Sprintf("|" + strings.ToUpper(fmt.Sprintf("[%s] -> ", i)))
				},
			}).Level(zerolog.InfoLevel).With().Timestamp().Caller().Logger()
	})
	return log
}

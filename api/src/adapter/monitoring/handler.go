package monitoring

import (
	"log"
	"ravxcheckout/src/adapter/config"
	"time"

	"github.com/getsentry/sentry-go"
)

func StartMonitoring() {
	cfg := config.GetConfig()

	err := sentry.Init(sentry.ClientOptions{
		Dsn:         cfg.SentryDNS,
		Environment: cfg.Environment,
	})

	if err != nil {
		log.Fatalf("sentry.Init failed: ERROR %s", err.Error())
	}

}

func CaptureError(err error) {
	sentry.CaptureException(err)
}

func Flush() {
	sentry.Flush(10 * time.Second)
}


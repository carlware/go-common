package sentry

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/carlware/go-common/log"

	"github.com/getsentry/raven-go"
)

var (
	sentryDSN   string
	debugMode   bool
	environment string
)

func init() {
	Setup("", "local", true)
}

// CaptureRuntime ...
func CaptureRuntime(tags map[string]string) {
	err := recover()
	client, _ := raven.NewWithTags(sentryDSN, nil)
	client.SetEnvironment(environment)

	switch rval := err.(type) {
	case nil:
		return
	case error:
		log.Sugar().Infof("Recover from panic (%s) %s", rval, debug.Stack())
		if !debugMode {
			ex := raven.NewException(rval, raven.NewStacktrace(2, 3, client.IncludePaths()))
			packet := raven.NewPacket(rval.Error(), ex)
			rid, e := client.Capture(packet, tags)
			<-e
			log.Sugar().Infof("Sending error to sentry with id: %s", rid)

		}
	default:
		log.Sugar().Infof("Recover from panic (%s) %s", rval, debug.Stack())
		if !debugMode {
			rvalStr := fmt.Sprint(rval)
			ex := raven.NewException(errors.New(rvalStr), raven.NewStacktrace(2, 3, client.IncludePaths()))
			packet := raven.NewPacket(rvalStr, ex)
			rid, e := client.Capture(packet, tags)
			<-e
			log.Sugar().Infof("Sending error to sentry with id: %s", rid)
		}
	}
}

// Setup ...
func Setup(dsn, env string, debug bool) {
	sentryDSN = dsn
	debugMode = debug
	environment = env
}

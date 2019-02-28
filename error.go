package main

import (
	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
	"log"
	"time"
)

// check every fucking err
func hardFail(e error) {
	defer func() {
		if r := recover(); r != nil {
			exception := appinsights.NewExceptionTelemetry(r)
			exception.SeverityLevel = appinsights.Critical
			exception.Frames = appinsights.GetCallstack(0)
			if conf.Telemetry {
				client.Track(exception)
			}
			panic(r)
		}
	}()

	if e != nil {
		panic(e)
	}
}

// check every fucking err and forget them
func softFail(e error) {
	if e != nil {
		exception := appinsights.NewExceptionTelemetry(e)
		exception.SeverityLevel = appinsights.Warning
		exception.Frames = appinsights.GetCallstack(0)
		if conf.Telemetry {
			client.Track(exception)
		}
		log.Printf("[ERROR] %s", e)
	}
}

func flushTelemetry() {
	select {
	case <-client.Channel().Close(10 * time.Second):
		// Ten second timeout for retries.

		// If we got here, then all telemetry was submitted
		// successfully, and we can proceed to exiting.
	case <-time.After(30 * time.Second):
		// Thirty second absolute timeout.  This covers any
		// previous telemetry submission that may not have
		// completed before Close was called.

		// There are a number of reasons we could have
		// reached here.  We gave it a go, but telemetry
		// submission failed somewhere.  Perhaps old events
		// were still retrying, or perhaps we're throttled.
		// Either way, we don't want to wait around for it
		// to complete, so let's just exit.
	}
}

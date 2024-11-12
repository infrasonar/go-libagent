package libagent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type checkFun func(*Check) (map[string][]map[string]any, error)

type Check struct {
	Key             string
	Collector       *Collector
	Asset           *Asset
	IntervalEnv     string
	DefaultInterval int
	NoCount         bool
	SetTimestamp    bool
	Fn              checkFun
	//Interval must not be configured on init but is calculated after calling Plan(..)
	Interval int
	//Data is a placeholder for custom data or simple `nil` when unused
	Data any
}

func (check *Check) Plan(quit chan bool) {
	s := os.Getenv(check.IntervalEnv)
	if check.Interval > 0 {
		log.Fatal("Interval must not be configured; Instead use IntervalEnv;")
	}
	if s == "" {
		if check.DefaultInterval > 0 {
			check.Interval = check.DefaultInterval
		} else {
			check.Interval = 300
		}
	} else {
		checkInterval, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
		check.Interval = checkInterval
	}

	if check.Interval < 0 {
		log.Fatal("Error: Invalid interval time")
	} else if check.Interval == 0 {
		log.Printf("Warning: %s is disabled (%s=0)\n", check.Key, check.IntervalEnv)
		return // Do not plan this check
	} else if check.Interval < 60 {
		log.Printf("Warning: %s should be at least one minute (60 seconds)\n", check.IntervalEnv)
		// Run the check immediatly as a check.Interval < 60 is only for testing
		check.run()
	} else {
		// We should initially wait for at least a minute and add a little random
		// to avoid different checks to run on the same time
		initWait := randInt(60, 120)
		timer := time.NewTimer(time.Duration(initWait) * time.Second)

		log.Printf("Scheduled: %s: %d / Inital wait: %d\n", check.IntervalEnv, check.Interval, initWait)

		select {
		case <-timer.C:
			check.run()
			break
		case <-quit:
			timer.Stop()
			return
		}
	}

	ticker := time.NewTicker(time.Duration(check.Interval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				check.run()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (check *Check) run() {
	start := time.Now()
	result, err := check.Fn(check)
	runtime := time.Since(start)
	check.handleResult(runtime, result, err)
}

func (check *Check) handleResult(runtime time.Duration, result map[string][]map[string]any, err error) {
	body := map[string]any{
		"runtime": runtime.Seconds(),
		"version": check.Collector.Version,
	}
	if result != nil {
		body["data"] = result
	}
	if check.SetTimestamp {
		body["timestamp"] = time.Now().Unix()
	}
	if check.NoCount {
		body["no_count"] = true
	}
	if err != nil {
		var severity string
		re, ok := err.(*CheckError)
		if ok {
			severity = string(re.Sev)
		} else {
			severity = string(Medium)
		}
		body["error"] = map[string]string{
			"severity": severity,
			"message":  err.Error(),
		}
	}

	data, err := json.Marshal(body)

	if err == nil {
		h := GetHelper()
		uri := fmt.Sprintf("/asset/%d/collector/%s/check/%s", check.Asset.Id, check.Collector.Key, check.Key)
		reader := bytes.NewReader(data)
		err := h.PostRaw(uri, reader, 2)
		if err != nil {
			log.Printf("Error while sending check data: %s\n", err)
		}
	} else if result != nil {
		check.handleResult(runtime, nil, &CheckError{Sev: High, Err: err})
	} else {
		log.Printf("Unexcpected JSON pack error: %s", err)
	}
}

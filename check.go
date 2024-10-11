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

type checkFun func() (map[string][]map[string]any, error)

type Check struct {
	Key          string
	Collector    *Collector
	Asset        *Asset
	IntervalEnv  string
	NoCount      bool
	SetTimestamp bool
	Fn           checkFun
}

func (check *Check) Plan(quit chan bool) {
	s := os.Getenv(check.IntervalEnv)
	if s == "" {
		s = "300"
	}
	checkInterval, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	if checkInterval < 1 {
		log.Fatal("Error: Invalid interval time")
	} else if checkInterval < 60 {
		log.Printf("Warning: %s should be at least one minute (60 seconds)\n", check.IntervalEnv)
		// Run the check immediatly as a checkInterval < 60 is only for testing
		check.run()

	} else {
		// We should initially wait for at least a minute and add a little random
		// to avoid different checks to run on the same time
		initWait := randInt(60, 120)
		timer := time.NewTimer(time.Duration(initWait) * time.Second)

		log.Printf("Scheduled: %s: %d / Inital wait: %d\n", check.IntervalEnv, checkInterval, initWait)

		select {
		case <-timer.C:
			check.run()
			break
		case <-quit:
			timer.Stop()
			return
		}
	}

	ticker := time.NewTicker(time.Duration(checkInterval) * time.Second)
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
	result, err := check.Fn()
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
		err := h.PostRaw(uri, reader)
		if err != nil {
			log.Printf("Error while sending check data: %s\n", err)
		}
	} else if result != nil {
		check.handleResult(runtime, nil, &CheckError{Sev: High, Err: err})
	} else {
		log.Printf("Unexcpected JSON pack error: %s", err)
	}
}

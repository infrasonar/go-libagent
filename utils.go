package libagent

import (
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/exp/rand"
)

// RandInit initializes random
func RandInit() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

// SigHandler is required for handling signals
func SigHandler(cs chan bool) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	signal := <-c
	log.Printf("Signal: %s\n", signal.String())

	switch signal {
	case syscall.SIGINT, syscall.SIGTERM:
		cs <- true
	case syscall.SIGHUP:
		cs <- false
	}
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func fqdn() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	addrs, err := net.LookupIP(hostname)
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			ip, err := ipv4.MarshalText()
			if err != nil {
				return "", err
			}
			hosts, err := net.LookupAddr(string(ip))
			if err != nil {
				return "", err
			}
			if len(hosts) == 0 {
				return "", errors.New("No hosts found for ip: " + string(ip))
			}
			return strings.TrimSuffix(hosts[0], "."), nil
		}
	}
	return "", errors.New("no FQDN could be found")
}

func getConfigPath() (string, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		// Fallback to `STORAGE_PATH`
		configPath = os.Getenv("STORAGE_PATH")
	}
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			configPath = path.Join(homeDir, ".infrasonar")
		} else {
			configPath = "/etc/infrasonar"
		}

	}
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(configPath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	return configPath, nil
}

type IFloat64 float64

func (f IFloat64) MarshalJSON() ([]byte, error) {
	if float64(f) == float64(int(f)) {
		return []byte(strconv.FormatFloat(float64(f), 'f', 1, 64)), nil
	}
	return []byte(strconv.FormatFloat(float64(f), 'f', -1, 64)), nil
}

type IFloat32 float32

func (f IFloat32) MarshalJSON() ([]byte, error) {
	if float32(f) == float32(int(f)) {
		return []byte(strconv.FormatFloat(float64(f), 'f', 1, 32)), nil
	}
	return []byte(strconv.FormatFloat(float64(f), 'f', -1, 32)), nil
}

package libagent

import (
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"golang.org/x/exp/rand"
)

func randInit() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func sigHandler(cs chan bool) {
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

func getStoragePath() (string, error) {
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", errors.New("failed to find the use HOME path, try to set `STORAGE_PATH` to work arround this issue")
		}
		storagePath = path.Join(homeDir, ".infrasonar")
	}
	_, err := os.Stat(storagePath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(storagePath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	return storagePath, nil
}

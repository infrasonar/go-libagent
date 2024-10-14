package libagent

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type helper struct {
	authorization string
	apiUri        string
}

var helperInstance *helper

func GetHelper() *helper {
	if helperInstance == nil {
		skipVerify := os.Getenv("SKIP_VERIFY")
		if skipVerify != "" {
			http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}

		token := os.Getenv("TOKEN")
		if token == "" {
			log.Fatal("Missing TOKEN environment variable")
		}
		apiUri := os.Getenv("API_URI")
		if apiUri == "" {
			apiUri = "https://api.infrasonar.com"
		}
		if !strings.HasPrefix(apiUri, "https://") {
			log.Fatal("Invalid API_URI environment variable")
		}
		helperInstance = &helper{
			authorization: "Bearer " + token,
			apiUri:        apiUri,
		}
	}
	return helperInstance
}

func handleResp(resp *http.Response, t any) error {
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
	if t == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(&t)
}

func doReq(req *http.Request, t any) error {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	return handleResp(resp, t)
}

func (h *helper) Get(uri string, t any) error {
	req, err := http.NewRequest("GET", h.apiUri+uri, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", h.authorization)

	return doReq(req, t)
}

func (h *helper) Post(uri string, t any, v any) error {
	b := new(bytes.Buffer)
	if v != nil {
		err := json.NewEncoder(b).Encode(v)
		if err != nil {
			return err
		}
	}
	req, err := http.NewRequest("POST", h.apiUri+uri, b)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", h.authorization)
	req.Header.Set("Content-Type", "application/json")

	return doReq(req, t)
}

func (h *helper) PostRaw(uri string, body io.Reader) error {
	req, err := http.NewRequest("POST", h.apiUri+uri, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", h.authorization)
	req.Header.Set("Content-Type", "application/json")

	return doReq(req, nil)
}

func (h *helper) Patch(uri string, t any, v any) error {
	b := new(bytes.Buffer)
	if v != nil {
		err := json.NewEncoder(b).Encode(v)
		if err != nil {
			return err
		}
	}
	req, err := http.NewRequest("PATCH", h.apiUri+uri, b)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", h.authorization)
	req.Header.Set("Content-Type", "application/json")

	return doReq(req, t)
}

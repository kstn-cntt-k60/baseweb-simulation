package util

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
)

func getHostname() string {
	hostname := os.Getenv("BASEWEB_HOSTNAME")
	if hostname == "" {
		hostname = "localhost:8080"
	}
	return hostname
}

func Url(pathname string) string {
	return fmt.Sprintf("http://%s%s", getHostname(), pathname)
}

func Login() (string, error) {
	basicAuth := "admin:admin"
	basicAuth = base64.StdEncoding.EncodeToString([]byte(basicAuth))
	basicAuth = "Basic " + basicAuth

	client := http.Client{}
	body := bytes.NewBuffer([]byte{})
	req, err := http.NewRequest("POST", Url("/api/login"), body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", basicAuth)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", errors.New("not 200")
	}

	return res.Header.Get("X-Auth-Token"), nil
}

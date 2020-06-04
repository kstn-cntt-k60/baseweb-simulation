package util

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
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

func Login(client *http.Client) (string, error) {
	basicAuth := "admin:admin"
	basicAuth = base64.StdEncoding.EncodeToString([]byte(basicAuth))
	basicAuth = "Basic " + basicAuth

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
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", errors.New("not 200")
	}

	return res.Header.Get("X-Auth-Token"), nil
}

func Post(
	client *http.Client,
	token string,
	pathname string,
	data []byte,
) (*http.Response, error) {

	body := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", Url(pathname), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Auth-Token", token)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return res, err
	}

	if res.StatusCode != 200 {
		return res, errors.New(
			fmt.Sprintf("%s: status code %d", pathname, res.StatusCode))
	}

	return res, nil
}

func Get(
	client *http.Client,
	token string,
	pathname string,
) (*http.Response, error) {
	body := bytes.NewBuffer([]byte{})
	req, err := http.NewRequest("GET", Url(pathname), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Auth-Token", token)

	res, err := client.Do(req)
	if err != nil {
		return res, err
	}

	if res.StatusCode != 200 {
		return res, errors.New(
			fmt.Sprintf("%s: status code %d", pathname, res.StatusCode))
	}

	return res, nil
}

func RandomBetween(begin int, end int) int {
	return rand.Intn(end-begin) + begin
}

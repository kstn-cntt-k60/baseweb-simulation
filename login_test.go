package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"net/http"
	"testing"
)

func TestLogin(t *testing.T) {
	basicAuth := "admin:admin"
	basicAuth = base64.StdEncoding.EncodeToString([]byte(basicAuth))
	basicAuth = "Basic " + basicAuth

	client := http.Client{}
	body := bytes.NewBuffer([]byte{})
	req, err := http.NewRequest("POST", "http://localhost:8080/api/login", body)
	if err != nil {
		t.Log(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", basicAuth)

	res, err := client.Do(req)
	if err != nil {
		t.Log(err)
		return
	}

	if res.StatusCode != 200 {
		t.Log("not 200")
		return
	}

	t.Log(res.Header.Get("X-Auth-Token"))
}

func Login() (string, error) {
	basicAuth := "admin:admin"
	basicAuth = base64.StdEncoding.EncodeToString([]byte(basicAuth))
	basicAuth = "Basic " + basicAuth

	client := http.Client{}
	body := bytes.NewBuffer([]byte{})
	req, err := http.NewRequest("POST", "http://localhost:8080/api/login", body)
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

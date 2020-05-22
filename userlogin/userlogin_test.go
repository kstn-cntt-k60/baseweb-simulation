package userlogin

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"baseweb-simulation/util"

	"github.com/bxcodec/faker/v3"
)

type UserLogin struct {
	UserName string `faker:"username" json:"username"`
	Password string `faker:"-" json:"password"`
	PersonId string `faker:"-" json:"personId"`
}

var ListPersonId = []string{"7036e7c2-9c01-11ea-89ab-14dda9bea6d7", "703699ea-9c01-11ea-89ab-14dda9bea6d7",
	"70370ece-9c01-11ea-89ab-14dda9bea6d7", "738ffe46-9c01-11ea-89ab-14dda9bea6d7", "738f3b0e-9c01-11ea-89ab-14dda9bea6d7",
	"739c7f8d-9c01-11ea-89ab-14dda9bea6d7", "739c0ab9-9c01-11ea-89ab-14dda9bea6d7", "73a24b1e-9c01-11ea-89ab-14dda9bea6d7"}

func RandomPersonId() string {
	return ListPersonId[rand.Intn(len(ListPersonId))]
}

func GenerateUserLogin(userlogin *UserLogin) error {
	err := faker.FakeData(userlogin)
	if err != nil {
		return err
	}
	userlogin.Password = "admin"
	userlogin.PersonId = RandomPersonId()
	return err
}

func AddUserLogin(token string) error {
	userlogin := UserLogin{}
	err := GenerateUserLogin(&userlogin)
	if err != nil {
		return err
	}

	data, err := json.Marshal(userlogin)
	if err != nil {
		return err
	}

	client := http.Client{}
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", "http://localhost:8080/api/account/add-user-login", body)
	if err != nil {
		return err
	}

	req.Header.Set("X-Auth-Token", token)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("add user login not 200")
	}

	log.Println("userlogin: ", userlogin)

	return nil
}

type PostResult struct {
	err      error
	duration time.Duration
}

func RunAddUserLogin(token string, ch chan<- PostResult) {
	begin := time.Now()
	err := AddUserLogin(token)
	end := time.Now()
	duration := end.Sub(begin)

	if err != nil {
		ch <- PostResult{err: err}
	} else {
		ch <- PostResult{err: nil, duration: duration}
	}
}

func TestAddUserLogin(t *testing.T) {
	loopCount := 20

	token, err := util.Login()
	if err != nil {
		t.Error(err)
	}

	ch := make(chan PostResult)
	defer close(ch)

	for i := 0; i < loopCount; i++ {
		go RunAddUserLogin(token, ch)
	}

	var sum int64 = 0
	for i := 0; i < loopCount; i++ {
		result := <-ch
		sum += result.duration.Microseconds()
		t.Log(result.err, result.duration.Microseconds())
	}
	t.Log("Avg: ", (sum / int64(loopCount)))
}

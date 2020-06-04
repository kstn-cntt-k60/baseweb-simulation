package customer

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"baseweb-simulation/util"

	"github.com/bxcodec/faker/v3"
)

type Customer struct {
	PartyTypeId  int    `faker:"-" json:"partyTypeId"`
	Description  string `faker:"sentence" json:"description"`
	CustomerName string `faker:"name" json:"customerName"`
}

var List = []int{1, 2}

func RandomGenderId() int {
	return List[rand.Intn(len(List))]
}

func GenerateCustomer(customer *Customer) error {
	err := faker.FakeData(customer)
	if err != nil {
		return err
	}
	customer.PartyTypeId = 2
	customer.Description = "Customer " + customer.Description
	return err
}

func AddCustomer(token string) error {
	customer := Customer{}
	err := GenerateCustomer(&customer)
	if err != nil {
		return err
	}

	data, err := json.Marshal(customer)
	if err != nil {
		return err
	}

	client := http.Client{}
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", "http://localhost:8080/api/account/add-party", body)
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
		return errors.New("add customer not 200")
	}

	return nil
}

type PostResult struct {
	err      error
	duration time.Duration
}

func RunAddCustomer(token string, ch chan<- PostResult) {
	begin := time.Now()
	err := AddCustomer(token)
	end := time.Now()
	duration := end.Sub(begin)
	if err != nil {
		ch <- PostResult{err: err}
	} else {
		ch <- PostResult{err: nil, duration: duration}
	}
}

func TestAddCustomer(t *testing.T) {
	loopCount := 10

	client := &http.Client{}
	defer client.CloseIdleConnections()

	token, err := util.Login(client)
	if err != nil {
		t.Error(err)
	}

	ch := make(chan PostResult)
	defer close(ch)

	for i := 0; i < loopCount; i++ {
		go RunAddCustomer(token, ch)
	}

	var sum int64 = 0
	for i := 0; i < loopCount; i++ {
		result := <-ch
		sum += result.duration.Microseconds()
		t.Log(result.err, result.duration.Microseconds())
	}
	t.Log("Avg: ", (sum / int64(loopCount)))
}

package person

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

type Person struct {
	FirstName   string `faker:"first_name" json:"firstName"`
	MiddleName  string `faker:"first_name_male" json:"middleName"`
	LastName    string `faker:"last_name" json:"lastName"`
	PartyTypeId int    `faker:"-" json:"partyTypeId"`
	Description string `faker:"sentence" json:"description"`
	GenderId    int    `faker:"-" json:"genderId"`
	BirthDate   string `faker:"timestamp" json:"birthDate"`
}

var List = []int{1, 2}

func RandomGenderId() int {
	return List[rand.Intn(len(List))]
}

func GeneratePerson(person *Person) error {
	err := faker.FakeData(person)
	if err != nil {
		return err
	}
	person.PartyTypeId = 1
	person.GenderId = RandomGenderId()
	person.Description = "Person " + person.Description
	return err
}

func TestPerson(t *testing.T) {
	rand.Seed(100)
	for i := 0; i < 10; i++ {
		person := Person{}
		err := GeneratePerson(&person)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(person)
	}
}

func AddPerson(token string) error {
	person := Person{}
	err := GeneratePerson(&person)
	if err != nil {
		return err
	}

	data, err := json.Marshal(person)
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
		return errors.New("add person not 200")
	}

	log.Println("person: ", person.FirstName, person.MiddleName, person.LastName)

	return nil
}

type PostResult struct {
	err      error
	duration time.Duration
}

func RunAddPerson(token string, ch chan<- PostResult) {
	begin := time.Now()
	err := AddPerson(token)
	end := time.Now()
	duration := end.Sub(begin)

	if err != nil {
		ch <- PostResult{err: err}
	} else {
		ch <- PostResult{err: nil, duration: duration}
	}
}

func TestAddPerson(t *testing.T) {
	loopCount := 5

	token, err := util.Login()
	if err != nil {
		t.Error(err)
	}

	ch := make(chan PostResult)
	defer close(ch)

	for i := 0; i < loopCount; i++ {
		go RunAddPerson(token, ch)
	}

	var sum int64 = 0
	for i := 0; i < loopCount; i++ {
		result := <-ch
		sum += result.duration.Microseconds()
		t.Log(result.err, result.duration.Microseconds())
	}
	t.Log("Avg: ", (sum / int64(loopCount)))
}

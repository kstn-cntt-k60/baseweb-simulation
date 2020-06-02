package person

import (
	"baseweb-simulation/util"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"time"

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

var GenderIdList = []int{1, 2}

func RandomGenderId() int {
	return GenderIdList[rand.Intn(len(GenderIdList))]
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
	req, err := http.NewRequest("POST", util.Url("/api/account/add-party"), body)
	if err != nil {
		return err
	}

	req.Header.Set("X-Auth-Token", token)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New("add person not 200")
	}

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

func AddPersonBenchmark(loop int) {
	token, err := util.Login()
	if err != nil {
		log.Panic(err)
	}

	ch := make(chan PostResult)
	defer close(ch)

	for i := 0; i < loop; i++ {
		go RunAddPerson(token, ch)
	}

	var errorCount int64 = 0
	for i := 0; i < loop; i++ {
		result := <-ch
		if result.err != nil {
			errorCount++
		} else {
			log.Println(result.duration.Microseconds())
		}
	}

	log.Println("Error Count", errorCount)
}

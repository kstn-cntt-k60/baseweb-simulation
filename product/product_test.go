package product

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/shopspring/decimal"
)

type Product struct {
	Name        string          `faker:"name" json:"name"`
	Description string          `faker:"sentence" json:"description"`
	Weight      decimal.Decimal `faker:"-" json:"weight"`
	WeightUomId string          `faker:"-" json:"weightUomId"`
	UnitUomId   string          `faker:"-" json:"unitUomId"`
}

var UnitUomIdList = []string{"package", "box", "bottle"}
var WeightUomIdList = []string{"kg", "g", "mg"}

func RandomUnitUom() string {
	return UnitUomIdList[rand.Intn(len(UnitUomIdList))]
}

func RandomWeightUomId() string {
	return WeightUomIdList[rand.Intn(len(WeightUomIdList))]
}

func RandomWeight() (decimal.Decimal, error) {
	type Weight struct {
		Value int64 `faker:"boundary_start=1, boundary_end=1000"`
		Exp   int32 `faker:"boundary_start=-2, boundary_end=0"`
	}
	w := Weight{}
	err := faker.FakeData(&w)
	if err != nil {
		return decimal.Zero, err
	}
	return decimal.New(w.Value, w.Exp), err
}

func GenerateProduct(product *Product) error {
	err := faker.FakeData(product)
	if err != nil {
		return err
	}
	product.UnitUomId = RandomUnitUom()
	product.WeightUomId = RandomWeightUomId()
	product.Weight, err = RandomWeight()
	return err
}

func TestProduct(t *testing.T) {
	rand.Seed(200)

	for i := 0; i < 10; i++ {
		p := Product{}
		err := GenerateProduct(&p)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(p)
	}
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
	req.Header.Set("Authorization", basicAuth)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", errors.New("can't login")
	}

	return res.Header.Get("X-Auth-Token"), nil
}

func AddProduct(token string) error {
	product := Product{}
	err := GenerateProduct(&product)
	if err != nil {
		return err
	}
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	client := http.Client{}
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST",
		"http://localhost:8080/api/product/add-product", body)
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
		return errors.New("can't add product")
	}

	return nil
}

type PostResult struct {
	err      error
	duration time.Duration
}

func RunAddProduct(token string, ch chan<- PostResult) {
	begin := time.Now()
	err := AddProduct(token)
	end := time.Now()
	duration := end.Sub(begin)

	if err != nil {
		ch <- PostResult{err: err}
	} else {
		ch <- PostResult{err: nil, duration: duration}
	}
}

func TestClient(t *testing.T) {
	loopCount := 20

	token, err := Login()
	if err != nil {
		t.Log(err)
		return
	}

	ch := make(chan PostResult)
	defer close(ch)

	begin := time.Now()

	for i := 0; i < loopCount; i++ {
		go RunAddProduct(token, ch)
	}

	var sum int64 = 0
	for i := 0; i < loopCount; i++ {
		result := <-ch
		sum += result.duration.Microseconds()
		t.Log(result.err, result.duration.Microseconds())
	}
	t.Log("Avg: ", sum/int64(loopCount))
	t.Log("Time: ", time.Now().Sub(begin).Microseconds())
}

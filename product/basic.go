package product

import (
	"baseweb-simulation/util"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
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

func RandomWeight() decimal.Decimal {
	value := util.RandomBetween(1, 1000)
	exp := util.RandomBetween(-2, 0)
	return decimal.New(int64(value), int32(exp))
}

func GenerateProduct(product *Product) error {
	err := faker.FakeData(product)
	if err != nil {
		return err
	}
	product.UnitUomId = RandomUnitUom()
	product.WeightUomId = RandomWeightUomId()
	product.Weight = RandomWeight()
	return nil
}

func AddProduct(client *http.Client, token string) error {
	product := Product{}
	err := GenerateProduct(&product)
	if err != nil {
		return err
	}
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	res, err := util.Post(
		client, token, "/api/product/add-product", data)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

type PostResult struct {
	err      error
	duration time.Duration
}

func RunAddProduct(client *http.Client, token string, ch chan<- PostResult) {
	begin := time.Now()
	err := AddProduct(client, token)
	end := time.Now()
	duration := end.Sub(begin)

	if err != nil {
		ch <- PostResult{err: err}
	} else {
		ch <- PostResult{err: nil, duration: duration}
	}
}

func AddProductBenchmark(loop int) {
	client := &http.Client{}
	defer client.CloseIdleConnections()

	token, err := util.Login(client)
	if err != nil {
		log.Panic(err)
	}

	ch := make(chan PostResult)
	defer close(ch)

	for i := 0; i < loop; i++ {
		go RunAddProduct(client, token, ch)
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

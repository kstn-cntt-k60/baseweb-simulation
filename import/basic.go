package importProduct

import (
	"baseweb-simulation/util"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type InventoryItem struct {
	WarehouseId   uuid.UUID       `json:"warehouseId"`
	ProductId     int             `json:"productId"`
	Quantity      decimal.Decimal `json:"quantity"`
	UnitCost      decimal.Decimal `json:"unitCost"`
	CurrencyUomId string          `json:"currencyUomId"`
}

func GetProducts(client *http.Client, token string) []Product {
	pathname := fmt.Sprintf(
		"/api/product/view-product?page=0&pageSize=%d&sortedBy=createdAt&sortOrder=desc&query=", 100)

	response, err := util.Get(client, token, pathname)
	if err != nil {
		log.Panic(err)
	}
	defer response.Body.Close()

	type Response struct {
		ProductList  []Product `json:"productList"`
		ProductCount int       `json:"productCount"`
	}

	res := Response{}
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		log.Panic(err)
	}

	return res.ProductList
}

func RandomQuantity() decimal.Decimal {
	value := util.RandomBetween(1, 1000)
	return decimal.New(int64(value), 0)
}

func RandomCost() decimal.Decimal {
	value := util.RandomBetween(2, 200)
	return decimal.New(int64(value), 3)
}

func RandomProductId(products []Product) int {
	return products[rand.Intn(len(products))].Id
}

func GenerateInventoryItem(
	warehouseId uuid.UUID,
	products []Product,
	item *InventoryItem,
) error {
	item.Quantity = RandomQuantity()
	item.ProductId = RandomProductId(products)
	item.WarehouseId = warehouseId
	item.UnitCost = RandomCost()
	item.CurrencyUomId = "vnd"

	return nil
}

func AddInventoryItem(
	client *http.Client,
	token string,
	warehouseId uuid.UUID,
	products []Product,
) error {
	item := InventoryItem{}
	err := GenerateInventoryItem(warehouseId, products, &item)
	if err != nil {
		return err
	}

	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	res, err := util.Post(client, token, "/api/import/add-inventory-item", data)
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

func RunAddInventoryItem(
	client *http.Client,
	token string,
	warehouseId uuid.UUID,
	products []Product,
	ch chan<- PostResult,
) {
	begin := time.Now()
	err := AddInventoryItem(client, token, warehouseId, products)
	end := time.Now()
	duration := end.Sub(begin)

	if err != nil {
		ch <- PostResult{err: err}
	} else {
		ch <- PostResult{err: nil, duration: duration}
	}
}

func AddInventoryItemBenchmark(loop int) {
	warehouseId, err := uuid.Parse("28fb8f4a-5a02-11ea-b26e-14dda9bea6d7")
	if err != nil {
		log.Panic(err)
	}

	client := &http.Client{}
	defer client.CloseIdleConnections()

	token, err := util.Login(client)
	if err != nil {
		log.Panic(err)
	}

	products := GetProducts(client, token)

	ch := make(chan PostResult)
	defer close(ch)

	for i := 0; i < loop; i++ {
		go RunAddInventoryItem(client, token, warehouseId, products, ch)
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

func ViewInventoryItem(client *http.Client, token string) {
	page, err := strconv.Atoi(os.Getenv("PAGE"))
	if err != nil {
		page = 0
	}

	pathname := fmt.Sprintf(
		"/api/import/view-inventory-by-warehouse?page=%d&pageSize=10&sortedBy=createdAt&sortOrder=desc&warehouseId=28fb8f4a-5a02-11ea-b26e-14dda9bea6d7",
		page)

	log.Println("PATHNAME", pathname)

	response, err := util.Get(client, token, pathname)
	if err != nil {
		log.Panic(err)
	}
	defer response.Body.Close()
}

func ViewInventoryItemBenchmark() {
	client := &http.Client{}
	defer client.CloseIdleConnections()

	token, err := util.Login(client)
	if err != nil {
		log.Panic(err)
	}

	begin := time.Now()
	ViewInventoryItem(client, token)
	end := time.Now()
	duration := end.Sub(begin)

	log.Println(duration.Microseconds())
}

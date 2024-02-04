package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/nats-io/stan.go"
)

type Order struct {
	OrderUID          string   `json:"order_uid"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerID        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service"`
	ShardKey          string   `json:"shardkey"`
	SmID              int      `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func main() {
	gofakeit.Seed(0)

	sc, err := stan.Connect("wbCluster", "client-id", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	for {
		var order Order = Order{
			OrderUID:    gofakeit.UUID(),
			TrackNumber: gofakeit.UUID(),
			Entry:       gofakeit.Word(),
			Delivery: Delivery{
				Name:    gofakeit.Name(),
				Phone:   gofakeit.Phone(),
				Zip:     gofakeit.Zip(),
				City:    gofakeit.City(),
				Address: gofakeit.Address().Address,
				Region:  gofakeit.State(),
				Email:   gofakeit.Email(),
			},
			Payment: Payment{
				Transaction:  gofakeit.UUID(),
				RequestID:    gofakeit.UUID(),
				Currency:     gofakeit.CurrencyShort(),
				Provider:     gofakeit.CreditCardType(),
				Amount:       gofakeit.Number(100, 10000),
				PaymentDt:    gofakeit.Date().Day(),
				Bank:         gofakeit.NewCrypto().Car().Brand,
				DeliveryCost: gofakeit.Number(10, 500),
				GoodsTotal:   gofakeit.Number(500, 10000),
				CustomFee:    gofakeit.Number(0, 500),
			},
			Items:             generateRandomItems(rand.Intn(100) + 1),
			Locale:            gofakeit.Date().Local().GoString(),
			InternalSignature: gofakeit.HackerPhrase(),
			CustomerID:        gofakeit.UUID(),
			DeliveryService:   gofakeit.Company(),
			ShardKey:          gofakeit.Word(),
			SmID:              rand.Intn(1000),
			DateCreated:       gofakeit.Date().Format("2006-01-02"),
			OofShard:          gofakeit.Word(),
		}

		jsondata, _ := json.MarshalIndent(order, "", " ")
		fmt.Println("go go)")

		err = sc.Publish("wbCluster", jsondata)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(10 * time.Second)
	}

}
func generateRandomItems(numItems int) []Item {
	items := make([]Item, numItems)
	for i := range items {
		items[i] = Item{
			ChrtID:      rand.Intn(100000),
			TrackNumber: gofakeit.UUID(),
			Price:       gofakeit.Number(100, 10000),
			Rid:         gofakeit.UUID(),
			Name:        gofakeit.ProductName(),
			Sale:        gofakeit.Number(0, 100),
			Size:        gofakeit.RandomString([]string{"S", "M", "L", "XL", "XXL"}),
			TotalPrice:  gofakeit.Number(100, 10000),
			NmID:        rand.Intn(100000),
			Brand:       gofakeit.NewCrypto().Car().Brand,
			Status:      gofakeit.Number(1, 5),
		}
	}
	return items
}

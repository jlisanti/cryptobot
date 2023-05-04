package cryptobot

import (
	"fmt"

	"github.com/preichenberger/go-coinbasepro/v2"
	"github.com/shopspring/decimal"
)

type TransactionRequest struct {
	ProductID string
	Side      string
	Size      string
}

func Connect() (client *coinbasepro.Client) {

	client = coinbasepro.NewClient()

	//Sandbox key
	client.UpdateConfig(&coinbasepro.ClientConfig{
		BaseURL:    "https://api-public.sandbox.pro.coinbase.com",
		Key:        "89b1f52167924567c1a41b42a236d8a1",
		Passphrase: "puj31du7a4j",
		Secret:     "RUWPjira048friEd52Z34ptpYdeFnop1PrucxrvGRlZhUtNuM71Iub+HwTu7X2Gg8OjVkuFIW1iPm5C8qzamgw==",
	})

	println("account connection established\n")

	return client
}

func Print(client *coinbasepro.Client) {
	accounts, err := client.GetAccounts()
	if err != nil {
		println(err.Error())
	}

	for _, a := range accounts {
		println(a.Currency, a.Balance)
	}

	var orders []coinbasepro.Order

	fmt.Println("Printing open orders...")
	for _, o := range orders {
		fmt.Println(o.ID, o.CreatedAt.Time())
	}

	fmt.Println("Finished printing open orders")
	transfer := coinbasepro.Transfer{
		Type:   "deposit",
		Amount: "1000.00",
	}

	savedTransfer, err := client.CreateTransfer(&transfer)
	if err != nil {
		println(err.Error())
	}
	println(savedTransfer.Amount)

}

func Order(client *coinbasepro.Client, request TransactionRequest) (orderID string) {

	var orders []coinbasepro.Order
	//cursor := client.ListOrders()

	//fmt.Println("Canceling open orders...")
	//for _, o := range orders {
	//	client.CancelOrder(o.ID)
	//	}

	fmt.Println("Printing open orders...")
	for _, o := range orders {
		fmt.Println(o.ID, o.CreatedAt.Time())
	}

	fmt.Println("Finished printing open orders")

	book, err := client.GetBook("BTC-USD", 1)
	if err != nil {
		println(err.Error())
	}

	lastPrice, err := decimal.NewFromString(book.Bids[0].Price)
	if err != nil {
		println(err.Error())
	}

	fmt.Printf("Current price: $%s \n", lastPrice)

	//order := coinbasepro.Order{
	//	Price:     lastPrice.Add(decimal.NewFromFloat(1.00)).String(),
	//Size:      request.Size,
	//Side:      request.Side,
	//ProductID: request.ProductID,
	//}

	order := coinbasepro.Order{
		Price:     lastPrice.Add(decimal.NewFromFloat(1.00)).String(),
		Size:      "0.00005",
		Side:      "buy",
		ProductID: "BTC-USD",
	}

	savedOrder, err := client.CreateOrder(&order)
	if err != nil {
		println(err.Error())
	}

	return savedOrder.ID
}

package cryptobot

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/preichenberger/go-coinbasepro/v2"
)

func Track(client *coinbasepro.Client) {

	assets := pullAssets(client)
	//for i, a := range assets {
	//println(i, a.BuyDate.Hour(), a.BuyDate.Minute(), a.BuyDate.Second(), a.BuyDate.Day(), a.BuyDate.Month(), a.BuyDate.Year(), a.Currency, a.Quantity)
	//}
	fmt.Print(assets)

}

func pullAssets(client *coinbasepro.Client) (assets []Asset) {

	acounts, err := client.GetAccounts()
	if err != nil {
		println(err.Error())
	}

	var ledgers []coinbasepro.LedgerEntry

	for _, a := range acounts {
		cursor := client.ListAccountLedger(a.ID)
		for cursor.HasMore {
			if err := cursor.NextPage(&ledgers); err != nil {
				println(err.Error())
			}

			for _, e := range ledgers {
				fmt.Println(e.Amount, e.CreatedAt.Time(), e.Type)
				if e.Type == "match" {
					currencies := strings.Split(e.Details.ProductID, "-")
					// Determine if it was a buy or sell
					transferAmount, _ := strconv.ParseFloat(e.Amount, 64)
					if transferAmount > float64(0.0) {
						// record asset
						asset := Asset{
							ID:       e.Details.TradeID,
							Currency: currencies[0],
							Quantity: e.Amount,
							BuyDate:  time.Time(e.CreatedAt.Time()),
							BuyPrice: "",
							Cost:     "",
						}
						assets = append(assets, asset)
					} else {
						for i, asset := range assets {
							if asset.ID == e.Details.TradeID {
								(assets)[i].BuyPrice = e.Amount
							}
						}
					}
				} else if e.Type == "fee" {
					for i, asset := range assets {
						if asset.ID == e.Details.TradeID {
							(assets)[i].Cost = e.Amount
						}
					}
				}
			}
		}
	}
	return assets
}

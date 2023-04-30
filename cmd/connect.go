/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jlisanti/cryptobot/pkg/cryptobot"
	"github.com/manifoldco/promptui"
	"github.com/preichenberger/go-coinbasepro/v2"
	"github.com/spf13/cobra"
)

type promptContent struct {
	errorMsg string
	label    string
}

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to trading accounts",
	Long:  `Connect to trading accounts`,
	Run: func(cmd *cobra.Command, args []string) {
		connect()
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}

func promptGetInput(pc promptContent) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.errorMsg)
		}
		return nil

	}
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }}",
		Valid:   "{{ . | green}}",
		Invalid: "{{ . | red}}",
		Success: "{{ . | bold}}",
	}
	prompt := promptui.Prompt{
		Label:     pc.label,
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Input: %s\n", result)
	return result
}

func connect() {
	client := cryptobot.Connect()

	actionPromptContent := promptContent{
		"Acount connection established",
		"Please select an action",
	}

	action := promptGetSelect(actionPromptContent)

	if action == "print" {
		cryptobot.Print(client)
	} else if action == "buy" {
		buyPrompt(client)
	} else if action == "sell" {
		sellPrompt(client)
	}
}

func promptGetSelect(pc promptContent) string {

	items := []string{"print", "buy", "sell", "track", "simulate"}
	index := -1

	var result string
	var err error

	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label:    pc.label,
			Items:    items,
			AddLabel: "other",
		}

		index, result, err = prompt.Run()

		if index == -1 {
			items = append(items, result)
		}
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	fmt.Printf("Input: %s\n", result)

	return result
}

func buyPrompt(client *coinbasepro.Client) {
	buyRequest := cryptobot.TransactionRequest{}

	bookPromptContent := promptContent{
		"Enter crypto to purchase",
		"Buy crypto: ",
	}
	productID := promptGetInput(bookPromptContent) + "-USD"

	buyRequest.ProductID = strings.ReplaceAll(productID, " ", "")

	buyQuantityPromptContent := promptContent{
		"Enter purchase quantity",
		"Purchase amount: ",
	}
	size := promptGetInput(buyQuantityPromptContent)
	buyRequest.Size = strings.ReplaceAll(size, " ", "")

	buyRequest.Side = "buy"

	cryptobot.Order(client, buyRequest)
}

func sellPrompt(client *coinbasepro.Client) {
	sellRequest := cryptobot.TransactionRequest{}

	bookPromptContent := promptContent{
		"Enter crypto to sell",
		"Sell crypto: ",
	}
	productID := promptGetInput(bookPromptContent) + "-USD"

	sellRequest.ProductID = strings.ReplaceAll(productID, " ", "")

	sellQuantityPromptContent := promptContent{
		"Enter sell quantity",
		"Purchase amount: ",
	}
	size := promptGetInput(sellQuantityPromptContent)
	sellRequest.Size = strings.ReplaceAll(size, " ", "")

	sellRequest.Side = "sell"

	cryptobot.Order(client, sellRequest)
}

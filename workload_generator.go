package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
		panic(err)
	}
}

type Add struct {
	UserId string
	Amount int
}
type Quote struct {
	UserId      string
	StockSymbol string
}

type Default struct {
	UserId      string
	StockSymbol string
	Amount      int
}

type User struct {
	UserId string
}

func add(r []string) {
	toWebServer := Add{}
	toWebServer.UserId = r[1]
	s, err := strconv.ParseFloat(r[2], 64)
	if err != nil {
		//handle error
		return
	}
	toWebServer.Amount = int(s * 10)

	sendToWebServer(toWebServer, "AddFunds")
}

func quote(r []string) {
	toWebServer := Quote{}
	toWebServer.UserId = r[1]
	toWebServer.StockSymbol = r[2]

	sendToWebServer(toWebServer, "GetQuote")
}

func buy(r []string) {
	toWebServer := Default{}
	toWebServer.UserId = r[1]
	toWebServer.StockSymbol = r[2]
	s, err := strconv.ParseFloat(r[3], 64)
	if err != nil {
		//handle error
		return
	}
	toWebServer.Amount = int(s * 10)

	sendToWebServer(toWebServer, "BuyStock")
}

func commitBuy(r []string) {
	toWebServer := User{}
	toWebServer.UserId = r[1]

	sendToWebServer(toWebServer, "CommitBuy")
}

func cancelBuy(r []string) {
	toWebServer := User{}
	toWebServer.UserId = r[1]

	sendToWebServer(toWebServer, "CancelBuy")
}

func sell(r []string) {

	toWebServer := Default{}
	toWebServer.UserId = r[1]
	toWebServer.StockSymbol = r[2]
	s, err := strconv.ParseFloat(r[3], 64)
	if err != nil {
		//handle error
		fmt.Println("sell error")

		return
	}
	toWebServer.Amount = int(s * 10)

	sendToWebServer(toWebServer, "SellStock")
}

func commitSell(r []string) {
	toWebServer := User{}
	toWebServer.UserId = r[1]

	sendToWebServer(toWebServer, "CommitSell")
}

func cancelSell(r []string) {
	toWebServer := User{}
	toWebServer.UserId = r[1]

	sendToWebServer(toWebServer, "CancelSell")
}

func setBuyAmount(r []string) {
	toWebServer := Default{}
	toWebServer.UserId = r[1]
	toWebServer.StockSymbol = r[2]
	s, err := strconv.ParseFloat(r[3], 64)
	if err != nil {
		//handle error
		return
	}
	toWebServer.Amount = int(s * 10)

	sendToWebServer(toWebServer, "SetBuyAmount")
}

func setBuyTrigger(r []string) {
	toWebServer := Default{}
	toWebServer.UserId = r[1]
	toWebServer.StockSymbol = r[2]
	s, err := strconv.ParseFloat(r[3], 64)
	if err != nil {
		//handle error
		return
	}
	toWebServer.Amount = int(s * 10)

	sendToWebServer(toWebServer, "SetBuyTrigger")
}

func cancelSetBuy(r []string) {
	toWebServer := Quote{}
	toWebServer.UserId = r[1]
	toWebServer.StockSymbol = r[2]

	sendToWebServer(toWebServer, "CancelSetBuy")
}

func setSellAmount(r []string) {
	toWebServer := Default{}
	toWebServer.UserId = r[1]
	toWebServer.StockSymbol = r[2]
	s, err := strconv.ParseFloat(r[3], 64)
	if err != nil {
		//handle error
		return
	}
	toWebServer.Amount = int(s * 10)

	sendToWebServer(toWebServer, "SetSellAmount")
}

func setSellTrigger(r []string) {
	toWebServer := Default{}
	toWebServer.UserId = r[1]
	toWebServer.StockSymbol = r[2]
	s, err := strconv.ParseFloat(r[3], 64)
	if err != nil {
		//handle error
		return
	}
	toWebServer.Amount = int(s * 10)

	sendToWebServer(toWebServer, "SetSellTrigger")
}

func cancelSetSell(r []string) {
	toWebServer := Quote{}
	toWebServer.UserId = r[1]
	toWebServer.StockSymbol = r[2]

	sendToWebServer(toWebServer, "CancelSetSell")
}

func dumplog(r []string) {
	if len(r) == 2 {
		//dumplog without username
		//dumplogall
	} else {
		//dumplog user
	}
}

func displaySummary(r []string) {
	toWebServer := User{}
	toWebServer.UserId = r[1]

	sendToWebServer(toWebServer, "DisplaySummary")
}

func sendToWebServer(r interface{}, s string) {
	jsonValue, _ := json.Marshal(r)
	resp, err := http.Post("http://localhost:8080/"+s, "application/json", bytes.NewBuffer(jsonValue))
	failOnError(err, "Error sending request")
	defer resp.Body.Close()
}

func main() {
	fmt.Println("Parsing workload file...")
	file, err := os.Open("workload.txt")
	failOnError(err, "Could not open file!")

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		commandText := scanner.Text()
		splitCommandText := strings.Fields(commandText)
		command := splitCommandText[1]
		//commandBytes := []byte(command)

		result := strings.Split(command, ",")
		//fmt.Println(result[0])
		switch result[0] {
		case "ADD":
			add(result)
		case "QUOTE":
			quote(result)
		case "BUY":
			buy(result)
		case "COMMIT_BUY":
			commitBuy(result)
		case "CANCEL_BUY":
			cancelBuy(result)
		case "SELL":
			sell(result)
		case "COMMIT_SELL":
			commitSell(result)
		case "CANCEL_SELL":
			cancelSell(result)
		case "SET_BUY_AMOUNT":
			setBuyAmount(result)
		case "CANCEL_SET_BUY":
			cancelSetBuy(result)
		case "SET_BUY_TRIGGER":
			setBuyTrigger(result)
		case "SET_SELL_AMOUNT":
			setSellAmount(result)
		case "CANCEL_SET_SELL":
			cancelSetSell(result)
		case "SET_SELL_TRIGGER":
			setSellTrigger(result)
		case "DUMPLOG":
			dumplog(result)
		case "DISPLAY_SUMMARY":
			displaySummary(result)
		}

	}

	failOnError(scanner.Err(), "Error reading file")

	fmt.Println("Done parsing workload file.")
}

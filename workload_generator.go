package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	client = &http.Client{}
	config = func() WorkloadGeneratorConfig {
		commandLineArgs := os.Args[1:]

		if len(commandLineArgs) < 1 || len(commandLineArgs) > 2 {
			fmt.Println("***ERROR INCORRECT ARGS*** Usage: numUsers webServerURL(optional)")
			os.Exit(0)
		}

		numUsers, _ := strconv.Atoi(commandLineArgs[0])

		if len(commandLineArgs) == 1 {
			return WorkloadGeneratorConfig{numUsers, "http://localhost:8080/"}

		}

		return WorkloadGeneratorConfig{numUsers, commandLineArgs[1]}
	}()
)

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
		panic(err)
	}
}

type WorkloadGeneratorConfig struct {
	NumUsers     int
	WebServerURL string
}

type Add struct {
	UserId         string
	Amount         int
	TransactionNum int
}
type Quote struct {
	UserId         string
	StockSymbol    string
	TransactionNum int
}

type Default struct {
	UserId         string
	StockSymbol    string
	Amount         int
	TransactionNum int
}

type User struct {
	UserId         string
	TransactionNum int
}

type Dumplog struct {
	Filename       string
	TransactionNum int
	Username       string
}

func add(r []string) {
	toWebServer := Add{}
	toWebServer.TransactionNum, _ = strconv.Atoi(r[0])
	toWebServer.UserId = r[2]
	toWebServer.Amount = floatStringToCents(r[3])

	sendToWebServer(toWebServer, "AddFunds")
}

func quote(r []string) {
	toWebServer := Quote{}
	toWebServer.TransactionNum, _ = strconv.Atoi(r[0])
	toWebServer.UserId = r[2]
	toWebServer.StockSymbol = r[3]

	sendToWebServer(toWebServer, "GetQuote")
}

func buy(r []string) {
	toWebServer := Default{}
	toWebServer.TransactionNum, _ = strconv.Atoi(r[0])
	toWebServer.UserId = r[2]
	toWebServer.StockSymbol = r[3]
	toWebServer.Amount = floatStringToCents(r[4])

	sendToWebServer(toWebServer, "BuyStock")
}

func commitBuy(r []string) {
	toWebServer := User{}
	toWebServer.TransactionNum, _ = strconv.Atoi(r[0])
	toWebServer.UserId = r[2]

	sendToWebServer(toWebServer, "CommitBuy")
}

func cancelBuy(r []string) {
	toWebServer := User{}
	toWebServer.TransactionNum, _ = strconv.Atoi(r[0])
	toWebServer.UserId = r[2]

	sendToWebServer(toWebServer, "CancelBuy")
}

func sell(r []string) {

	toWebServer := Default{}
	toWebServer.TransactionNum, _ = strconv.Atoi(r[0])
	toWebServer.UserId = r[2]
	toWebServer.StockSymbol = r[3]
	toWebServer.Amount = floatStringToCents(r[4])

	sendToWebServer(toWebServer, "SellStock")
}

func commitSell(r []string) {
	toWebServer := User{}
	toWebServer.TransactionNum, _ = strconv.Atoi(r[0])
	toWebServer.UserId = r[2]

	sendToWebServer(toWebServer, "CommitSell")
}

func cancelSell(r []string) {
	toWebServer := User{}
	toWebServer.TransactionNum, _ = strconv.Atoi(r[0])
	toWebServer.UserId = r[2]

	sendToWebServer(toWebServer, "CancelSell")
}

func setBuyAmount(r []string) {
	toWebServer := Default{}
	toWebServer.TransactionNum, _ = strconv.Atoi(r[0])
	toWebServer.UserId = r[2]
	toWebServer.StockSymbol = r[3]
	toWebServer.Amount = floatStringToCents(r[4])

	sendToWebServer(toWebServer, "SetBuyAmount")
}

func setBuyTrigger(r []string) {
	toWebServer := Default{}
	toWebServer.TransactionNum, _ = strconv.Atoi(r[0])
	toWebServer.UserId = r[2]
	toWebServer.StockSymbol = r[3]
	toWebServer.Amount = floatStringToCents(r[4])

	sendToWebServer(toWebServer, "SetBuyTrigger")
}

func cancelSetBuy(r []string) {
	toWebServer := Quote{}
	toWebServer.TransactionNum, _ = strconv.Atoi(r[0])
	toWebServer.UserId = r[2]
	toWebServer.StockSymbol = r[3]

	sendToWebServer(toWebServer, "CancelSetBuy")
}

func setSellAmount(r []string) {
	toWebServer := Default{}
	toWebServer.TransactionNum, _ = strconv.Atoi(r[0])
	toWebServer.UserId = r[2]
	toWebServer.StockSymbol = r[3]
	toWebServer.Amount = floatStringToCents(r[4])

	sendToWebServer(toWebServer, "SetSellAmount")
}

func setSellTrigger(r []string) {
	toWebServer := Default{}
	toWebServer.TransactionNum, _ = strconv.Atoi(r[0])
	toWebServer.UserId = r[2]
	toWebServer.StockSymbol = r[3]
	toWebServer.Amount = floatStringToCents(r[4])

	sendToWebServer(toWebServer, "SetSellTrigger")
}

func cancelSetSell(r []string) {
	toWebServer := Quote{}
	toWebServer.TransactionNum, _ = strconv.Atoi(r[0])
	toWebServer.UserId = r[2]
	toWebServer.StockSymbol = r[3]

	sendToWebServer(toWebServer, "CancelSetSell")
}

func dumplog(r []string) {
	if len(r) == 3 {
		toWebServer := Dumplog{}
		toWebServer.TransactionNum, _ = strconv.Atoi(r[0])
		toWebServer.Filename = r[2]
		sendToWebServer(toWebServer, "Dumplog")
		//dumplog without username
		//dumplogall
	} else if len(r) == 4 {
		toWebServer := Dumplog{}
		toWebServer.TransactionNum, _ = strconv.Atoi(r[0])
		toWebServer.Filename = r[2]
		toWebServer.Username = r[3]
		sendToWebServer(toWebServer, "Dumplog")
		//dumplog user
	} else {
		//what
		panic("unexpected number of params for dump log\n")
	}
}

func displaySummary(r []string) {
	toWebServer := User{}
	toWebServer.TransactionNum, _ = strconv.Atoi(r[0])
	toWebServer.UserId = r[2]

	sendToWebServer(toWebServer, "DisplaySummary")
}

func sendToWebServer(r interface{}, s string) {
	jsonValue, _ := json.Marshal(r)
	req, err := http.NewRequest("POST", config.WebServerURL+s, bytes.NewBuffer(jsonValue))
	req.Close = true

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if resp != nil {
		resp.Body.Close()
	}
	if err != nil {
		failOnError(err, "Error sending request")
	}
}

func floatStringToCents(val string) int {
	cents, _ := strconv.Atoi(strings.Replace(val, ".", "", 1))
	return cents
}

func doDumplog() {
	file, err := os.Open("split/testLOG")
	failOnError(err, "Could not open file!")

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		commandText := scanner.Text()
		commandText = strings.Replace(commandText, "[", "", 1)
		commandText = strings.Replace(commandText, "]", ",", 1)
		//commandBytes := []byte(command)
		result := strings.Split(commandText, ",")
		for index := range result {
			result[index] = strings.Replace(result[index], " ", "", 1)
		}
		dumplog(result)
	}
}

func handleUserFile(name string, wg *sync.WaitGroup) {

	file, err := os.Open("split/" + name)
	failOnError(err, "Could not open file!")

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		commandText := scanner.Text()
		commandText = strings.Replace(commandText, "[", "", 1)
		commandText = strings.Replace(commandText, "]", ",", 1)
		//commandBytes := []byte(command)
		result := strings.Split(commandText, ",")
		for index := range result {
			result[index] = strings.Replace(result[index], " ", "", 1)
		}
		commandText = strings.Replace(result[1], " ", "", 1)
		switch commandText {
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
	wg.Done()

}

func main() {
	fmt.Println("Parsing workload file...")

	files, err := ioutil.ReadDir("./split")
	if err != nil {
		log.Fatal(err)
	}

	startTime := time.Now()

	var wg sync.WaitGroup
	wg.Add(config.NumUsers)

	for _, f := range files {
		if f.Name() != "testLOG" {
			go handleUserFile(f.Name(), &wg)
		}
		fmt.Println(f.Name())
	}

	wg.Wait()

	doDumplog()

	elapsedTime := time.Since(startTime)

	fmt.Println("Done parsing workload file.")
	fmt.Printf("Took %s\n", elapsedTime)
}

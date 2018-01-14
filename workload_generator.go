package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
		panic(err)
	}
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
		commandBytes := []byte(command)

		resp, err := http.Post("http://localhost:8080", "application/text", bytes.NewBuffer(commandBytes))
		failOnError(err, "Error sending request")
		defer resp.Body.Close()
	}

	failOnError(scanner.Err(), "Error reading file")

	fmt.Println("Done parsing workload file.")
}

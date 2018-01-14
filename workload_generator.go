package main

import (
	"bufio"
	"fmt"
	"os"
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
		fmt.Println(scanner.Text())
	}

	failOnError(scanner.Err(), "Error reading file")
}

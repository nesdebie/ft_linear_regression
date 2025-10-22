package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func estimatePrice(mileage, theta0, theta1 float64) float64 {
	return theta0 + theta1*mileage
}

func readModel() (float64, float64) {
	file, err := os.Open("model.txt")
	if err != nil {
		fmt.Println("File not fount. Default values : theta0 = 0, theta1 = 0.")
		return 0.0, 0.0
	}
	defer file.Close()

	var theta0, theta1 float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "theta0=") {
			theta0, _ = strconv.ParseFloat(strings.TrimPrefix(line, "theta0="), 64)
		}
		if strings.HasPrefix(line, "theta1=") {
			theta1, _ = strconv.ParseFloat(strings.TrimPrefix(line, "theta1="), 64)
		}
	}
	return theta0, theta1
}

func main() {
	theta0, theta1 := readModel()

	fmt.Print("Input the mileage: ")
	var input string
	fmt.Scan(&input)
	mileage, err := strconv.ParseFloat(input, 64)
	if err != nil {
		fmt.Printf("Invalid mileage: %v\n", err)
		os.Exit(1)
	}

	price := estimatePrice(mileage, theta0, theta1)
	fmt.Printf("Estimated price: %.2f\n", price)
}

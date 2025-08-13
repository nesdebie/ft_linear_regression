package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
	"strconv"
)

type Model struct {
    Theta0 float64 `json:"theta0"`
    Theta1 float64 `json:"theta1"`
}

func estimatePrice(mileage, theta0, theta1 float64) float64 {
    return theta0 + theta1*mileage
}

func main() {
    modelFile, err := os.Open("model.json")
    if err != nil {
        log.Fatal(err)
    }
    defer modelFile.Close()

    var model Model
    err = json.NewDecoder(modelFile).Decode(&model)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Print("Input the mileage : ")
	var input string
    fmt.Scan(&input)
	mileage, err := strconv.ParseFloat(input, 64)
	if err != nil {
		log.Fatalf("Invalid mileage: %v", err)
	}
    price := estimatePrice(mileage, model.Theta0, model.Theta1)
    fmt.Printf("Estimated price : %.2f\n", price)
}

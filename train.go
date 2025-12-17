package main

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"os"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)


func estimatePrice(mileage, theta0, theta1 float64) float64 {
	return theta0 + theta1*mileage
}


func errorAndExit(reason string) {
	fmt.Println("error:", reason)
	os.Exit(1)
}


func createModelFile(theta0, theta1 float64) {
	f, err := os.Create("model.txt")
	if err != nil {
		errorAndExit(err.Error())
	}
	defer f.Close()
	fmt.Fprintf(f, "theta0=%.10f\n", theta0)
	fmt.Fprintf(f, "theta1=%.10f\n", theta1)
    fmt.Println("---------------------------\n")
	fmt.Println("Model saved to ./model.txt")
}


func createGraphicPlot(pts plotter.XYs, theta0, theta1 float64, input string) {
	p := plot.New()
	p.Title.Text = "Car Price vs Mileage"
	p.X.Label.Text = "Mileage (km)"
	p.Y.Label.Text = "Price"

	scatter, err := plotter.NewScatter(pts)
	if err != nil {
		errorAndExit(err.Error())
	}
	p.Add(scatter)

	line := plotter.NewFunction(func(x float64) float64 {
		return theta0 + theta1*x
	})
	line.Color = color.RGBA{R: 255, A: 255}
	line.Width = vg.Points(2)
	p.Add(line)

    filePath := "img/price_vs_mileage_" + input + ".png"
	if err := p.Save(8*vg.Inch, 5*vg.Inch, filePath); err != nil {
		errorAndExit(err.Error())
	}
	fmt.Println("Plot created ! Check ", filePath, "\n")
}


func parseData(file *os.File) ([]float64, []float64, plotter.XYs) {
	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		errorAndExit(err.Error())
	}
	if len(rows) < 1 {
		errorAndExit("no data")
	}
	var mileages, prices []float64
	pts := make(plotter.XYs, 0, len(rows)-1)

	for i, row := range rows {
		if i == 0 {
			continue
		}
		mileage, err := strconv.ParseFloat(row[0], 64)
		if err != nil {
			errorAndExit(err.Error())
		}
		price, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			errorAndExit(err.Error())
		}
		mileages = append(mileages, mileage)
		prices = append(prices, price)
		pts = append(pts, plotter.XY{X: mileage, Y: price})
	}

	if len(mileages) == 0 {
		errorAndExit("no data")
	}
	return mileages, prices, pts
}


func linearRegression(mileages, prices, normMil []float64, minMil, maxMil float64) (string, float64, float64) {
	m := float64(len(mileages))
	theta0 := 0.0
	theta1 := 0.0
	denormalizedTheta0 := 0.0
	denormalizedTheta1 := 0.0
	learningRate := 0.1

	fmt.Print("Amount of iterations: ")
	var input string
	fmt.Scan(&input)
	iterations, err := strconv.Atoi(input)
	if err != nil {
		fmt.Printf("Invalid integer: %v\n", err)
		os.Exit(1)
	}
    if iterations < 0 {
		fmt.Println("Invalid integer: Must be greater or equal to 0.")
		os.Exit(1)        
    }

	for i := 0; i < iterations; i++ {
		sum0, sum1 := 0.0, 0.0
		for j := 0; j < int(m); j++ {
			pred := estimatePrice(normMil[j], theta0, theta1)
			sum0 += pred - prices[j]
			sum1 += (pred - prices[j]) * normMil[j]
		}
		theta0 -= learningRate * (1/m) * sum0
		theta1 -= learningRate * (1/m) * sum1
    
        denormalizedTheta1 = theta1 / (maxMil - minMil)
		denormalizedTheta0 = theta0 - denormalizedTheta1*minMil
        fmt.Printf("Iteration #%d: Theta0 = %.5f, Theta1 = %.5f\n", i + 1, denormalizedTheta0, denormalizedTheta1)
	}
	return input, denormalizedTheta0, denormalizedTheta1
}


func main() {
	file, err := os.Open("data.csv")
	if err != nil {
		errorAndExit(err.Error())
	}
	defer file.Close()
	mileages, prices, pts := parseData(file)

	minMil, maxMil := mileages[0], mileages[0]
	for _, v := range mileages {
		if v < minMil {
			minMil = v
        }
		if v > maxMil {
			maxMil = v
        }
	}

	normMil := make([]float64, len(mileages))
	for i := range mileages {
		normMil[i] = (mileages[i] - minMil) / (maxMil - minMil)
    }

	input, theta0, theta1 := linearRegression(mileages, prices, normMil, minMil, maxMil)
	createModelFile(theta0, theta1)
	createGraphicPlot(pts, theta0, theta1, input)

    fmt.Println("-----theta0 and theta1-----")
	fmt.Printf("Theta0 = %.5f\nTheta1 = %.5f\n", theta0, theta1)
}

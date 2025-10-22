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

func main() {
	file, err := os.Open("data.csv")
	if err != nil {
		errorAndExit(err.Error())
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		errorAndExit(err.Error())
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

	m := float64(len(mileages))
	theta0 := 0.0
	theta1 := 0.0
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
    
        tmpTheta1 := theta1 / (maxMil - minMil)
        fmt.Printf("Iteration #%d: Theta0 = %.5f, Theta1 = %.5f\n", i + 1, theta0 - tmpTheta1*minMil, theta1 / (maxMil - minMil))
	}

	realTheta1 := theta1 / (maxMil - minMil)
	realTheta0 := theta0 - realTheta1*minMil

	f, err := os.Create("model.txt")
	if err != nil {
		errorAndExit(err.Error())
	}
	defer f.Close()
	fmt.Fprintf(f, "theta0=%.10f\n", realTheta0)
	fmt.Fprintf(f, "theta1=%.10f\n", realTheta1)
	fmt.Fprintf(f, "min_mileage=%.2f\n", minMil)
	fmt.Fprintf(f, "max_mileage=%.2f\n", maxMil)
    fmt.Println("---------------------------\n")
	fmt.Println("Model saved to ./model.txt")

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
		return realTheta0 + realTheta1*x
	})
	line.Color = color.RGBA{R: 255, A: 255}
	line.Width = vg.Points(2)
	p.Add(line)

    filePath := "img/price_vs_mileage_" + input + ".png"
	if err := p.Save(8*vg.Inch, 5*vg.Inch, filePath); err != nil {
		errorAndExit(err.Error())
	}
	fmt.Println("Plot created ! Check ./img/price_vs_mileage.png\n")

    fmt.Println("-----theta0 and theta1-----")
	fmt.Printf("Theta0 = %.5f\nTheta1 = %.5f\n", realTheta0, realTheta1)
}

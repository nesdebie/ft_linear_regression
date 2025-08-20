package main

import (
    "encoding/csv"
    "encoding/json"
    "fmt"
    "os"
    "strconv"

    "gonum.org/v1/plot"
    "gonum.org/v1/plot/plotter"
    "gonum.org/v1/plot/vg"
)

type Model struct {
    Theta0     float64 `json:"theta0"`
    Theta1     float64 `json:"theta1"`
    MinMileage float64 `json:"min_mileage"`
    MaxMileage float64 `json:"max_mileage"`
}

func estimatePrice(mileage, theta0, theta1 float64) float64 {
    return theta0 + theta1*mileage
}

func errorAndExit(reason string) {
    fmt.Println(reason)
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
        errorAndExit("error")
    }

    p := plot.New()
    p.Title.Text = "Car Price vs Mileage"
    p.X.Label.Text = "Mileage (km)"
    p.Y.Label.Text = "Price"

    scatter, err := plotter.NewScatter(pts)
    if err != nil {
        errorAndExit(err.Error())
    }
    p.Add(scatter)

    if err := p.Save(6*vg.Inch, 4*vg.Inch, "img/price_vs_mileage.png"); err != nil {
        errorAndExit(err.Error())
    }
    fmt.Println("Plot saved to img/price_vs_mileage.png")

    minMil, maxMil := mileages[0], mileages[0]
    for _, v := range mileages {
        if v < minMil {
            minMil = v
        }
        if v > maxMil {
            maxMil = v
        }
    }
    for i := range mileages {
        mileages[i] = (mileages[i] - minMil) / (maxMil - minMil)
    }

    m := float64(len(mileages))
    theta0 := 0.0
    theta1 := 0.0
    learningRate := 0.01
    iterations := 1000

    for i := 0; i < iterations; i++ {
        sum0 := 0.0
        sum1 := 0.0
        for j := 0; j < int(m); j++ {
            prediction := estimatePrice(mileages[j], theta0, theta1)
            sum0 += (prediction - prices[j])
            sum1 += (prediction - prices[j]) * mileages[j]
        }
        tmp0 := theta0 - learningRate*(1/m)*sum0
        tmp1 := theta1 - learningRate*(1/m)*sum1
        theta0, theta1 = tmp0, tmp1

        if i%100 == 0 {
            fmt.Printf("Iteration #%d -> theta0=%.5f theta1=%.5f\n", i, theta0, theta1)
        }
    }

    model := Model{
        Theta0:     theta0,
        Theta1:     theta1,
        MinMileage: minMil,
        MaxMileage: maxMil,
    }

    modelFile, err := os.Create("model.json")
    if err != nil {
        errorAndExit(err.Error())
    }
    defer modelFile.Close()
    enc := json.NewEncoder(modelFile)
    enc.SetIndent("", "  ")
    if err := enc.Encode(model); err != nil {
        errorAndExit(err.Error())
    }

    fmt.Println("Model trained:")
    fmt.Printf("Theta0 = %.5f, Theta1 = %.5f\n", theta0, theta1)
    fmt.Printf("MinMileage = %.2f, MaxMileage = %.2f\n", minMil, maxMil)
}

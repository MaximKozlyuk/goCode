package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

/*
	Ва 12
	моменты времени ∆ti:  1 2 3 4 5 6 7
	Кол-во ошибок ∆Ni:    4 4 3 3 2 1 0
*/

/*
//enter Ni:
fmt.Println("∆Ni:")
for i := 0; i < len(Ni); i++ {
	fmt.Scanf("%d", &Ni[i])
}
fmt.Println(Ni)
*/

func init() {

}

func main() {

	eps := 0.00001

	// 12 variant
	//Ni := make([]int, 7,7)
	//Ni[0] = 4; Ni[1] = 4; Ni[2] = 3 ;Ni[3] = 3; Ni[4] = 2; Ni[5] = 1; Ni[6] = 0

	// test variant // todo init all values {}
	// todo make Ni arr float64, recode with links and draw init data
	Ni := make([]int, 6, 6)
	Ni[0] = 3
	Ni[1] = 4
	Ni[2] = 2
	Ni[3] = 1
	Ni[4] = 1
	Ni[5] = 0

	K, counter := halfLengthMethod(0, 100, eps, Ni)
	fmt.Println("K found after", counter, "iterations")
	fmt.Print("k = ", K, "\nNo = ")
	No := calc_No(K, 1, Ni)
	fmt.Println(No)

	// ∆ni расчетное :
	fmt.Println()
	var errorsLeft float64
	arrY := make([]float64, 0, 0)
	arrX := make([]float64, 0, 0)
	for i := 0; i < len(Ni); i++ {
		fmt.Print("∆n")
		fmt.Println(i, " = ", No*math.Pow(math.E, -K*float64(i+1))*K)
		errorsLeft += No * math.Pow(math.E, -K*float64(i+1)) * K

		arrY = append(arrY, No*math.Pow(math.E, -K*float64(i+1))*K)
		arrX = append(arrX, float64(i))
	}
	fmt.Println("Errors left: ", No-errorsLeft)

	renderPlot(&arrX, &arrY)

}

// ∑ ∆ni * e^(-K * ti)
func rightEquationPart(k float64, ni []int) float64 {
	var sum float64
	for i := 0; i < len(ni); i++ {
		sum += float64(ni[i]) * math.Pow(math.E, -k*float64(i+1))
	}
	return sum
}

/*
					∑ ∆ni * e^(-Kti) * ti
	(∑ e^(-2Kti)) * ---------------------
					   ∑ e^(-2Kti) * ti
*/
func leftEquationPart(k float64, ni []int) float64 {
	var leftPart, numerator, denominator float64
	for i := 0; i < len(ni); i++ {
		leftPart += math.Pow(math.E, -2*k*float64(i+1))
		numerator += float64(ni[i]) * math.Pow(math.E, -k*float64(i+1)) * float64(i+1)
		denominator += math.Pow(math.E, -2*k*float64(i+1)) * float64(i+1)
	}
	return leftPart * (numerator / denominator)
}

func difference(k float64, ni []int) float64 {
	return leftEquationPart(k, ni) - rightEquationPart(k, ni)
}

/*
	метод деления отрезка пополам, возвращает K и кол-во итераций поиска
	func difference здесь играет роль искомого F(k)
*/
func halfLengthMethod(a, b, eps float64, ni []int) (float64, int) {
	var (
		c       float64
		counter int
	)
	for (b - a) >= 2*eps {
		c = (a + b) / 2
		if difference(a, ni)*difference(c, ni) > 0 {
			a = c
		} else {
			b = c
		}
		counter++
	}
	return (a + b) / 2, counter
}

/*
            ∑ ∆ni * e^(-Kti) * ti
	N0 = ---------------------------
		 K * ∆t * ∑ (e^(-2Kti) * ti)
*/
func calc_No(k, Dt float64, ni []int) float64 {
	var numerator, denominator float64
	for i := 0; i < len(ni); i++ {
		numerator += float64(ni[i]) * float64(i+1) * math.Pow(math.E, -k*float64(i+1))
		denominator += math.Pow(math.E, -2*k*float64(i+1)) * float64(i+1)
	}
	return numerator / (k * Dt * denominator)
}

func renderPlot(arrX, arrY *[]float64) {
	// Get values
	values := packArrayForPlot(arrX, arrY)

	// Create a new plot, set its title and
	// axis labels.
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Модель джевинского-Моранды"
	p.X.Label.Text = "∆t"
	p.Y.Label.Text = "n"
	// Draw a grid behind the data
	p.Add(plotter.NewGrid())

	// Make a scatter plotter and set its style.
	s, err := plotter.NewLine(values)
	if err != nil {
		panic(err)
	}
	s.LineStyle.Color = color.RGBA{R: 255, B: 128, A: 255}

	// Add the plotters to the plot, with a legend
	// entry for each
	p.Add(s)
	p.Legend.Add("Решение моджели", s)

	// Save the plot to a PNG file.
	if err := p.Save(12*vg.Inch, 12*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}

func packArrayForPlot(arrX, arrY *[]float64) plotter.XYs {
	pts := make(plotter.XYs, len(*arrX))
	for i := range *arrX {
		pts[i].X = (*arrX)[i]
		pts[i].Y = (*arrY)[i]
	}
	return pts
}

// randomPoints returns some random x, y points for plot
func randomPoints(n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {
		rand.Seed(time.Now().UnixNano())
		pts[i].X = rand.Float64() * float64(rand.Intn(100))
		pts[i].Y = rand.Float64() * float64(rand.Intn(100))
	}
	return pts
}

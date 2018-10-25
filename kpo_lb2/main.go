package main

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
	"math"
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
	Ni := make([]int, 7, 7)
	Ni[0] = 4
	Ni[1] = 4
	Ni[2] = 3
	Ni[3] = 3
	Ni[4] = 2
	Ni[5] = 1
	Ni[6] = 0

	// todo make Ni arr float64, recode with links
	//b := [2]string{"Penn", "Teller"}
	//Ni := make([]int, 6,6)
	//Ni[0] = 3
	//Ni[1] = 4
	//Ni[2] = 2
	//Ni[3] = 1
	//Ni[4] = 1
	//Ni[5] = 0

	K, counter := halfLengthMethod(0, 100, eps, Ni)
	fmt.Println("K found after", counter, "iterations")
	fmt.Print("k = ", K, "\nNo = ")
	No := calc_No(K, 1, Ni)
	fmt.Println(No)

	// ∆ni расчетное :
	fmt.Println()
	var errorsLeft float64
	// errors by the init values
	arrY := make([]float64, 0, 0)
	arrX := make([]float64, 0, 0)
	for i := 0; i < len(Ni); i++ {
		fmt.Print("∆n")
		fmt.Println(i, " = ", errorsLeftInDt(No, K, float64(i+1)))
		errorsLeft += errorsLeftInDt(No, K, float64(i+1))

		arrY = append(arrY, errorsLeftInDt(No, K, float64(i+1)))
		arrX = append(arrX, float64(i))
	}
	fmt.Println("Errors left: ", No-errorsLeft)

	// arrays with init data
	initX := make([]float64, 0, 0)
	initY := make([]float64, 0, 0)
	for i := 0; i < len(Ni); i++ {
		initX = append(initX, float64(i))
		initY = append(initY, float64(Ni[i]))
	}

	// generate arrays for extrapolation
	extraX := make([]float64, 0, 0)
	extraY := make([]float64, 0, 0)
	for i := 0.0; i < 10; i += 0.01 {
		extraX = append(extraX, i)
		extraY = append(extraY, No*math.Pow(math.E, -K*i)*K)
	}

	// получение точек невязок
	nevs := getNevsY(No, K, 1, Ni)
	news2 := make(plotter.XYs, len(Ni))
	var y, nevs2Sum float64
	for i := 0; i < nevs.Len(); i++ {
		news2[i].X = float64(i + 1)
		_, y = nevs.XY(i)
		news2[i].Y = math.Pow(y, 2)
		nevs2Sum += math.Pow(y, 2)
	}
	fmt.Println("S =", nevs2Sum)

	renderPlot(
		*packArrayForPlot(&extraX, &extraY),
		*packArrayForPlot(&initX, &initY),
		*nevs,
		news2)

}

func getNevsY(No, K, Dt float64, ni []int) *plotter.XYs {
	pts := make(plotter.XYs, len(ni))
	for i := 0; i < len(ni); i++ {
		pts[i].X = float64(i + 1)
		pts[i].Y = float64(ni[i]) - No*K*Dt*math.Pow(math.E, -K*float64(i+1))
	}
	return &pts
}

func errorsLeftInDt(No, K, Dt float64) float64 {
	return No * math.Pow(math.E, -K*Dt) * K
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

func renderPlot(modelSolve, initData, nevs, nevs2 plotter.XYs) {
	// Create a new plot, set its title and
	// axis labels.
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Модель Джелинского-Моранды"
	p.X.Label.Text = "∆t"
	p.Y.Label.Text = "n"
	// Draw a grid behind the data
	p.Add(plotter.NewGrid())

	// Make a scatter plotter and set its style.
	line, err := plotter.NewLine(modelSolve)
	if err != nil {
		panic(err)
	}
	line.LineStyle.Color = color.RGBA{B: 255, A: 255}

	dots, err := plotter.NewScatter(initData)
	if err != nil {
		panic(err)
	}
	//dots.GlyphStyle.Color = color.RGBA{R: 255, G: 0, B: 0, A: 1}
	dots.Color = color.RGBA{R: 255, A: 255}

	nev, err := plotter.NewScatter(nevs)
	if err != nil {
		panic(err)
	}
	nev.Shape = draw.PyramidGlyph{}
	nev.Color = color.RGBA{G: 255, A: 255}

	nev2, err := plotter.NewScatter(nevs2)
	if err != nil {
		panic(err)
	}
	nev2.Shape = draw.BoxGlyph{}
	nev2.Color = color.RGBA{B: 255, A: 255}

	// Add the plotters to the p, with a legend
	// entry for each
	p.Add(line, dots, nev, nev2)
	p.Legend.Add("Решение моджели", line)
	p.Legend.Add("Исходные данные об ошибках", dots)
	p.Legend.Add("Значения невязок", nev)
	p.Legend.Add("Квадраты невязок", nev2)

	// Save the plot to a PNG file.
	if err := p.Save(10*vg.Inch, 10*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}

func packArrayForPlot(arrX, arrY *[]float64) *plotter.XYs {
	pts := make(plotter.XYs, len(*arrX))
	for i := range *arrX {
		pts[i].X = (*arrX)[i]
		pts[i].Y = (*arrY)[i]
	}
	return &pts
}

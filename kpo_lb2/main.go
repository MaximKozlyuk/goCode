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

func main() {

	eps := 0.00001
	Ni := make([]int, 7, 7)
	Ni[0] = 4
	Ni[1] = 4
	Ni[2] = 3
	Ni[3] = 3
	Ni[4] = 2
	Ni[5] = 1
	Ni[6] = 0

	// test
	//Ni := make([]int, 6,6)
	//Ni[0] = 3
	//Ni[1] = 4
	//Ni[2] = 2
	//Ni[3] = 1
	//Ni[4] = 1
	//Ni[5] = 0

	K, counter := halfLengthMethod(0, 100, eps, Ni)
	fmt.Println("K найдена после", counter, "итераций чм")
	fmt.Print("k = ", K, "\nNo = ")
	No := calc_No(K, 1, Ni)
	fmt.Println(No)

	fmt.Println()
	var errorsLeft float64
	// расчет точек кол-ва оствшихся ошибок
	for i := 0; i < len(Ni); i++ {
		fmt.Print("∆n")
		fmt.Println(i, " = ", errorsLeftInDt(No, K, float64(i+1)))
		errorsLeft += errorsLeftInDt(No, K, float64(i+1))
	}
	fmt.Println("Осталось ошибок: ", No-errorsLeft)

	// исходные данные для графика
	initX := make([]float64, 0, 0)
	initY := make([]float64, 0, 0)
	for i := 0; i < len(Ni); i++ {
		initX = append(initX, float64(i+1))
		initY = append(initY, float64(Ni[i]))
	}

	// массивы с решением модели и экстраполяцией
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
	fmt.Println("S =", nevs2Sum) // вывод суммы квадратов невязок

	// нахождение коэффициэнтов линейной аппроксимации, уравнение прямой: y = ax + b
	var a, b float64
	lineApprox(&initX, &Ni, &a, &b)
	fmt.Println("a =", a, " b =", b)
	lineApproaxPoints := make(plotter.XYs, 2, 2)
	lineApproaxPoints[0].X = 0
	lineApproaxPoints[1].X = 10
	lineApproaxPoints[0].Y = a + b*lineApproaxPoints[0].X
	lineApproaxPoints[1].Y = a + b*lineApproaxPoints[1].X

	// отрисовка графика
	renderPlot(
		*packArrayForPlot(&extraX, &extraY),
		*packArrayForPlot(&initX, &initY),
		*nevs,
		news2,
		lineApproaxPoints)
}

//линейная апроксимация xy[0] содержит иксы и [1] соответственно игрики
func lineApprox(x *[]float64, y *[]int, a, b *float64) {
	var sx, sx2, sxy, sy float64
	for i := 0; i < len(*x); i++ {
		sx += (*x)[i]
		sx2 += math.Pow((*x)[i], 2)
		sxy += (*x)[i] * float64((*y)[i])
		sy += float64((*y)[i])
	}
	*a = (sx2*sy - sx*sxy) / (float64(len(*x))*sx2 - sx*sx)
	*b = (float64(len(*x))*sxy - sx*sy) / (float64(len(*x))*sx2 - sx*sx)
}

// возвращает массив точек невязок
func getNevsY(No, K, Dt float64, ni []int) *plotter.XYs {
	pts := make(plotter.XYs, len(ni))
	for i := 0; i < len(ni); i++ {
		pts[i].X = float64(i + 1)
		pts[i].Y = float64(ni[i]) - No*K*Dt*math.Pow(math.E, -K*float64(i+1))
	}
	return &pts
}

// расчет кол-ва оставшихся ошибок No * e^(-K∆t) * K
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

//метод деления отрезка пополам, возвращает K и кол-во итераций поиска
//func difference здесь играет роль искомого F(k)
func halfLengthMethod(a, b, eps float64, ni []int) (float64, int) {
	fmt.Println("Метод деления отрезка пополам:")
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
		fmt.Println(counter, ":", (a+b)/2)
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

func renderPlot(modelSolve, initData, nevs, nevs2, lineApproaxPoints plotter.XYs) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Модель Джелинского-Моранды"
	p.X.Label.Text = "t"
	p.Y.Label.Text = "∆n"
	p.Add(plotter.NewGrid())

	// решение модели с экстраполяцией
	line, err := plotter.NewLine(modelSolve)
	if err != nil {
		panic(err)
	}
	line.LineStyle.Color = color.RGBA{B: 255, A: 255}

	// апроксимация прямой
	lineApproax, err := plotter.NewLine(lineApproaxPoints)
	if err != nil {
		panic(err)
	}
	lineApproax.LineStyle.Color = color.RGBA{R: 255, A: 255}

	// исходные данные
	dots, err := plotter.NewScatter(initData)
	if err != nil {
		panic(err)
	}
	dots.Color = color.RGBA{R: 255, A: 255}

	// значения невязок
	nev, err := plotter.NewScatter(nevs)
	if err != nil {
		panic(err)
	}
	nev.Shape = draw.PyramidGlyph{}
	nev.Color = color.RGBA{G: 255, A: 255}

	// квадраты невязок
	nev2, err := plotter.NewScatter(nevs2)
	if err != nil {
		panic(err)
	}
	nev2.Shape = draw.BoxGlyph{}
	nev2.Color = color.RGBA{B: 255, A: 255}

	// добавление графиков и их легенд
	p.Add(line, dots, nev, nev2, lineApproax)
	p.Legend.Add("Решение моджели", line)
	p.Legend.Add("Апроксимация решения прямой", lineApproax)
	p.Legend.Add("Исходные данные об ошибках", dots)
	p.Legend.Add("Значения невязок", nev)
	p.Legend.Add("Квадраты невязок", nev2)

	// сохранение графика в PNG файл
	if err := p.Save(10*vg.Inch, 10*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}

// создание структуры plotter.XYs с точками для отрисовки графика
func packArrayForPlot(arrX, arrY *[]float64) *plotter.XYs {
	pts := make(plotter.XYs, len(*arrX))
	for i := range *arrX {
		pts[i].X = (*arrX)[i]
		pts[i].Y = (*arrY)[i]
	}
	return &pts
}

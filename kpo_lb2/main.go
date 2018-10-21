package main

import (
	"fmt"
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
	//Ni := make([]int, 7,7)
	//Ni[0] = 4; Ni[1] = 4; Ni[2] = 3 ;Ni[3] = 3; Ni[4] = 2; Ni[5] = 1; Ni[6] = 0

	// test variant
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
	for i := 0; i < len(Ni); i++ {
		fmt.Print("∆n")
		fmt.Println(i, " = ", No*math.Pow(math.E, -K*float64(i+1))*K)
		errorsLeft += No * math.Pow(math.E, -K*float64(i+1)) * K
	}
	fmt.Println("Errors left: ", No-errorsLeft)

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

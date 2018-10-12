package main

import "fmt"

func main() {

	mapa := make(map[int]int)

	fmt.Println(len(mapa))

	mapa[1] = 1
	mapa[1] = 1

	fmt.Println(len(mapa))

	mapa[2] = 2
	mapa[3] = 3

	fmt.Println(len(mapa))

}

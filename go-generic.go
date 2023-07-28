package main

import "log"

func sumAnyType[T int | float64](i []T) (o T) {
	for _, v := range i {
		o += v
	}
	return
}

func main() {

	log.Println(sumAnyType([]int{1, 2, 3})) // prints 6

	log.Println(sumAnyType([]float64{1.2, 2.5, 3.9})) // prints 7.6

	// log.Println(sumAnyType([]float32{1.2, 2.5, 3.9})) // does not compile
}

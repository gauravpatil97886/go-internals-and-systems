// Sum of Slice

// Problem
// nums := []int{1,2,3,4,5}

// Return sum:

// 15



package main

import "fmt"

func main() {
	sums := 0
	nums := []int{1, 2, 3, 4, 5}

	for _, v := range nums {
		sums = sums + v
	}
	fmt.Println(sums)

}

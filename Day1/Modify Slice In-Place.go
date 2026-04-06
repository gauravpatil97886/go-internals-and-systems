// Problem 2 — Modify Slice In-Place
// Problem
// nums := []int{1, 2, 3, 4}

// Convert it to:

// [2, 4, 6, 8]

package main

import "fmt"

func main() {
	nums := []int{1, 2, 3, 4}

	for i, _ := range nums {
		nums[i] = nums[i] * 2
	}
	fmt.Print(nums)
}

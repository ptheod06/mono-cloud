package main

import (
	"sort"
	"fmt"
)


func main() {

	var arr = []int{10, 5, 6, 2, 8 ,53, 12 ,46, 23}

	var largest = []int{}

	for i := 0; i < 3; i++ {
		largest = append(largest, arr[i])
	}

	sort.Ints(largest)

	for i := 3; i < len(arr); i++ {

		if (largest[0] < arr[i]) {
			largest[0] = arr[i]
			sort.Ints(largest)
		}

	}

	fmt.Println(largest)
}

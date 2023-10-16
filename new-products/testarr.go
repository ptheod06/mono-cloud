package main

import (
	"sort"
	"fmt"
)


func main() {

	var arr = []int{10, 5, 6, 2, 8 ,53, 12 ,46, 23}

	var largest = []int{}

	for i := 0; i< 3; i++ {
		largest = append(largest, arr[i])
	}

	sort.Ints(largest)

	fmt.Println(largest)


//	fmt.Println(append(arr[:3], 1000))


//	fmt.Println(arr[0:0])
	for i := 3; i < len(arr); i++ {
		index := -1
		var remaining = []int{}
//		fmt.Println("largest: ",largest)
//		fmt.Println(i)
		if (arr[i] > largest[0]) {
			for j := 1; j < len(largest); j++ {
				fmt.Println("largest ", largest)
				fmt.Println("comparing ", arr[i], largest[j])
				if (arr[i] < largest[j]) {
					index = j-1
					fmt.Println("insert at ", index)
					break
				}
				if (j == 2) {

					largest = append(largest[1:],arr[i])
				}
			}

			if (index == -1) {
				continue
			}
			remaining = append(remaining[:], largest[index+1:]...)
			fmt.Println("remainng", remaining)
			largest = append(largest[1:index+1], arr[i])
			fmt.Println("before final", largest)
			fmt.Println("app", remaining)
			largest = append(largest[:], remaining...)
			
		}

	}

	fmt.Println(largest)

	fmt.Println(arr)

}

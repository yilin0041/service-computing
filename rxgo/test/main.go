package main

import (
	"fmt"
	"time"

	rxgo "github.com/yilin0041/service-computing/rxgo"
)

func main() {
	fmt.Println("测试数据：0,1,2,3,4,5,3,4,5")
	res := []int{}
	ob := rxgo.Just(0, 1, 2, 3, 4, 5, 3, 4, 5).Map(func(x int) int {
		return x
	}).Debounce(999999)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("Debounce(999999): ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n")
	res = []int{}
	ob = rxgo.Just(0, 1, 2, 3, 4, 5, 3, 4, 5).Map(func(x int) int {
		return x
	}).Distinct()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("Distinct(): ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n")
	res = []int{}
	ob = rxgo.Just(0, 1, 2, 3, 4, 5, 3, 4, 5).Map(func(x int) int {
		return x
	}).ElementAt(5)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("ElementAt(5): ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n")

	res = []int{}
	ob = rxgo.Just(0, 1, 2, 3, 4, 5, 3, 4, 5).Map(func(x int) int {
		return x
	}).First()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("First(): ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n")

	res = []int{}
	ob = rxgo.Just(0, 1, 2, 3, 4, 5, 3, 4, 5).Map(func(x int) int {
		return x
	}).Last()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("Last(): ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n")

	res = []int{}
	rxgo.Just(0, 1, 2, 3, 4, 5, 3, 4, 5).Map(func(x int) int {
		time.Sleep(4 * time.Millisecond)
		return x
	}).Sample(40 * time.Millisecond).Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("Sample(40): (sleep (4))")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n")

	res = []int{}
	ob = rxgo.Just(0, 1, 2, 3, 4, 5, 3, 4, 5).Map(func(x int) int {
		return x
	}).Skip(4)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("Skip(4): ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n")
	res = []int{}
	ob = rxgo.Just(0, 1, 2, 3, 4, 5, 3, 4, 5).Map(func(x int) int {
		return x
	}).SkipLast(2)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("SkipLast(2): ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n")

	res = []int{}
	ob = rxgo.Just(0, 1, 2, 3, 4, 5, 3, 4, 5).Map(func(x int) int {
		return x
	}).Take(4)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("Take(4) :")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n")
	res = []int{}
	ob = rxgo.Just(0, 1, 2, 3, 4, 5, 3, 4, 5).Map(func(x int) int {
		return x
	}).TakeLast(2)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("TakeLast(2): ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n")
}

package rxgo_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yilin0041/service-computing/rxgo"
)

func TestDebounce(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(0, 1, 2, 3, 4, 5).Map(func(x int) int {
		return x
	}).Debounce(100 * time.Millisecond)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})

	assert.Equal(t, []int{}, res, "Debounce Test Error!")
}

func TestDistinct(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(0, 1, 2, 3, 4, 5, 3, 4, 5, 6).Map(func(x int) int {
		return x
	}).Distinct()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6}, res, "Distinct Test Error!")
}

func TestElementAt(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(0, 1, 2, 3, 4, 5).Map(func(x int) int {
		return x
	}).ElementAt(5)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{4}, res, "ElementAt Test Error!")
}

func TestIgnoreElement(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(0, 1, 2, 3, 4, 5).Map(func(x int) int {
		return x
	}).IgnoreElement()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{}, res, "IgnoreElement Test Error!")
}
func TestFirst(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(0, 1, 2, 3, 4, 5).Map(func(x int) int {
		return x
	}).First()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})

	assert.Equal(t, []int{0}, res, "First Test Error!")
}

func TestLast(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(0, 1, 2, 3, 4, 5).Map(func(x int) int {
		return x
	}).Last()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})

	assert.Equal(t, []int{5}, res, "Last Test Error!")
}

func TestSample(t *testing.T) {
	res := []int{}
	rxgo.Just(0, 1, 2, 3, 4, 5).Map(func(x int) int {
		time.Sleep(2 * time.Millisecond)
		return x
	}).Sample(20 * time.Millisecond).Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{}, res, "SkipLast Test Error!")

}

func TestSkip(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(0, 1, 2, 3, 4, 5).Map(func(x int) int {
		return x
	}).Skip(2)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{2, 3, 4, 5}, res, "Skip Test Error!")
}

func TestSkipLast(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(0, 1, 2, 3, 4, 5).Map(func(x int) int {
		return x
	}).SkipLast(3)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{0, 1, 2}, res, "SkipLast Test Error!")
}

func TestTake(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(0, 1, 2, 3, 4, 5).Map(func(x int) int {
		return x
	}).Take(2)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{0, 1}, res, "Take Test Error!")
}

func TestTakeLast(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(0, 1, 2, 3, 4, 5).Map(func(x int) int {
		return x
	}).TakeLast(3)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{3, 4, 5}, res, "TakeLast Test Error!")
}

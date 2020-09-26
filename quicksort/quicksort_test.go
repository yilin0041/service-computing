package quicksort

import "testing"

func TestQuickSort(t *testing.T) {
	var arr = []int{9, 5, 6, 7, 8, 1, 0, 2, 4, 3}
	QuickSort(arr, 0, 9)
	expected := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i := 0; i < 10; i++ {
		if arr[i] != expected[i] {
			t.Errorf("arr[%d] expected %d,but got %d\n", i, expected[i], arr[i])
		}
	}
}

func TestBubbleSort(t *testing.T) {
	var arr = []int{9, 5, 6, 7, 8, 1, 0, 2, 4, 3}
	BubbleSort(arr, 10)
	expected := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i := 0; i < 10; i++ {
		if arr[i] != expected[i] {
			t.Errorf("arr[%d] expected %d,but got %d\n", i, expected[i], arr[i])
		}
	}
}

func BenchmarkQuickSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var arr = []int{
			9, 5, 6, 7, 8, 1, 0, 2, 4, 3,
			9, 5, 6, 7, 8, 1, 0, 2, 4, 3,
			9, 5, 6, 7, 8, 1, 0, 2, 4, 3,
			9, 5, 6, 7, 8, 1, 0, 2, 4, 3,
			9, 5, 6, 7, 8, 1, 0, 2, 4, 3,
			9, 5, 6, 7, 8, 1, 0, 2, 4, 3,
			9, 5, 6, 7, 8, 1, 0, 2, 4, 3,
			9, 5, 6, 7, 8, 1, 0, 2, 4, 3,
			9, 5, 6, 7, 8, 1, 0, 2, 4, 3,
			9, 5, 6, 7, 8, 1, 0, 2, 4, 3}
		QuickSort(arr, 0, 99)
	}
}

func BenchmarkBubbleSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var arr = []int{
			9, 5, 6, 7, 8, 1, 0, 2, 4, 3,
			9, 5, 6, 7, 8, 1, 0, 2, 4, 3,
			9, 5, 6, 7, 8, 1, 0, 2, 4, 3,
			9, 5, 6, 7, 8, 1, 0, 2, 4, 3,
			9, 5, 6, 7, 8, 1, 0, 2, 4, 3,
			9, 5, 6, 7, 8, 1, 0, 2, 4, 3,
			9, 5, 6, 7, 8, 1, 0, 2, 4, 3,
			9, 5, 6, 7, 8, 1, 0, 2, 4, 3,
			9, 5, 6, 7, 8, 1, 0, 2, 4, 3,
			9, 5, 6, 7, 8, 1, 0, 2, 4, 3}
		BubbleSort(arr, 100)
	}
}

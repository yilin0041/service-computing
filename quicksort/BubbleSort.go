package quicksort

//BubbleSort :the sort function
func BubbleSort(arr []int, len int) {
	for i := 0; i < len-1; i++ {
		for j := i + 1; j < len; j++ {
			if arr[j] < arr[i] {
				temp := arr[j]
				arr[j] = arr[i]
				arr[i] = temp
			}
		}
	}
}

package quicksort

//QuickSort :the sort function
func QuickSort(arr []int, s int, e int) {
	pivot := arr[s]
	pos := s
	i := s
	j := e
	for i <= j {
		for j >= pos && arr[j] >= pivot {
			j--
		}
		if j >= pos {
			arr[pos] = arr[j]
			pos = j
		}
		for i <= pos && arr[i] <= pivot {
			i++
		}
		if i <= pos {
			arr[pos] = arr[i]
			pos = i
		}
	}
	arr[pos] = pivot
	if pos-s > 1 {
		QuickSort(arr, s, pos-1)
	}
	if e-pos > 1 {
		QuickSort(arr, pos+1, e)
	}
}

package util

// HeapSort 堆排序
func HeapSort(arr []int) {
	heapAdjust := func(v []int, i, length int) {
		temp := v[i]                                //先取出当前元素i
		for k := i*2 + 1; k < length; k = k*2 + 1 { //从i结点的左子结点开始，也就是2i+1处开始
			if k+1 < length && v[k] < v[k+1] { //如果左子结点小于右子结点，k指向右子结点
				k++
			}
			if v[k] > temp { //如果子节点大于父节点，将子节点值赋给父节点（不用进行交换）
				v[i] = v[k]
				i = k
			} else {
				break
			}
		}
		v[i] = temp //将temp值放到最终的位置
	}
	//从第一个非叶子结点从下至上，从右至左调整结构
	for i := len(arr)/2 - 1; i >= 0; i-- {
		heapAdjust(arr, i, len(arr))
	}
	//调整堆结构+交换堆顶元素与末尾元素
	for i := len(arr) - 1; i > 0; i-- {
		arr[0], arr[i] = arr[i], arr[0] //将堆顶元素与末尾元素进行交换
		heapAdjust(arr, 0, i)           //重新对堆进行调整
	}
}

// MergeSort 归并排序
func MergeSort(arr []int) {
	temp := []int{}
	i := 0
	for i < len(arr) {
		temp = append(temp, 0)
		i++
	}
	merdySort(arr, temp, 0, len(arr)-1)
}

func mergy(arr []int, temp []int, left, mid, right int) {
	i := left
	j := mid + 1
	t := 0
	for i <= mid && j <= right {
		if arr[i] <= arr[j] {
			temp[t] = arr[i]
			i++
		} else {
			temp[t] = arr[j]
			j++
		}
		t++
	}
	for i <= mid {
		temp[t] = arr[i]
		t++
		i++
	}
	for j <= right {
		temp[t] = arr[j]
		t++
		j++
	}
	t = 0
	for left <= right {
		arr[left] = temp[t]
		t++
		left++
	}
}

func merdySort(arr []int, temp []int, left, right int) {
	if left < right {
		mid := (left + right) / 2
		merdySort(arr, temp, left, mid)
		merdySort(arr, temp, mid+1, right)
		mergy(arr, temp, left, mid, right)
	}
}

// QuickSort 快排
func QuickSort(arr []int) {
	quickSort(arr, 0, len(arr)-1)
}

func quickSort(arr []int, left, right int) {
	if left < right {
		mid := (left + right) / 2
		if arr[left] > arr[mid] {
			arr[left], arr[mid] = arr[mid], arr[left]
		}
		if arr[left] > arr[right] {
			arr[left], arr[right] = arr[right], arr[left]
		}
		if arr[right] < arr[mid] {
			arr[right], arr[mid] = arr[mid], arr[right]
		}
		pivot := arr[mid]
		i := left
		j := right - 1
		for i < j {
			for arr[j] > pivot && i < j {
				j--
			}
			for arr[i] <= pivot && i < j {
				i++
			}
			arr[i], arr[j] = arr[j], arr[i]
		}
		quickSort(arr, left, i)
		quickSort(arr, i+1, right)
	}
}

// QuickSort2 快排
func QuickSort2(arr []int) {
	quickSort2(arr, 0, len(arr)-1)
}
func quickSort2(arr []int, left, right int) {
	if left < right {
		pivot := arr[left]
		i := left
		j := right
		for i < j {
			for arr[j] >= pivot && i < j {
				j--
			}
			if i < j {
				arr[i] = arr[j]
			}
			for arr[i] <= pivot && i < j {
				i++
			}
			if i < j {
				arr[j] = arr[i]
			}
		}
		arr[i] = pivot
		quickSort(arr, left, i-1)
		quickSort(arr, i+1, right)
	}
}

package util

import (
	"log"
	"sync"
	"testing"
)

type TestContainer struct {
}

var instance *TestContainer
var once sync.Once

func GetTestContainer() *TestContainer {
	once.Do(func() {
		log.Println("do once")
		instance = &TestContainer{}
	})
	return instance
}

func Test_HeapSort(t *testing.T) {

	tc := GetTestContainer()
	log.Println(tc)

	arr := []int{2, 5, 9, 6, 7, 2, 8, 4, 0, 11, 6, 1}
	HeapSort(arr)
	log.Println(arr)

	arr2 := []int{2, 5, 9, 6, 7, 2, 8, 4, 0, 11, 6, 1}
	MergeSort(arr2)
	log.Println(arr2)

	arr3 := []int{2, 5, 9, 6, 7, 2, 8, 4, 0, 11, 6, 1}
	QuickSort(arr3)
	log.Println(arr3)

	arr4 := []int{2, 5, 9, 6, 7, 2, 8, 4, 0, 11, 6, 1}
	QuickSort2(arr4)
	log.Println(arr4)
}

func Test_MergySort(t *testing.T) {

	tc := GetTestContainer()
	log.Println(tc)

	arr := []int{2, 5, 9, 6, 7, 2, 8, 4, 0, 11, 6, 1}
	HeapSort(arr)
	log.Println(arr)

	arr2 := []int{2, 5, 9, 6, 7, 2, 8, 4, 0, 11, 6, 1}
	MergeSort(arr2)
	log.Println(arr2)

	arr3 := []int{2, 5, 9, 6, 7, 2, 8, 4, 0, 11, 6, 1}
	QuickSort(arr3)
	log.Println(arr3)

	arr4 := []int{2, 5, 9, 6, 7, 2, 8, 4, 0, 11, 6, 1}
	QuickSort2(arr4)
	log.Println(arr4)
}

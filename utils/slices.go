package utils

import "golang.org/x/exp/constraints"

// Filter does what array.filter() does in JS.
func Filter[T any](slice []T, test func(T) bool) (ret []T) {
	for _, item := range slice {
		if test(item) {
			ret = append(ret, item)
		}
	}
	return
}

// Map does what array.map() does in JS.
func Map[T any, U any](slice []T, test func(T) U) (ret []U) {
	for _, item := range slice {
		ret = append(ret, test(item))
	}
	return
}

// RemoveDuplicates returns a new slice with duplicate items removed.
func RemoveDuplicates[T comparable](slice []T) []T {
	allKeys := make(map[T]bool)
	var list []T
	for _, item := range slice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// Reverse reverses a slice in place. Neither Go nor the exp/slices package have a method for this. Thanks!
func Reverse[T any](slice []T) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// AllSame returns true if all elements of the slice are the same, false otherwise.
func AllSame[T comparable](slice []T) bool {
	for i := 1; i < len(slice); i++ {
		if slice[i] != slice[0] {
			return false
		}
	}
	return true
}

func NumberDistribution[T constraints.Integer](slice []T) map[T]int {
	distMap := make(map[T]int)

	for _, t := range slice {
		if _, ok := distMap[t]; ok {
			// value already in map
			distMap[t]++
		} else {
			// need to put value into map
			distMap[t] = 1
		}
	}

	return distMap
}

// Sum returns the sum of all numbers in a slice.
func Sum[T constraints.Integer | constraints.Float](slice []T) T {
	var sum T

	for _, t := range slice {
		sum += t
	}

	return sum
}

// Average returns the sum of all numbers in a slice divided by the amount of numbers.
func Average[T constraints.Integer | constraints.Float](slice []T) float64 {
	sum := Sum(slice)

	return float64(sum) / float64(len(slice))
}

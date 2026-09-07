package leetcodetwosum

import "fmt"

// apenas para passar O(nˆ2)
func twoSum(nums []int, target int) []int {
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}

	return []int{0, 0}
}

// ideia propria para melhorar, mergesort O(nlogn)
func twoSum2(nums []int, target int) []int {
	ordered := sort(nums)
	for i, j := 0, len(ordered)-1; i < j; {
		if ordered[i]+ordered[j] == target {
			fmt.Println(i, j)
			firstIndex, secondIndex := -1, -1
			for k := 0; k < len(ordered); k++ {
				if nums[k] == ordered[i] && firstIndex == -1 {
					firstIndex = k
				}
				if nums[k] == ordered[j] && firstIndex != k {
					secondIndex = k
				}
			}
			return []int{firstIndex, secondIndex}
		} else if ordered[i]+ordered[j] > target {
			j--
		} else {
			i++
		}
	}
	return []int{0, 0}
}

func sort(nums []int) []int {
	length := len(nums)
	if length == 1 {
		return nums
	}
	firstHalf := nums[:length/2]
	secondHalf := nums[length/2:]

	orderedFirst := sort(firstHalf)
	orderedSecond := sort(secondHalf)

	ordered := []int{}

	i, j := 0, 0
	for i < len(orderedFirst) && j < len(orderedSecond) {
		if orderedFirst[i] <= orderedSecond[j] {
			ordered = append(ordered, orderedFirst[i])
			i++
		} else {
			ordered = append(ordered, orderedSecond[j])
			j++
		}
	}
	if i >= len(orderedFirst) {
		for j < len(orderedSecond) {
			ordered = append(ordered, orderedSecond[j])
			j++
		}
	} else {
		for i < len(orderedFirst) {
			ordered = append(ordered, orderedFirst[i])
			i++
		}
	}

	return ordered
}

// solucao otima usando hash O(n)
func twoSum3(nums []int, target int) []int {
	hash := make(map[int]int)

	for i := range len(nums) {
		if complemento, ok := hash[nums[i]]; ok {
			return []int{i, complemento}
		} else {
			hash[target-nums[i]] = i
		}
	}

	return []int{0, 0}
}

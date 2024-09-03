/*
Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.

You may assume that each input would have exactly one solution, and you may not use the same element twice.

You can return the answer in any order.
*/

package main

func Week1TwoSum(nums []int, target int) []int {
	mp := make(map[int]int)

	for i, num := range nums {
		complement := target - num
		if index, found := mp[complement]; found {
			return []int{index, i}
		}
		mp[num] = i
	}

	return nil
}

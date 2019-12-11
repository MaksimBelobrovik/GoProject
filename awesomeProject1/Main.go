package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	var str string
	fmt.Scan(&str)
	strs := strings.Split(str, "")
	mp := make(map[int]int)
	for _, s := range strs {
		num, err := strconv.Atoi(s)
		if err == nil {
			mp[num]++
		}
	}
	for k := 0; k < 9; k++ {
		maxInd := 9
		max := mp[9]
		for i := len(mp); i >= 0; i-- {
			if mp[i] > max {
				max = mp[i]
				maxInd = i
			}
		}
		for j := 1; j <= mp[maxInd]; j++ {
			fmt.Print(maxInd)
		}
		mp[maxInd] = 0
	}
}

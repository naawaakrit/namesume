package main

import "fmt"

var meanings = map[int]string{
	// 15: "ใส่ความหมายของเลข 15 ที่นี่",
	// 20: "ใส่ความหมายของเลข 20 ที่นี่",
	// 35: "ใส่ความหมายของเลข 35 ที่นี่",
}

func meaningFor(n int) string {
	if m, ok := meanings[n]; ok && m != "" {
		return m
	}
	return fmt.Sprintf("(ยังไม่มีความหมายของเลข %d ในตาราง meanings)", n)
}

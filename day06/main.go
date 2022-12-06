package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func findFirstDiff(s string, n int) int {
	for i := 0; i < len(s)-n; i++ {
		found := true
	dupcheck:
		for j := i; j < n+i-1; j++ {
			for k := j + 1; k < n+i; k++ {
				if s[j] == s[k] {
					found = false
					i = j // we can skip ahead this far
					break dupcheck
				}
			}
		}
		if found {
			return i + n
		}
	}
	return -1
}

func main() {
	f, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")
	for _, l := range lines {
		fmt.Println(findFirstDiff(l, 4), findFirstDiff(l, 14))
	}
}

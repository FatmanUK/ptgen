package main

import (
	"fmt"
)

// loop and copy nonempty - O(n)
// join and split - O(???)
func removeEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func stringFromInterface(any interface{}) string {
	return fmt.Sprintf("%v", any)
}

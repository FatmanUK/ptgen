package main

import (
	"fmt"
	"strconv"
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

func rowNumFromRowStr(rowStr string, isHex bool) (uint16, error) {
	rowBase := 10
	if isHex {
		rowBase = 16
	}
	u, err := strconv.ParseUint(rowStr, rowBase, 16)
	if err != nil {
		return 0, err
	}
	return uint16(u), nil
}

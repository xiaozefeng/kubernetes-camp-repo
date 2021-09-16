package practice_test

import (
	"fmt"
	"testing"
)

func Test_change_slice_element(t *testing.T) {
	var s = []string{"I", "am", "stupid", "and", "weak"}
	var mapping = map[string]string{
		"stupid": "smart",
		"weak":   "strong",
	}
	for i := 0; i < len(s); i++ {
		if val, ok := mapping[s[i]]; ok {
			s[i] = val
		}
	}
	fmt.Println(s)
}

package common

import "fmt"

func StringJoin(list []string, sep string) string {

	var str string

	for i, s := range list {
		if i == 0 {
			str = s
		} else {
			str = fmt.Sprintf("%s%s%s", str, sep, s)
		}
	}

	return str
}

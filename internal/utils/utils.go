package utils

import "fmt"

func FormatUrl(addr, typeM, nameM, valueM string) string {
	return fmt.Sprintf("http://%v/update/%v/%v/%v", addr, typeM, nameM, valueM)
}

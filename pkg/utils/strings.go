package utils

import "fmt"

func BlueColor(str string) string {
	return fmt.Sprintf("\033[1;36m%s\033[0m", str)
}

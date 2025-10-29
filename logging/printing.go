package logging

import (
	"fmt"
)

func printPrefix() string {
	color := string(config.Color)
	return color + bold + "[" + string(config.Role) + "]" + reset + color
}

func Println(a ...any) {
	println(printPrefix(), fmt.Sprint(a...), reset)
}

func Printf(format string, a ...any) {
	println(printPrefix() + " " + fmt.Sprintf(format, a...) + reset)
}

func PrintErrStr(a ...any) {
	println(printPrefix() + " Error: " + bold + fmt.Sprint(a...) + reset)
}

func PrintErr(err error) {
	println(printPrefix() + " Error: " + bold + err.Error() + reset)
}

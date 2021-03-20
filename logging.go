package DCP

import "fmt"

func logLn(loggerDisabled bool, a ...interface{}) {
	if !loggerDisabled {
		fmt.Println(a...)
	}
}

func logf(loggerDisabled bool, format string, a ...interface{}) {
	if !loggerDisabled {
		fmt.Printf(format, a...)
	}
}

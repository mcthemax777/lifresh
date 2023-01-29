package define

import (
	"fmt"
	"runtime"
)

const (
	OsTypeLinux   = 1
	OsTypeWindows = 2
)

var OsType int

func init() {
	os := runtime.GOOS
	switch os {
	case "windows":
		OsType = OsTypeWindows
	case "linux":
		OsType = OsTypeLinux
	default:
		panic(0)
	}

	fmt.Printf("os type - %d\n", OsType)
}

const (
	CategoryTypeMoneyMinus = 0
	CategoryTypeMoneyPlus  = 1
)

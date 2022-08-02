package help

import (
	"fmt"
	"os"
)

func Erase(path, lau string) error {
	switch lau {
	case "c++":
		path += "\\main.exe"
	case "go":
		path += "\\main.exe"
	case "java":
		path += "\\Main.class"
	default:
		return nil
	}
	fmt.Println(path)
	return os.Remove(path)
}

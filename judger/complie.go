package judger

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

func Complie(language, path string) (string, error) {
	var stderr bytes.Buffer
	if language == "go" {
		cmd := exec.Command("go", "build", "main.go")
		cmd.Dir = path
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			return path, errors.New(stderr.String())
		}
		return path, nil
	} else if language == "c++" {
		cmd := exec.Command("g++", "-o", "main.exe", "main.cpp")
		cmd.Dir = path
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			return path, errors.New(stderr.String())
		}
		return path, nil
	} else if language == "java" {
		cmd := exec.Command("javac", "Main.java")
		cmd.Stderr = &stderr
		cmd.Dir = path
		err := cmd.Run()
		if err != nil {
			fmt.Println(stderr.String())
			return path, errors.New(stderr.String())
		}
		return path, nil
	} else if language == "python" {
		return path, nil
	}
	return "", errors.New("no such language")

}

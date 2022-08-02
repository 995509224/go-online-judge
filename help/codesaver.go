package help

import "os"

func CodeSave(code, lau string) (string, error) {
	dirname := "code/" + Getuuid()
	path := dirname + "/"
	switch lau {
	case "go":
		path += "main.go"
	case "c++":
		path += "main.cpp"
	case "java":
		path += "Main.java"
	case "python":
		path += "main.py"
	}
	err := os.Mkdir(dirname, 0777)
	if err != nil {
		return "", err
	}
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	f.Write([]byte(code))
	defer f.Close()
	return dirname, nil

}

package judger

import (
	"bytes"
	"errors"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"
	"io"
	"log"
	"os/exec"
	"runtime"
	"time"
)

func Running(language, path, input, output string, maxmem, maxtime int) (string, int, error, int, int) {
	var cmd *exec.Cmd
	if language == "go" || language == "c++" {
		cmd = exec.Command("./main.exe")
	} else if language == "java" {
		cmd = exec.Command("java", "Main")
	} else if language == "python" {
		cmd = exec.Command("Python", "main.py")
	} else {
		return "", 0, errors.New("no such language"), 0, 0
	}
	cmd.Dir = path
	var out, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &out
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalln(err)
	}
	io.WriteString(stdinPipe, input)
	var flag bool = false
	var status int
	go func() {
		if err := cmd.Run(); err != nil {
			log.Println(err, stderr.String())
			flag = true
			if err.Error() == "exit status 2" {
				status = 5
				return
			}
		}
		flag = true
		time.Sleep(time.Millisecond * 100)
	}()
	var em runtime.MemStats
	var mem float32
	var time *cpu.TimesStat
	var ls float64
	for flag == false {
		if status == 5 {
			return "", 4, nil, 0, 0
		}
		if cmd.Process == nil {
			continue
		}
		p, _ := process.NewProcess(int32(cmd.Process.Pid))
		percent, _ := p.MemoryPercent()
		runtime.ReadMemStats(&em)
		mem = float32(em.Alloc / 1024 / 1024)
		mem *= percent
		time, err = p.Times()
		if err != nil {
			continue
		}
		ls = time.System + time.User
		//fmt.Println(p)
		//fmt.Println(time)
		//fmt.Println(cmd.Process)
		if mem > float32(maxmem) {
			cmd.Process.Kill()
			return "", 3, nil, int(ls * 1000), maxmem
		}
		if ls > float64(maxtime) {
			cmd.Process.Kill()
			return "", 2, nil, maxtime * 1000, int(mem)
		}

		//time.Sleep(time.Millisecond * 5)
	}
	// 答案错误
	if compare(out.String(), output) {
		status = 1
	} else {
		status = -1
	}
	return out.String(), status, nil, int(ls * 1000), int(mem)
}

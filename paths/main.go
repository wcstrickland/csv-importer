package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	buuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		fmt.Println("error:", err)
	}
	uuid := string(buuid)
	uuid = strings.TrimSuffix(uuid, "\n")

	cores := runtime.NumCPU()
	chunks := fmt.Sprintf("l/%d", cores)
	split, _ := exec.LookPath("split")
	splitCmd := &exec.Cmd{
		Path:   split,
		Args:   []string{"split", "-n", chunks, "-d", "test.txt", uuid},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}
	err = splitCmd.Run()
	if err != nil {
		fmt.Println("error:", err)
	}
	files := make([]string, cores)
	for i := 0; i < cores; i++ {
		file := fmt.Sprintf("%s0%d", uuid, i)
		files[i] = file
	}
	fmt.Println(files)
	for _, v := range files {
		if err = os.Remove(v); err != nil {
			fmt.Println(err)
		}
	}
}

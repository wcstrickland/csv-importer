package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	split, _ := exec.LookPath("split")

	splitCmd := &exec.Cmd{
		Path:   split,
		Args:   []string{"split", "-l100", "test.txt", "cus_suff"},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	fmt.Println(splitCmd.String())
	err := splitCmd.Run()
	if err != nil {
		fmt.Println("error:", err)
	}
}

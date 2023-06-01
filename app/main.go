package main

import (
	"fmt"
	"os"
	"os/exec"
	// Uncomment this block to pass the first stage!
	// "os"
	// "os/exec"
)

// Usage: your_docker.sh run <image> <command> <arg1> <arg2> ...
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage!
	//
	command := os.Args[3]
	args := os.Args[3:len(os.Args)]

	cmd := &exec.Cmd{
		Path:   command,
		Args:   args,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	// output, err := cmd.Output()
	err := cmd.Run()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			fmt.Println("is exec error")
			os.Exit(e.ExitCode())
		}
		fmt.Printf("Err: %v", err)
		os.Exit(1)
	}

	// fmt.Println(string(output))
}

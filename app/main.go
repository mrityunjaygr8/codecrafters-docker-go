package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"syscall"
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

	chrootDir, err := ioutil.TempDir("", "")
	if err != nil {
		fmt.Printf("error creating chroot dir: %v", err)
		os.Exit(1)
	}

	if err = copyExecutableIntoDir(chrootDir, command); err != nil {
		fmt.Printf("error copying executable into chroot: %v", err)
		os.Exit(1)
	}

	if err = createDevNull(chrootDir); err != nil {
		fmt.Printf("error creating /dev/null: %v", err)
		os.Exit(1)
	}

	if err = syscall.Chroot(chrootDir); err != nil {
		fmt.Printf("chroot error: %v", err)
		os.Exit(1)
	}

	cmd := &exec.Cmd{
		Path:   command,
		Args:   args,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	// output, err := cmd.Output()
	err = cmd.Run()
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

func copyExecutableIntoDir(chrootDir, command string) error {
	executablePathInChrootDir := path.Join(chrootDir, command)
	if err := os.MkdirAll(path.Dir(executablePathInChrootDir), 0750); err != nil {
		return err
	}

	return copyFile(command, executablePathInChrootDir)
}

func copyFile(sourceFilePath, destFilePath string) error {
	sourceFileStat, err := os.Stat(sourceFilePath)
	if err != nil {
		return err
	}

	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}

	defer sourceFile.Close()

	destinationFile, err := os.OpenFile(destFilePath, os.O_RDWR|os.O_CREATE, sourceFileStat.Mode())
	if err != nil {
		return err
	}

	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

func createDevNull(chrootDir string) error {
	if err := os.MkdirAll(path.Join(chrootDir, "dev"), 0750); err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(chrootDir, "dev", "null"), []byte{}, 0644)
}

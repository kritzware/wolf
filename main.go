package main

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		panic("must specify filepath")
	}

	filePath := args[0]
	fmt.Println("reading", filePath)

	latestHash, err := getFileHash(filePath)
	if err != nil {
		panic(err)
	}
	fmt.Println("initial hash:", latestHash)

	dir := getDir()
	runPath := fmt.Sprintf("%s%s", dir, filePath[1:])

	initialOutput, err := runFile(runPath)
	if err != nil {
		panic(initialOutput)
	}
	fmt.Print(initialOutput)

	for {
		changed, newHash, err := checkForChanges(latestHash, filePath)
		if err != nil {
			fmt.Println("error reading file, exiting")
			break
		}

		if changed {
			latestHash = newHash
			fmt.Println("\nfound changes, running:")
			run(filePath)
		} else {
			fmt.Println("no changes found, watching..")
			time.Sleep(time.Second)
		}
	}
}

func checkForChanges(latestHash uint64, filePath string) (bool, uint64, error) {
	hash, err := getFileHash(filePath)
	if err != nil {
		return false, 0, err
	}
	if hash != latestHash {
		return true, hash, nil
	}
	return false, 0, nil
}

func run(path string) {
	output, err := runFile(path)
	if err != nil {
		panic(output)
	}
	fmt.Print(output)
}

func runFile(path string) (string, error) {
	command := fmt.Sprintf("node %s", path)
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	cmd.Run()
	formatOutput := string(output[:])
	if err != nil {
		return fmt.Sprintf("\n%s", formatOutput), err
	}
	return formatOutput, nil
}

func getFileHash(filePath string) (uint64, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return 0, err
	}
	fileContents := string(file[:len(file)-1])
	return hash64(fileContents), nil
}

func hash64(text string) uint64 {
	algorithm := fnv.New64a()
	algorithm.Write([]byte(text))
	return algorithm.Sum64()
}

func getDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}

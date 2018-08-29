package main

import (
	"bufio"
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
	latestHash, err := getFileHash(filePath)
	if err != nil {
		panic(err)
	}

	dir := getDir()
	runPath := fmt.Sprintf("%s%s", dir, filePath[1:])
	fmt.Println(latestHash, runPath)

	clear()
	printDefault(fmt.Sprintf("Watching file %s\n\n", filePath))
	runWithOutput(runPath)

	// initialOutput, err := runFile(runPath)
	// if err != nil {
	// 	printError("An error occured, watching for changes..")
	// }
	// fmt.Print(initialOutput)

	shownMessage := false

	for {
		changed, newHash, err := checkForChanges(latestHash, filePath)
		if err != nil {
			fmt.Println("error reading file, exiting")
			continue
		}

		if changed {
			latestHash = newHash
			// run(runPath)
			runWithOutput(runPath)
			shownMessage = false
		} else {

			if !shownMessage {
				printDefault("No changes found, watching..\n")
				shownMessage = true
			}

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
		printError("An error occured, watching for changes..")
	}
	fmt.Print(output)
}

func runWithOutput(path string) {
	command := fmt.Sprintf("node %s", path)
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)
	reader, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(reader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	// if err := cmd.Wait(); err != nil {
	// 	panic(err)
	// }
}

func runFile(path string) (string, error) {
	command := fmt.Sprintf("node %s", path)
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)
	// output, err := cmd.CombinedOutput()
	// go cmd.Output()

	reader, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(reader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	if err := cmd.Start(); err != nil {
		panic(err)
	}
	if err := cmd.Wait(); err != nil {
		panic(err)
	}

	return "", nil
	// formatOutput := string(output[:])
	// if err != nil {
	// 	return fmt.Sprintf("\n%s", formatOutput), err
	// }
	// return formatOutput, nil
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

func printError(out string) {
	// clear()
	ts := formatTS()
	fmt.Printf("\x1b[31;1m[%s] %s\x1b[0m\n", ts, out)
}

func printDefault(out string) {
	// clear()
	ts := formatTS()
	fmt.Printf("\n\x1b[34;1m[%s] %s\x1b[0m", ts, out)
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func formatTS() string {
	ts := time.Now()
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		ts.Year(), ts.Month(), ts.Day(),
		ts.Hour(), ts.Minute(), ts.Second())
}

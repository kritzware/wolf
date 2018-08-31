package main

import (
	"bufio"
	"bytes"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No file to watch specified e.g. \"wolf index.js\"\n")
		os.Exit(1)
	}

	filePath := args[0]
	latestHash, err := getFileHash(filePath)
	if err != nil {
		panic(err)
	}

	dir := getDir()
	runPath := fmt.Sprintf("%s%s", dir, filePath[1:])
	logChannel := createLogChannel()

	watchingMessage := fmt.Sprintf("Running and watching for changes @%s\n\n", filePath)
	printDefault(watchingMessage, false, true)
	go runWithOutput(runPath, logChannel)

	for {
		changed, newHash, err := checkForChanges(latestHash, filePath)
		if err != nil {
			log.Fatal("Error reading file", filePath)
		}
		if changed {
			latestHash = newHash
			close(logChannel)
			logChannel = createLogChannel()
			changedMessage := fmt.Sprintf("Detected changes, watching @%s\n\n", filePath)
			printDefault(changedMessage, false, true)
			go runWithOutput(runPath, logChannel)
			// fmt.Println("changes detected, killing and creating new log channel")
		}
		time.Sleep(time.Second)
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

func runWithOutput(path string, quit chan struct{}) {
	command := fmt.Sprintf("node %s", path)
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)

	reader, err := cmd.StdoutPipe()
	if err != nil {
		panic(err.Error())
	}

	scanner := bufio.NewScanner(reader)
	go func() {
		for {
			select {
			case <-quit:
				if err := cmd.Process.Kill(); err != nil {
					log.Fatal("failed to kill process: ", err)
				}
				return
			default:
				for scanner.Scan() {
					select {
					case <-quit:
						if err := cmd.Process.Kill(); err != nil {
							log.Fatal("failed to kill process: ", err)
						}
						return
					default:
						fmt.Println(scanner.Text())
					}
				}
			}
		}
	}()

	if err := cmd.Start(); err != nil {
		panic(err.Error())
	}
	// if err := cmd.Wait(); err != nil {
	// 	panic(err.Error())
	// }
	printDebug(fmt.Sprintf("pid=%d", cmd.Process.Pid))
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

func printDefault(out string, withNewLine bool, withClear bool) {
	var b bytes.Buffer

	if withClear {
		clear()
	}
	if withNewLine {
		b.WriteString("\n")
	}
	str := "\x1b[34;1m[%s] %s\x1b[0m"
	ts := formatTS()
	b.WriteString(fmt.Sprintf(str, ts, out))
	fmt.Print(b.String())
}

func printDebug(out string) {
	ts := formatTS()
	fmt.Printf("\x1b[31;1m[DEBUG][%s] %s\x1b[0m\n", ts, out)
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

func createLogChannel() chan struct{} {
	return make(chan struct{})
}

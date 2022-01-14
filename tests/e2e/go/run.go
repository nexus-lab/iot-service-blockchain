package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"
)

func pipeStd(steamCreator func() (io.ReadCloser, error), print func(line string)) io.ReadCloser {
	stream, err := steamCreator()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		reader := bufio.NewScanner(stream)
		for reader.Scan() {
			print(reader.Text())
		}
	}()

	return stream
}

func runInGoroutine(waitGroup *sync.WaitGroup, label, binary string, args ...string) {
	fmt.Printf("%s Starting %s %s\n", label, binary, strings.Join(args, " "))

	defer waitGroup.Done()

	cmd := exec.Command(binary, args...)

	print := func(line string) {
		fmt.Printf("%s %s\n", label, line)
	}

	stdout := pipeStd(cmd.StdoutPipe, print)
	stderr := pipeStd(cmd.StderrPipe, print)
	defer stdout.Close()
	defer stderr.Close()

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var wg sync.WaitGroup

	wg.Add(3)
	go runInGoroutine(&wg, "[eventhub]", "go", "run", "tests/e2e/go/eventhub/main.go")
	go runInGoroutine(&wg, "[device]  ", "go", "run", "tests/e2e/go/device/main.go")
	go runInGoroutine(&wg, "[app]     ", "go", "run", "tests/e2e/go/app/main.go")

	wg.Wait()
}

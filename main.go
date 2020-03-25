package main

// https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}

func main() {
	// cmd := exec.Command("ls", "-lah")
	// if runtime.GOOS == "windows" {
	// 	cmd = exec.Command("tasklist")
	// }

	cmdName := "ping 127.0.0.1"
	cmdArgs := strings.Fields(cmdName)

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)

	// var stdout, stderr []byte
	var errStdout, errStderr error
	stdoutPipe, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	// cmd.Wait() should be called only after we finish reading
	// from stdoutIn and stderrIn.
	// wg ensures that we finish
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		//stdout, errStdout =
		copyAndCapture(os.Stdout, stdoutPipe)
		wg.Done()
	}()

	//stderr, errStderr =
	copyAndCapture(os.Stderr, stderrIn)

	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}
	// outStr, errStr := string(stdout), string(stderr)
	// fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
}

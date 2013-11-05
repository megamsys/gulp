

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"runtime"
	"time"
)

var workers = runtime.NumCPU()

type Result struct {
	filename string
	lino     int
	line     string
}

type Job struct {
	filename string
	results  chan<- Result
}

func (job Job) Do(lineRx *regexp.Regexp) {
	file, err := os.Open(job.filename)
	if err != nil {
		log.Printf("error: %s\n", err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for lino := 1; ; lino++ {
		line, err := reader.ReadBytes('\n')
		line = bytes.TrimRight(line, "\n\r")
		if lineRx.Match(line) {
			job.results <- Result{job.filename, lino, string(line)}
		}
		if err != nil {
			if err != io.EOF {
				log.Printf("error:%d: %s\n", lino, err)
			}
			break
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // Use all the machine's cores
	log.SetFlags(0)
	var timeoutOpt *int64 = flag.Int64("timeout", 0,
		"seconds (0 means no timeout)")
	flag.Parse()
	if *timeoutOpt < 0 || *timeoutOpt > 240 {
		log.Fatalln("timeout must be in the range [0,240] seconds")
	}
	args := flag.Args()
	if len(args) < 1 {
		log.Fatalln("a regexp to match must be specified")
	}
	pattern := args[0]
	files := args[1:]
	if len(files) < 1 {
		log.Fatalln("must provide at least one filename")
	}
	if lineRx, err := regexp.Compile(pattern); err != nil {
		log.Fatalf("invalid regexp: %s\n", err)
	} else {
		var timeout int64 = 1e9 * 60 * 10 // 10 minutes!
		if *timeoutOpt != 0 {
			timeout = *timeoutOpt * 1e9
		}
		grep(timeout, lineRx)
	}
}

func grep(timeout int64, lineRx *regexp.Regexp) {
	jobs := make(chan Job, workers)
	done := make(chan struct{}, workers)

	for i := 0; i < workers; i++ {
		go doJobs(done, lineRx, jobs)
	}
	waitAndProcessResults(timeout, done)
}

func doJobs(done chan<- struct{}, lineRx *regexp.Regexp, jobs <-chan Job) {
	for job := range jobs {
		job.Do(lineRx)
	}
	done <- struct{}{}
}

func waitAndProcessResults(timeout int64, done <-chan struct{}) {
	finish := time.After(time.Duration(timeout))
	for working := workers; working > 0; {
		select { // Blocking
		case <-finish:
			fmt.Println("timed out")
			return // Time's up so finish with what results there were
		case <-done:
			working--
		}
	}
	for {
		select { // Nonblocking
		case <-finish:
			fmt.Println("timed out")
			return // Time's up so finish with what results there were
		default:
			return
		}
	}
}

func minimum(x int, ys ...int) int {
	for _, y := range ys {
		if y < x {
			x = y
		}
	}
	return x
}

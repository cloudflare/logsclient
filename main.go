package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	email = flag.String("auth.email", "", "authorization email")
	key   = flag.String("auth.key", "", "authorization key")

	baseURL = flag.String("url", "https://api.cloudflare.com/client/v4/logs", "URL for CloudFlare logs API")

	start = flag.Int64("start", -1, "the unix epoch timestamp to start downloading from")
	end   = flag.Int64("end", time.Now().Unix(), "the unix epoch timestamp to end downloading at")
	dir   = flag.String("dir", os.TempDir(), "directory to download logs in")
)

const (
	fileTimeLayout = "logs-2006_01_02-15_04_05.log.gz"
	fileInterval   = time.Minute
	checkpointFile = "checkpoint"
)

func main() {
	validateFlags()
	log.Printf("Downloading to %d", *dir)
	downloadLogs()
	saveCheckpoint()
}

// downloadLogs connects to log downloading service and saves
// logs to local disk.
func downloadLogs() {
	var (
		// Global interval to download.
		startT = time.Unix(*start, 0).UTC()
		endT   = time.Unix(*end, 0).UTC()

		// Interval to download for current file.
		s = startT
		e time.Time
	)
	for {
		if !s.Before(endT) {
			return
		}
		e = s.Add(fileInterval)
		if e.After(endT) {
			e = endT
		}
		saveLogs(s, e)
		s = s.Add(fileInterval)
	}
}

// saveLogs downloads logs for the period [s, e) and saves
// saves them to a file.
func saveLogs(s, e time.Time) {
	log.Printf("Downloading logs from %v to %v", s, e)
	u := fmt.Sprintf("%s?start=%d&end=%d", *baseURL, s.Unix(), e.Unix())
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Add("X-Auth-Email", *email)
	req.Header.Add("X-Auth-Key", *key)
	req.Header.Add("Accept-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Fatalf("Received non-2xx status code: %d", resp.StatusCode)
	}

	fname := filepath.Join(*dir, s.Format(fileTimeLayout))
	f, err := os.Create(fname)
	if err != nil {
		log.Fatalf("Failed to create log file (%s): %v", fname, err)
	}
	defer f.Close()

	io.Copy(f, resp.Body)
}

// saveCheckpoint saves the last downloaded state in a file
// to resume download from.
func saveCheckpoint() {
	fp := filepath.Join(*dir, checkpointFile)
	os.Remove(fp)
	f, err := os.Create(fp)
	if err != nil {
		log.Fatalf("Failed to create checkpoint file (%s): %v", fp, err)
	}
	defer f.Close()
	if _, err := f.WriteString(strconv.FormatInt(*end, 10)); err != nil {
		log.Fatalf("Failed to write to checkpoint file (%s): %v", fp, err)
	}
}

// validateFlags parses flags, initializes appropriate ones
// from checkpointFile, and performs some sanity checks.
func validateFlags() {
	flag.Parse()
	if len(*email) == 0 {
		log.Fatal("No auth.email provided")
	}
	if len(*key) == 0 {
		log.Fatal("No auth.key provided")
	}
	if len(*baseURL) == 0 {
		log.Fatal("No url provided")
	}
	if *start < 0 {
		fp := filepath.Join(*dir, checkpointFile)
		b, err := ioutil.ReadFile(fp)
		if err != nil {
			log.Fatalf("Failed to read checkpoint file (%s): %v", fp, err)
		}
		s, err := strconv.ParseInt(string(b), 10, 0)
		if err != nil || s < 0 {
			log.Fatalf("Corrupt checkpoint file (%s)", fp)
		}
		*start = s
	}
	if *end < 0 {
		log.Fatalf("The provided end (%d) is < 0", *end)
	}
	if !(*end > *start) {
		log.Fatalf("The provided end (%d) is not after start (%d)", *end, *start)
	}
	if _, err := os.Stat(*dir); os.IsNotExist(err) {
		log.Fatalf("The provided dir (%s) does not exist", *dir)
	}
}

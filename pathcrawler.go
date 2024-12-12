package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/rix4uni/pathcrawler/banner"
)

func main() {
	// Define flags
	onlyComplete := flag.Bool("only-complete", false, "Show only complete URLs starting with http:// or https://")
	completeURL := flag.Bool("complete-url", false, "Complete relative URLs based on the input URL")
	timeout := flag.Duration("timeout", 10*time.Second, "Timeout duration for HTTP requests")
	delay := flag.Duration("delay", 0, "Delay between requests")
	concurrency := flag.Int("concurrent", 50, "Number of concurrent requests")
	silent := flag.Bool("silent", false, "silent mode.")
	versionFlag := flag.Bool("version", false, "Print the version of the tool and exit.")
	flag.Parse()

	if *versionFlag {
		banner.PrintBanner()
		banner.PrintVersion()
		return
	}

	if !*silent {
		banner.PrintBanner()
	}

	// Read URLs from stdin
	var urls []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Create a wait group for concurrent processing
	var wg sync.WaitGroup
	urlChan := make(chan string, *concurrency)
	client := &http.Client{
		Timeout: *timeout,
	}

	// Worker function to process URLs
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for inputURL := range urlChan {
				processURL(inputURL, client, *onlyComplete, *completeURL)
				time.Sleep(*delay) // Apply delay
			}
		}()
	}

	// Send URLs to the channel
	for _, url := range urls {
		urlChan <- url
	}
	close(urlChan)

	// Wait for all workers to finish
	wg.Wait()
}

func processURL(inputURL string, client *http.Client, onlyComplete bool, completeURL bool) {
	// Make HTTP request
	resp, err := client.Get(inputURL)
	if err != nil {
		log.Printf("Error fetching %s: %v\n", inputURL, err)
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response from %s: %v\n", inputURL, err)
		return
	}

	// Prepare the grep command
	cmd := exec.Command("grep", "-oP", `["'\''"]\K[^"'\''=]+(?=["'\''"])`)

	// Set up a pipe to send data to the stdin of the grep command
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	// Set up a pipe to capture the stdout of the grep command
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	// Start the grep command
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	// Write the HTTP body content to stdin of grep
	go func() {
		_, err := stdin.Write(body)
		if err != nil {
			log.Fatal(err)
		}
		// Close stdin to signal grep that no more input will be sent
		err = stdin.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Read the output from stdout
	output, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}

	// Wait for the grep command to finish
	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}

	// Process each matched URL
	matches := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, match := range matches {
		if completeURL {
			match = completeRelativeURL(inputURL, match)
		}

		if onlyComplete && !(strings.HasPrefix(match, "http://") || strings.HasPrefix(match, "https://")) {
			continue
		}

		fmt.Println(match)
	}
}

func completeRelativeURL(inputURL, link string) string {
	parsedInput, err := url.Parse(inputURL)
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://"):
		// Case 1: Ignore URLs starting with http:// or https://
		return link

	case strings.HasPrefix(link, "//"):
		// Case 2: Add input scheme (http: or https:) to URLs starting with //
		return parsedInput.Scheme + ":" + link

	case strings.HasPrefix(link, "/"):
		// Case 3: Add input domain to URLs starting with /
		return parsedInput.Scheme + "://" + parsedInput.Host + link

	default:
		// Case 4: Add input domain + / for other cases
		return parsedInput.Scheme + "://" + parsedInput.Host + "/" + link
	}
}

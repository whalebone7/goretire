package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
	urls, err := readURLsFromInput()
	if err != nil {
		fmt.Println("Error reading URLs:", err)
		return
	}

	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36"

	for _, endpointURL := range urls {
		client := &http.Client{}
		req, err := http.NewRequest("GET", endpointURL, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}
		req.Header.Set("User-Agent", userAgent)

		response, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			continue
		}
		defer response.Body.Close()

		jsCode, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			continue
		}

		err = ioutil.WriteFile("javascript.js", jsCode, 0644)
		if err != nil {
			fmt.Println("Error writing JavaScript to file:", err)
			continue
		}

		cmd := exec.Command("retire", "--path", "javascript.js")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error running Retire.js for %s: %v\n", endpointURL, err)
			fmt.Println("Output:", string(output))
			continue
		}

		fmt.Printf("Retire.js Scan Results for %s:\n", endpointURL)

		// Extract CVEs from the output
		cves := extractCVEs(string(output))
		if len(cves) > 0 {
			fmt.Println("CVEs found:")
			for _, cve := range cves {
				fmt.Println(cve)
			}
		} else {
			fmt.Println("No CVEs found.")
		}

		err = os.Remove("javascript.js")
		if err != nil {
			fmt.Println("Error deleting JavaScript file:", err)
		}
	}
}

func extractCVEs(output string) []string {
	var cves []string
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "[*] ") {
			cves = append(cves, line[4:])
		}
	}
	return cves
}

func readURLsFromInput() ([]string, error) {
	var urls []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url != "" {
			urls = append(urls, url)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

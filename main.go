package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
        "net/http"
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
			fmt.Println("Error running Retire.js:", err)
			continue
		}

		fmt.Println("Retire.js Scan Results for", endpointURL + ":")
		fmt.Println(string(output))

		err = os.Remove("javascript.js")
		if err != nil {
			fmt.Println("Error deleting JavaScript file:", err)
		}
	}
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

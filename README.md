# README

Go Apache log parser.


## Example Data

The `access.log` file can be generated with:

```go
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	file, err := os.Create("access.log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	ips := []string{
		"192.168.1.1", "10.0.0.5", "172.16.0.10", "203.0.113.42",
		"198.51.100.8", "192.0.2.15", "192.168.100.50", "10.20.30.40",
	}

	paths := []string{
		"/", "/index.html", "/about", "/contact", "/api/users",
		"/api/posts", "/images/logo.png", "/style.css", "/script.js",
		"/products", "/products/1", "/checkout", "/login", "/admin",
	}

	methods := []string{"GET", "POST", "PUT", "DELETE", "HEAD"}
	statuses := []int{200, 201, 204, 301, 302, 400, 401, 403, 404, 500, 502, 503}
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
		"Mozilla/5.0 (X11; Linux x86_64)",
		"curl/7.64.1",
		"python-requests/2.25.1",
	}

	rand.Seed(time.Now().UnixNano())
	baseTime := time.Now().AddDate(0, 0, -7) // Start 7 days ago

	numLogs := 1000
	for range numLogs {
		ip := ips[rand.Intn(len(ips))]
		user := "-"
		if rand.Float32() > 0.8 {
			user = fmt.Sprintf("user%d", rand.Intn(100))
		}

		// Generate timestamp
		timestamp := baseTime.Add(time.Duration(rand.Intn(7*24)) * time.Hour)
		timestamp = timestamp.Add(time.Duration(rand.Intn(60)) * time.Minute)
		timeStr := timestamp.Format("02/Jan/2006:15:04:05 -0700")

		method := methods[rand.Intn(len(methods))]
		path := paths[rand.Intn(len(paths))]
		protocol := "HTTP/1.1"

		status := statuses[rand.Intn(len(statuses))]
		responseSize := rand.Intn(50000) + 100

		userAgent := userAgents[rand.Intn(len(userAgents))]
		referer := "-"
		if rand.Float32() > 0.6 {
			referer = "https://google.com"
		}

		// Format: IP - USER [TIMESTAMP] "METHOD PATH PROTOCOL" STATUS SIZE "REFERER" "USER_AGENT"
		logLine := fmt.Sprintf(
			`%s - %s [%s] "%s %s %s" %d %d "%s" "%s"`,
			ip, user, timeStr, method, path, protocol, status, responseSize, referer, userAgent,
		)

		fmt.Fprintln(file, logLine)
	}

	fmt.Printf("Generated %d log entries in access.log\n", numLogs)
}

```

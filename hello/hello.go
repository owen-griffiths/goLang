package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	filenameFlag := flag.String("tag_file", "tags.txt", "Path to file containing tags to monitor")
	flag.Parse()

	lines := readFileLines(*filenameFlag)
	fmt.Printf("%d tags loaded\n", len(lines))

	htmlForTag := checkTag(lines[0])
	processHtml(htmlForTag)
}

func readFileLines(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var result []string
	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, line)
	}
	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func checkTag(tag string) []byte {
	url := fmt.Sprintf("http://iconosquare.com/viewer.php#/tag/%s/list", tag)
	fmt.Printf("Checking '%s'\n", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received %d[b] for %s\n", len(body), tag)
	return body
}

func processHtml(html []byte) {
	htmlStr := string(html)
	fmt.Printf("Body:\n")
	fmt.Println(htmlStr)

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Doc size = %d\n", doc.Size())
	fmt.Printf("Doc data = %s\n", doc.Nodes[0])
	matches := doc.Find("div")
	fmt.Printf("Found %d matches\n", len(matches.Nodes))
}

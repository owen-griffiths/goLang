package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	const filename = "/Users/oweng/Documents/tags.txt"

	lines := readFileLines(filename)
	fmt.Printf("%d tags loaded\n", len(lines))

	checkTag(lines[0])
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

func checkTag(tag string) {
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
	fmt.Printf("Received %d[b]\n", len(body))
}

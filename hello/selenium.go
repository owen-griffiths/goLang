package main

import (
	"fmt"
	"github.com/sourcegraph/go-selenium"
)

func main() {
  ExampleFindElement()
}

func ExampleFindElement() {
  const url = "http://iconosquare.com/viewer.php#/tag/kayano/list"
  
	var webDriver selenium.WebDriver
	var err error
	caps := selenium.Capabilities(map[string]interface{}{"browserName": "firefox"})
	if webDriver, err = selenium.NewRemote(caps, "http://localhost:4444/wd/hub"); err != nil {
		fmt.Printf("Failed to open session: %s\n", err)
		return
	}
	defer webDriver.Quit()

	err = webDriver.Get(url)
	if err != nil {
		fmt.Printf("Failed to load page: %s\n", err)
		return
	}

	if title, err := webDriver.Title(); err == nil {
		fmt.Printf("Page title: %s\n", title)
	} else {
		fmt.Printf("Failed to get page title: %s", err)
		return
	}

	var elems []selenium.WebElement
	elems, err = webDriver.FindElements(selenium.ByCSSSelector, ".viewphoto")
	if err != nil {
		fmt.Printf("Failed to find element: %s\n", err)
		return
	}

  fmt.Printf("Found %d elements", len(elems))

	// output:
	// Page title: sourcegraph/go-selenium · GitHub
	// Repository: go-selenium
}

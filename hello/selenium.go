package main

import (
	"flag"
	"fmt"
	"github.com/sourcegraph/go-selenium"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

func main() {
	imageStoreFlag := flag.String("image_store", "images", "Path to output downloaded images to")
	flag.Parse()

	ExampleFindElement(*imageStoreFlag)
}

func ExampleFindElement(imageStore string) {
	const url = "http://iconosquare.com/viewer.php#/tag/kayano/list"

	var webDriver selenium.WebDriver
	var err error
	caps := selenium.Capabilities(map[string]interface{}{"browserName": "firefox"})
	if webDriver, err = selenium.NewRemote(caps, "http://localhost:4444/wd/hub"); err != nil {
		fmt.Printf("Failed to open session: %s\n", err)
		return
	}
	defer webDriver.Quit()
	fmt.Printf("CurrentWindowHandler = %s\n", formatResult(webDriver.CurrentWindowHandle))

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

	output, err := os.Create("Images.txt")
	if err != nil {
		panic(err)
	}
	defer output.Close()

	fmt.Printf("Found %d elements", len(elems))
	for i, e := range elems {
		fmt.Printf("Parsing Entry %d\n", i)
		imageDetails := parsePhoto(e)
		fmt.Fprintf(output, "%s\n", imageDetails.String())
		imageDetails.DownloandImage(imageStore)

		fmt.Printf("==============\n\n")
	}
}

func parsePhoto(parent selenium.WebElement) imageDetails {
	user := getElementText(parent, ".list-username-user")

	id, err := parent.GetAttribute("data-id")
	if err != nil {
		panic(err)
	}

	imageSource := ""
	image, err := parent.FindElement(selenium.ByCSSSelector, ".bloc-photo img")
	if err != nil {
		fmt.Printf("Failed to find image %s\n", err)
	} else {
		fmt.Printf("Image Src: %s\n", formatAttribute(image, "src"))
		fmt.Printf("Image Id: %s\n", formatAttribute(image, "image-id"))
		fmt.Printf("Image Original: %s\n", formatAttribute(image, "data-original"))
		fmt.Printf("Image Tag: %s\n", formatResult(image.TagName))
		fmt.Printf("Image Text: %s\n", formatResult(image.Text))

		imageSource = attributeOrEmpty(image, "data-original")
	}

	var tagValues []string
	tags, err := parent.QAll(".detail-tags-droite .unTag")
	if err != nil {
		fmt.Printf("Failed to get tag elements: %s\n", err)
	} else {
		fmt.Printf("%d tags found\n", len(tags))
		for _, tagElement := range tags {
			value, err := tagElement.Text()
			if err != nil {
				fmt.Printf("Error getting tag text: %s\n", err)
			} else {
				tagValues = append(tagValues, value)
			}
		}
	}

	return imageDetails{id, user, imageSource, tagValues}
}

func getElementText(parent selenium.WebElement, selector string) string {
	child, err := parent.Q(selector)
	if err != nil {
		fmt.Printf("Error finding child element %s: %s\n", selector, err)
		return ""
	}

	value, err := child.Text()
	if err != nil {
		fmt.Printf("Error getting text of child element: %s\n", err)
		return ""
	}

	return value
}

type getStringOrError func() (string, error)

func formatResult(fn getStringOrError) string {
	s, err := fn()
	if err != nil {
		return "Error: " + err.Error()
	} else {
		return "Value: " + s
	}
}

func getSrcAttribute(e selenium.WebElement) getStringOrError {
	return func() (string, error) {
		return e.GetAttribute("src")
	}
}

func formatAttribute(e selenium.WebElement, attribute string) string {
	s, err := e.GetAttribute(attribute)
	if err != nil {
		return "Error: " + err.Error()
	} else {
		return "Value: " + s
	}
}

func attributeOrEmpty(e selenium.WebElement, attributeName string) string {
	s, err := e.GetAttribute(attributeName)
	if err != nil {
		fmt.Printf("Error accessing attribute %s: %s\n", attributeName, err)
		return ""
	} else {
		return s
	}
}

type imageDetails struct {
	id string
	user string
	url  string
	tags []string
}

func (i imageDetails) String() string {
	return i.user + "," + i.url + "," + strings.Join(i.tags, ";")
}

func (i imageDetails) DownloandImage(directory string) {
	outputPath := path.Join(directory, i.id + ".jpg")
	output, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	resp, err := http.Get(i.url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	countBytesWritten, err := io.Copy(output, resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Downloaded %d[b] for from %s", countBytesWritten, i.url)
}

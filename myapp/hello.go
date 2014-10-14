package hello

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"mime"
	"net/http"

	"appengine"
	"appengine/user"
)

func init() {
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/image", imageHandler)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	} else {
		fmt.Fprintln(w, "Hello to ", u)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
		fmt.Fprintf(w, "%d [b] posted\n", len(body))
	} else {
		fmt.Fprintln(w, "Error reading post data: ", err)
	}

	rd := bytes.NewReader(body)
	textRd := bufio.NewReader(rd)

	for {
		line, err := textRd.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Fprint(w, line)
	}
	fmt.Fprintln(w, "----")
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	img := image.NewRGBA(image.Rect(0, 0, 640, 480))
	blue := color.RGBA{0, 0, 255, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

	red := color.RGBA{255, 0, 0, 255}
	innerBounds := image.Rect(40, 80, 600, 400)
	draw.Draw(img, innerBounds, &image.Uniform{red}, image.ZP, draw.Src)

	ctype := mime.TypeByExtension("png")
	png.Encode(w, img)

	w.Header().Set("Content-Type", ctype)
}

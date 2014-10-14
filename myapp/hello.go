package hello

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"appengine"
	"appengine/user"
)

func init() {
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/post", postHandler)
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

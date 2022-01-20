package main

import (
	_ "embed"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

// Todo:
// 1. image lazy download
// 2. github action
// 3. use go tpl

const (
	addr    = ":18989"
	magicV  = "magic_v"
	tplMark = "#CONTENT#"
)

var (
	//go:embed template.html
	tpl           string
	imageDir      = ""
	imageSuffixes = []string{"jpg", "jpeg", "webp", "png", "gif", "bmp"}
	videoSuffixes = []string{"mp4", "webm", "ts", "wmv", "mkv", "avi", "mts"}
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Expect an absolute picture path.")
	}
	imageDir = os.Args[1]

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		v := regexp.MustCompile(fmt.Sprintf(`^/%s/.*$`, magicV))

		switch {
		case v.MatchString(p):
			_, _ = w.Write([]byte(render(fmt.Sprintf(`<video src="/%s" width="100%%" controls></video>`, p[len(magicV)+2:]))))
		default:
			handle(w, r)
		}
	})

	_, _ = io.WriteString(os.Stdout, fmt.Sprintf("Start listening at: %s\n", addr))
	log.Fatal(http.ListenAndServe(addr, nil))
}

func isKindOf(name string, suffixes []string) bool {
	var fn string
	if strings.HasSuffix(name, "/") {
		fn = name[:len(name)-2]
	} else {
		fn = name
	}
	for _, suffix := range suffixes {
		if strings.HasSuffix(fn, suffix) || strings.HasSuffix(fn, strings.ToUpper(suffix)) {
			return true
		}
	}
	return false
}

var isImage = func(name string) bool { return isKindOf(name, imageSuffixes) }

var isVideo = func(name string) bool { return isKindOf(name, videoSuffixes) }

func handle(w http.ResponseWriter, r *http.Request) {
	reqPath := r.URL.Path

	absolutePath := path.Join(imageDir, reqPath)

	_, _ = io.WriteString(os.Stdout, fmt.Sprintf("Request: %s\n", absolutePath))

	info, err := os.Stat(absolutePath)
	checkErr(err, w)

	if info.Mode().IsRegular() {
		checkErr(err, w)
		// w.Header().Set("Content-Encoding", "gzip")
		w.Header().Add("Cache-Control", "public, max-age=31536000, immutable")
		w.Header().Add("Last-Modified", http.TimeFormat)

		if match := r.Header.Get("If-None-Match"); match != "" {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		if isVideo(info.Name()) {
			fmt.Println("Serve video")
			http.ServeFile(w, r, absolutePath)
		} else {
			// gw, err := gzip.NewWriterLevel(w, gzip.BestCompression)
			// checkErr(err, w)
			// defer gw.Close()
			b, err := os.ReadFile(absolutePath)
			checkErr(err, w)
			_, _ = w.Write(b)
		}

	} else {

		files, err := ioutil.ReadDir(absolutePath)
		checkErr(err, w)

		var content []string

		for _, e := range files {
			// ignore dotfile
			if strings.HasPrefix(e.Name(), ".") {
				continue
			}

			var filePath = reqPath

			if reqPath == "/" {
				filePath = ""
			}

			filePath = filePath + "/" + e.Name()

			if e.IsDir() {
				content = append(content, fmt.Sprintf(`<p><a href="%s">%s</a></p>`, url.PathEscape(filePath), e.Name()))
			} else if isImage(filePath) {
				content = append(content, fmt.Sprintf(`<img data-src="/%s"/>`, url.PathEscape(filePath[1:])))
			} else if isVideo(filePath) {
				content = append(content, fmt.Sprintf(`<p><a href="/%s/%s">%s</a></p>`, magicV, url.PathEscape(filePath[1:]), e.Name()))
			}
		}

		sort.Slice(content, func(i, j int) bool {
			const link = "<a"

			a, b := content[i], content[j]
			ac, bc := strings.Contains(a, link), strings.Contains(b, link)
			if ac == bc {
				return strings.Compare(a, b) < 0
			}
			if ac {
				return true
			}
			if bc {
				return false
			}
			// dead code
			return true
		})

		_, err = io.WriteString(w, render(strings.Join(content, "")))
		checkErr(err, w)
	}
}

func render(content string) string {
	return strings.Replace(tpl, tplMark, content, 1)
}

func checkErr(err error, w http.ResponseWriter) {
	if err != nil {
		fmt.Println(err.Error())
		_, _ = w.Write([]byte(err.Error()))
		panic(http.ErrAbortHandler)
	}
}

package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

//easyjson:json
type JSONData struct {
	Browsers []string
	Company  string
	Country  string
	Email    string
	Job      string
	Name     string
	Phone    string
}

func main() {

	/*fSlow, err := os.Create("/tmp/SlowSearch.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fSlow.Close()
	fs := bufio.NewWriter(fSlow) */
	fFast, err := os.Create("/tmp/FastSearch.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fFast.Close()
	ff := bufio.NewWriter(fFast)
	SlowSearch(os.Stdout)
	//	fs.Flush()
	FastSearch(ff)
	//	ff.Flush()
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	var dataPool = sync.Pool{
		New: func() interface{} {
			return new(JSONData)
		},
	}

	const numUsers = 1000 //кол-во юзеров

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	seenBrowsers := make([]string, 0, 200)
	//	seenBrowsers = append(seenBrowsers, "")

	uniqueBrowsers := 0
	//foundUser := ""
	foundusersBs := make([]byte, numUsers*5)
	foundusersBl := 0
	foundusersBl += copy(foundusersBs[foundusersBl:], "found users:\n")

	i := 0
	for scanner.Scan() {
		user := dataPool.Get().(*JSONData) //&JSONData{}
		err := user.UnmarshalJSON(scanner.Bytes())
		if err != nil {
			//panic(err)
			log.Fatal(err)
		}

		userHasAndroid := false
		userHasMSIE := false

		for _, browser := range user.Browsers {
			isAndroid := strings.Contains(browser, "Android")
			isMSIE := strings.Contains(browser, "MSIE")
			if isAndroid {
				userHasAndroid = true
			}
			if isMSIE {
				userHasMSIE = true
			}
			if isAndroid || isMSIE {
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
						break
					}
				}
				if notSeenBefore {
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}
		if userHasAndroid && userHasMSIE {
			email := strings.Replace(user.Email, "@", " [at] ", 1)
			//	foundUser = fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email)
			foundusersBl += copy(foundusersBs[foundusersBl:], "[")
			foundusersBl += copy(foundusersBs[foundusersBl:], strconv.Itoa(i))
			foundusersBl += copy(foundusersBs[foundusersBl:], "]")
			foundusersBl += copy(foundusersBs[foundusersBl:], " ")
			foundusersBl += copy(foundusersBs[foundusersBl:], user.Name)
			foundusersBl += copy(foundusersBs[foundusersBl:], " ")
			foundusersBl += copy(foundusersBs[foundusersBl:], "<")
			foundusersBl += copy(foundusersBs[foundusersBl:], email)
			foundusersBl += copy(foundusersBs[foundusersBl:], ">\n")
			//	foundusersBl += copy(foundusersBs[foundusersBl:], foundUser)
		}
		i++
		dataPool.Put(user)
	}

	//	fmt.Fprintln(out, "found users:")
	//	fmt.Fprintln(out, string(foundusersBs))
	foundusersBl += copy(foundusersBs[foundusersBl:], "\nTotal unique browsers ")
	foundusersBl += copy(foundusersBs[foundusersBl:], strconv.Itoa(uniqueBrowsers))
	foundusersBl += copy(foundusersBs[foundusersBl:], "\n")
	//fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
	out.Write(foundusersBs[:foundusersBl])
}

package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

var addr = flag.String("addr", ":8080", "http service address") // Q=17, R=18

var templ = template.Must(template.New("qr").Parse(templateStr))

func main() {
	flag.Parse()
	http.Handle("/", http.HandlerFunc(QR))
	http.Handle("/health", http.HandlerFunc(health))
	http.Handle("/app", http.HandlerFunc(app))
	fmt.Printf("start web server %s\n", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func QR(w http.ResponseWriter, req *http.Request) {
	templ.Execute(w, req.FormValue("s"))
}

func getHost() string {
	host, err := os.Hostname()
	if err != nil {
		return fmt.Sprintf("Error get host: %v", err.Error())
	}
	return host
}

func health(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(fmt.Sprintf("version 2 OK %s", getHost())))
}

func app(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("APPLICATION version 2" + getHost()))
}

const templateStr = `
<html>
<head>
<title>QR Link Generator</title>
</head>
<body>
{{if .}}
<img src="http://chart.apis.google.com/chart?chs=300x300&cht=qr&choe=UTF-8&chl={{.}}" />
<br>
{{.}}
<br>
<br>
{{end}}
<form action="/" name=f method="GET">
    <input maxLength=1024 size=70 name=s value="" title="Text to QR Encode">
    <input type=submit value="Show QR" name=qr>
</form>
</body>
</html>
`

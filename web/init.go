package web

import (
	"fmt"
	"net/http"

	"github.com/magicmonkey/cnc/gamepad/gcode"
)

func Initialise() {
	fmt.Println("Starting webserver")
	go func() {
		http.HandleFunc("/", handleIndex)
		http.HandleFunc("/files", handleFiles)
		http.HandleFunc("/refresh-files", handleRefreshFiles)
		http.HandleFunc("/run-file", handleRunFile)
		panic(http.ListenAndServe(":8080", nil))
	}()
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<html>`))
	w.Write([]byte(`<body>`))
	w.Write([]byte(`<a href="/files">Files listing</a>`))
	w.Write([]byte(`</body>`))
	w.Write([]byte(`</html>`))
}

func handleRefreshFiles(w http.ResponseWriter, r *http.Request) {
	gcode.SendGcode("M20 S2")
	w.Write([]byte(`<html>`))
	w.Write([]byte(`<body>`))
	w.Write([]byte(`<p>Refreshing files list...</p>`))
	w.Write([]byte(`<script>window.setTimeout(function() {window.location="/files"}, 1000);</script>`))
	w.Write([]byte(`</body>`))
	w.Write([]byte(`</html>`))

}

func handleRunFile(w http.ResponseWriter, r *http.Request) {
	fname := r.URL.Query().Get("filename")
	gcode.SendGcode("M23 " + fname)
	w.Write([]byte(`<html>`))
	w.Write([]byte(`<body>`))
	w.Write([]byte(`<p>Running ` + fname + `...</p>`))
	w.Write([]byte(`<script>window.setTimeout(function() {window.location="/files"}, 1000);</script>`))
	w.Write([]byte(`</body>`))
	w.Write([]byte(`</html>`))

}

func handleFiles(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<html>`))
	w.Write([]byte(`<body>`))
	w.Write([]byte(`<h1>CNC</h1>`))
	w.Write([]byte(`<button onclick="window.location='/refresh-files'">Refresh files</button>`))
	for i, fs := range gcode.FilesList {
		w.Write([]byte(`<h2>`))
		w.Write([]byte(i))
		w.Write([]byte(`</h2>`))
		w.Write([]byte(`<ul>`))
		for _, f := range fs {
			if f[0:1] == "*" {
				continue
			}
			w.Write([]byte(`<li>`))
			w.Write([]byte(fmt.Sprintf(`<button onclick="window.location='/run-file?filename=%s/%s'">Go</button>`, i, f)))
			w.Write([]byte(f))
			w.Write([]byte(`</li>`))
		}
		w.Write([]byte(`</ul>`))
	}
	w.Write([]byte(`</body>`))
	w.Write([]byte(`</html>`))
}

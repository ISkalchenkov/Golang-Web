package handlers

import (
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/static/html/index.html")
}

var Static = http.StripPrefix(
	"/static/",
	http.FileServer(http.Dir("./web/static")),
)

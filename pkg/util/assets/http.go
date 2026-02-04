package assets

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type StoreHandler struct {
	Store
}

func (sh StoreHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var filename = strings.TrimPrefix(req.URL.Path, "/")

	if filename == "" {
		filename = "index.html"
	}

	if r, err := sh.Store.Open(filename); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GET /%s 404\n", filename)
		w.WriteHeader(http.StatusNotFound)
	} else {
		_, _ = fmt.Fprintf(os.Stdout, "GET /%s 200\n", filename)
		w.WriteHeader(http.StatusOK)
		_, _ = io.Copy(w, r)
	}
}

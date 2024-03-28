package handler

import (
	"io"
	"marketplace/util"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join("static", "index.html")
	file, err := os.Open(filePath)
	if err != nil {
		util.SendJSONError(w, r, "Could not read file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		util.SendJSONError(w, r, "Could not read file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	w.WriteHeader(http.StatusOK)

	_, err = io.Copy(w, file)
	if err != nil {
		util.SendJSONError(w, r, "Could not write response", http.StatusInternalServerError)
	}
}

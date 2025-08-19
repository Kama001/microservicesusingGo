package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) HelloHandler(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello handler invoked")
	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Failed to read request body", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(rw, "received %s", d)
}

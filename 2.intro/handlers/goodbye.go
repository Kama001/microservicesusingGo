package handlers

import (
	"log"
	"net/http"
)

type Goodbye struct {
	l *log.Logger
}

func NewGoodbye(l *log.Logger) *Goodbye {
	return &Goodbye{l}
}

func (h *Goodbye) GoodbyeHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Goodbye"))
}

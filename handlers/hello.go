package handlers

import (
	"fmt"
	"io"
	"net/http"
)

type Hello struct {
}

func ResponseFunc() *Hello {
	return &Hello{}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("Hello from home page.")
	data, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(rw, "Error: %s", err)

	}
	fmt.Fprintf(rw, "Hello %s\n", data)
}

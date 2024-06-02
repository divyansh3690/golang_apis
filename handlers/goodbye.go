package handlers

import (
	"fmt"
	"io"
	"net/http"
)

type Goodbye struct {
}

func ResponseFuncGoodbye() *Goodbye {
	return &Goodbye{}
}

func (gg *Goodbye) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	fmt.Println("INTO GOODBYES")
	incoming, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err, "error")
	}

	fmt.Fprintf(rw, "Byee %s!\n", incoming)

}

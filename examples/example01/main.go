package main

import (
	"log"
	"net/http"
	
	MNY "github.com/Comcast/golang-money"
	"github.com/justinas/alice"
)

func responseHandler1(w http.ResponseWriter, r *http.Request) {
	msg := "Request Received 1\n"
	w.Write([]byte(msg))
}

func main() {
	h1 := http.HandlerFunc(responseHandler1)
	s := MNY.New()
	chain := alice.New(s.Decorate).Then(h1)
	
	if err := http.ListenAndServe(":12345", chain); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
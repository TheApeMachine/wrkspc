package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

/*
sendkeys takes a request from another docker container and
uses it as a trigger to send the keyfragments through the
LastPass cli.
*/
func sendkeys(w http.ResponseWriter, req *http.Request) {
	cmd := exec.Command("lpass", "share", "useradd", "--read-only=true", "--hidden=true", "<your email>")
	msg := "keyfragments sent"

	if err := cmd.Run(); err != nil {
		// Override the messge with the error.
		msg = err.Error()
	}

	// Output response to the caller.
	fmt.Fprintf(w, msg+"\n")
}

func main() {
	// Since LastPass only provides a cli and we are trying to
	// use it more as a service, here is a quick wrapper that
	// exposes an endpoint.
	log.Println("starting LastPass control endpoint")
	http.HandleFunc("/sendkeys", sendkeys)
	log.Fatal(http.ListenAndServe(":8090", nil))
}

package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/rpc/json"
	"github.com/lavab/rex/node/service"
	"github.com/namsral/flag"
)

var (
	apiURL = flag.String("api_url", "http://127.0.0.1:6003/rpc", "URL of the API to use")
)

func main() {
	flag.Parse()

	if len(os.Args) < 2 {
		log.Fatal("Not enough args")
	}

	body, err := json.EncodeClientRequest("Rexd.Execute", &service.ExecuteArgs{
		Token:  "nyraITVmr61ALZNdf9Ye",
		Name:   "deploy_api",
		Branch: "master",
		Args:   nil,
	})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(*apiURL, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	rb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(string(rb))
}

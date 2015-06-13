package main

import (
	"log"
	"net/http"

	r "github.com/dancannon/gorethink"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/namsral/flag"

	"github.com/lavab/rex/node/service"
)

var (
	bindAddress       = flag.String("bind_address", ":6003", "Bind address of the rexd HTTP server")
	rethinkdbAddress  = flag.String("rethinkdb_address", "172.16.0.1:28015", "RethinkDB address to use")
	rethinkdbDatabase = flag.String("rethinkdb_database", "rex", "RethinkDB database to use")
)

func main() {
	flag.Parse()

	session, err := r.Connect(r.ConnectOpts{
		Address:  *rethinkdbAddress,
		Database: *rethinkdbDatabase,
	})
	if err != nil {
		log.Fatal(err)
	}

	r.DB(*rethinkdbDatabase).TableCreate("scripts").Exec(session)
	r.DB(*rethinkdbDatabase).TableCreate("tokens").Exec(session)

	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(&service.Service{
		Session: session,
	}, "Rexd")
	http.Handle("/rpc", s)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("lavab/rexd 0.1.0\n"))
	})

	log.Printf("Binding to %s", *bindAddress)
	if err := http.ListenAndServe(*bindAddress, nil); err != nil {
		log.Fatal(err)
	}
}

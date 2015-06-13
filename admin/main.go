package main

import (
	"log"
	"net/http"

	r "github.com/dancannon/gorethink"
	"github.com/lavab/rex/node/service"
	"github.com/namsral/flag"
)

var (
	bindAddress       = flag.String("bind_address", ":6004", "Bind address of the rexd HTTP server")
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

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.RequestURI == "/favicon.ico" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		} else if req.RequestURI != "/" {
			cursor, err := r.Table("scripts").Get(req.RequestURI[1:]).Run(session)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer cursor.Close()
			var result *service.Script
			if err := cursor.One(&result); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if req.Method == "POST" {
				if err := req.ParseForm(); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				name := req.Form.Get("name")
				if name == "" {
					if err := r.Table("scripts").Get(req.RequestURI[1:]).Delete().Exec(session); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					http.Redirect(w, req, "/", 302)
					return
				}

				code := req.Form.Get("code")
				if code == "" {
					http.Error(w, "No code passed", http.StatusBadRequest)
					return
				}

				changedID := false
				if result.ID != name {
					result.ID = name
					changedID = true
				}

				if result.Code != code {
					result.Code = code
				}

				if err := r.Table("scripts").Get(req.RequestURI[1:]).Update(result).Exec(session); err != nil {
					http.Error(w, err.Error(), http.StatusConflict)
					return
				}

				if changedID {
					http.Redirect(w, req, "/"+name, 302)
					return
				}
			}

			err = edit.Execute(w, result)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			if req.Method == "POST" {
				if err := req.ParseForm(); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				name := req.Form.Get("name")
				if name == "" {
					http.Error(w, "No name passed", http.StatusBadRequest)
					return
				}

				if err := r.Table("scripts").Insert(&service.Script{
					ID: name,
				}).Exec(session); err != nil {
					http.Error(w, err.Error(), http.StatusConflict)
					return
				}

				http.Redirect(w, req, "/"+name, 302)
				return
			}

			cursor, err := r.Table("scripts").Pluck("id").Run(session)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer cursor.Close()
			var result []*service.Script
			if err := cursor.All(&result); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = index.Execute(w, result)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})

	log.Printf("Binding to %s", *bindAddress)
	if err := http.ListenAndServe(*bindAddress, nil); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"net/http"
	"short/server"
	"short/short"
	"short/utils"
)

func main() {
	utils.ParseConfig("config.ini")
	short.StartShorter()
	defer short.DefaultShorter.Close()
	http.HandleFunc("/short", server.Short)
	http.HandleFunc("/long", server.Long)
	http.HandleFunc("/", server.Redirect)
	log.Println(http.ListenAndServe(utils.Conf.Host, nil))
}

// curl -X POST -H "Content-Type:application/json" -d "{\"req_url\": \"http://www.google.com\"}" http://127.0.0.1:8080/short
// {"resp_url":"http://127.0.0.1:8080/66XtVqtzTJ"}%

// curl -X POST -H "Content-Type:application/json" -d "{\"req_url\": \"http://127.0.0.1:8080/66XtVqtzTJ\"}" http://127.0.0.1:8080/long
// {"resp_url":"http://www.google.com"}

// curl -L http://127.0.0.1:8080/\?shortURL\=66XtVqtzTJ

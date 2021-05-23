package main

import (
    "fmt"
    "github.com/gorilla/mux"
    "log"
    "net/http"
)

var Server *mux.Router

func main() {
    Server = mux.NewRouter().StrictSlash(true)
    Server.HandleFunc("/", mainHandler)

    Server.HandleFunc("/api/{module}", apiHandler)
    Server.HandleFunc("/api/{module}/{param}", apiHandler)
    Server.HandleFunc("/api/{module}/{param}/{id}", apiHandler)

    fileServer := http.FileServer(http.Dir("./public/"))
    Server.
        PathPrefix("/static/").
        Handler(http.StripPrefix("/static", fileServer))

    fmt.Println("Server is listening on port 8181...")
    err := http.ListenAndServe(":8181", Server)
    if err != nil {
        log.Fatalf("[FATAL] %s", err.Error())
    }
}

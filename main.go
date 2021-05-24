package main

import (
    "cfgEditor/common"
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

    serverAddress := fmt.Sprintf("%s:%d", common.AppConfig.Server.Host, common.AppConfig.Server.Port)
    fmt.Printf("Server is listening on %s...", serverAddress)
    err := http.ListenAndServe(serverAddress, Server)
    if err != nil {
        log.Fatalf("[FATAL] %s", err.Error())
    }
}

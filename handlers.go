package main

import (
	"cfgEditor/cfg"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

type mainData struct {
	PageData
	Message string
	Days    []string
	Params  Configuration
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	var (
		data mainData
		err  error
	)
	data.Reset()

	data.Title = "Список задач"
	data.Message = "Message"
	data.Days = append(data.Days, "Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс")

	config, err := cfg.Load()
	if err != nil || config == nil {
		log.Printf("[FATAL] %s", err.Error())
	}

	data.Params, _ = parseConfig(config)

	err = buildTemplate(w, &data, "template/index.html", "template/scheduler.html")
	if err != nil {
		log.Printf("[FATAL] %s", err.Error())
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		data   []byte
		result struct {
			Message string
			Path    string
		}
	)
	vars := mux.Vars(r)
	key := vars["action"]

	result.Message = key
	result.Path = strings.ReplaceAll(r.RequestURI, "/api/", "")

	data, err = json.Marshal(result)

	_, err = w.Write(data)
	if err != nil {
		return
	}

	return
}

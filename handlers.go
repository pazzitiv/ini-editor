package main

import (
    "cfgEditor/cfg"
    "log"
    "net/http"
)

type mainData struct {
    PageData
    Message string
    Days []string
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

    err = buildTemplate(w, &data, "template/index.html", "template/scheduler-api.html")
    if err != nil {
        log.Printf("[FATAL] %s", err.Error())
    }
}

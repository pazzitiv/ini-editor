package main

import (
    "cfgEditor/common"
    "encoding/json"
    "github.com/gorilla/mux"
    "log"
    "net/http"
    "strconv"
    "strings"
)

type mainData struct {
    common.PageData
    Message string
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
    var (
        data mainData
        err  error
    )
    data.Reset()

    data.Title = "Список задач"
    data.Message = "Message"

    err = common.BuildTemplate(w, &data, "template/index.html", "template/scheduler.html")
    if err != nil {
        log.Printf("[FATAL] %s", err.Error())
    }
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
    var (
        err    error
        data   []byte
        config common.Configuration
    )

    method := r.Method

    vars := mux.Vars(r)
    module := vars["module"]
    id := vars["id"]

    switch module {
    case "schedules":
        switch method {
        case http.MethodGet:
            sched := common.Scheduler{}

            data, err = json.Marshal(common.Anonymous{
                "data": sched.List(),
            })

        case http.MethodDelete:
            log.Printf("%s", id)

        default:
            w.WriteHeader(http.StatusMethodNotAllowed)
            err = mux.ErrMethodMismatch
            log.Printf("[WARNING] %s: %s /api/%s", mux.ErrMethodMismatch.Error(), method, module)
        }

    case "dictionaries":
        switch method {
        case http.MethodGet:
            dicts := common.Dictionary{}

            data, err = json.Marshal(common.Anonymous{
                "data": dicts.List(),
            })

        case http.MethodPost:
            err = r.ParseForm()
            if err != nil {
                log.Printf("%s", err.Error())
                break
            }

            for index, item := range r.PostForm {
                switch index {
                case "day":
                    var days []string

                    config = common.GetConfiguration()
                    items := strings.Split(item[0], ",")
                    days = make([]string, 0, 7)

                    for j := 0; j < 7; j++ {
                        enabled := "0"
                        for _, i := range items {
                            ind, _ := strconv.Atoi(i)

                            if ind == j {
                                enabled = "1"
                                break
                            }
                        }
                        days = append(days, enabled)
                    }

                    config.Scheduler.Day = append(config.Scheduler.Day, days)

                    err := common.SaveConfiguration(config)
                    if err != nil {
                        return 
                    }
                }
            }

        default:
            w.WriteHeader(http.StatusMethodNotAllowed)
            err = mux.ErrMethodMismatch
            log.Printf("[WARNING] %s: %s /api/%s", mux.ErrMethodMismatch.Error(), method, module)
        }

    default:
        log.Printf("[WARNING] %s: /api/%s", mux.ErrNotFound.Error(), module)
        w.WriteHeader(http.StatusNotFound)
        err = mux.ErrNotFound
    }

    if err != nil {
        data, _ = json.Marshal(common.Anonymous{
            "error": err.Error(),
        })
    }

    _, err = w.Write(data)
}

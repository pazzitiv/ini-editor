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
    param := vars["param"]
    id := vars["id"]

    switch module {
    case "schedules":
        switch method {
        case http.MethodGet:
            if param == "toggle" {
                config = common.GetConfiguration()
                TaskId, _ := strconv.Atoi(id)

                isEnabled := config.Map.Map_id[TaskId-1][5]

                if isEnabled == "1" {
                    config.Map.Map_id[TaskId-1][5] = "0"
                } else {
                    config.Map.Map_id[TaskId-1][5] = "1"
                }

                err = common.SaveConfiguration(config)
                if err != nil {
                    log.Printf("[FATAL] %s", err.Error())
                    break
                }

                break
            }

            sched := common.Scheduler{}

            data, err = json.Marshal(common.Anonymous{
                "data": sched.List(),
            })
        case http.MethodPost:
            var (
                Task struct {
                    Day     string `json:"day"`
                    Time    string `json:"time"`
                    Sender  string `json:"sender"`
                    Subject string `json:"subject"`
                    Phone   string `json:"phone"`
                }
                Map []string
            )

            config = common.GetConfiguration()

            r.Body = http.MaxBytesReader(w, r.Body, 1048576)
            decoder := json.NewDecoder(r.Body)
            decoder.DisallowUnknownFields()

            err = decoder.Decode(&Task)

            if err != nil {
                log.Printf("[FATAL] %s", err.Error())
                break
            }

            Map = append(Map, Task.Day, Task.Time, Task.Sender, Task.Subject, Task.Phone, "1")

            config.Map.Map_id = append(config.Map.Map_id, Map)

            err = common.SaveConfiguration(config)
            if err != nil {
                log.Printf("[FATAL] %s", err.Error())
                break
            }
            w.WriteHeader(http.StatusCreated)

        case http.MethodDelete:
            var (
                tasks [][]string
            )

            if id == "" {
                id = param
            }

            config = common.GetConfiguration()
            TaskId, _ := strconv.Atoi(id)

            for tInd, task := range config.Map.Map_id {
                if tInd != (TaskId - 1) {
                    tasks = append(tasks, task)
                }
            }

            config.Map.Map_id = tasks

            err = common.SaveConfiguration(config)
            if err != nil {
                log.Printf("[FATAL] %s", err.Error())
                break
            }

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
                case "time":
                    config = common.GetConfiguration()
                    items := strings.Split(item[0], ",")

                    config.Scheduler.Time = append(config.Scheduler.Time, strings.Join(items, "-"))

                    err := common.SaveConfiguration(config)
                    if err != nil {
                        return
                    }
                case "sender":
                    config = common.GetConfiguration()
                    config.Templates.Sender = append(config.Templates.Sender, item[0])

                    err := common.SaveConfiguration(config)
                    if err != nil {
                        return
                    }
                case "subject":
                    config = common.GetConfiguration()
                    config.Templates.Subject = append(config.Templates.Subject, item[0])

                    err := common.SaveConfiguration(config)
                    if err != nil {
                        return
                    }
                case "telnum":
                    config = common.GetConfiguration()
                    config.Destinations.Tel_num = append(config.Destinations.Tel_num, item[0])

                    err := common.SaveConfiguration(config)
                    if err != nil {
                        return
                    }
                }

            }
            w.WriteHeader(http.StatusCreated)
        case http.MethodDelete:
            config = common.GetConfiguration()

            switch param {
            case "day":
                ids := r.URL.Query()

                var (
                    cntSubs int
                    days    [][]string
                )

                for idx, item := range config.Scheduler.Day {
                    var (
                        hasItem bool
                        intId   int
                    )
                    for _, rid := range ids {
                        intId, err = strconv.Atoi(rid[0])

                        if idx == (intId - 1) {
                            hasItem = true
                            cntSubs++

                            for mInd, m := range config.Map.Map_id {
                                mIntInd, err := strconv.Atoi(m[0])
                                if err != nil {
                                    log.Printf("ERROR: %v \n %v - %i \n %s", m, m[0], intId, err.Error())
                                }

                                if mIntInd >= intId {
                                    config.Map.Map_id[mInd][0] = strconv.Itoa(mIntInd - cntSubs)
                                }

                            }

                            break
                        }
                    }
                    if !hasItem {
                        days = append(days, item)
                    }
                }

                config.Scheduler.Day = days
            case "time":
                ids := r.URL.Query()

                var (
                    cntSubs int
                    time    []string
                )

                for idx, item := range config.Scheduler.Time {
                    var (
                        hasItem bool = false
                        intId   int
                    )
                    for _, rid := range ids {
                        intId, err = strconv.Atoi(rid[0])

                        if idx == (intId - 1) {
                            hasItem = true
                            cntSubs++

                            for mInd, m := range config.Map.Map_id {
                                mIntInd, err := strconv.Atoi(m[1])
                                if err != nil {
                                    log.Printf("ERROR: %v \n %v - %i \n %s", m, m[1], intId, err.Error())
                                }

                                if mIntInd >= intId {
                                    config.Map.Map_id[mInd][1] = strconv.Itoa(mIntInd - cntSubs)
                                }
                            }

                            break
                        }
                    }
                    if !hasItem {
                        time = append(time, item)
                    }
                }

                config.Scheduler.Time = time
            case "sender":
                ids := r.URL.Query()

                var (
                    cntSubs int
                    snd     []string
                )

                for idx, item := range config.Templates.Sender {
                    var (
                        hasItem bool = false
                        intId   int
                    )
                    for _, rid := range ids {
                        intId, err = strconv.Atoi(rid[0])

                        if idx == (intId - 1) {
                            hasItem = true
                            cntSubs++

                            for mInd, m := range config.Map.Map_id {
                                mIntInd, err := strconv.Atoi(m[2])
                                if err != nil {
                                    log.Printf("ERROR: %v \n %v - %i \n %s", m, m[2], intId, err.Error())
                                }

                                if mIntInd >= intId {
                                    config.Map.Map_id[mInd][2] = strconv.Itoa(mIntInd - cntSubs)
                                }
                            }

                            break
                        }
                    }
                    if !hasItem {
                        snd = append(snd, item)
                    }
                }

                config.Templates.Sender = snd
            case "subject":
                ids := r.URL.Query()

                var (
                    cntSubs int
                    sbj     []string
                )

                for idx, item := range config.Templates.Subject {
                    var (
                        hasItem bool = false
                        intId   int
                    )
                    for _, rid := range ids {
                        intId, err = strconv.Atoi(rid[0])

                        if idx == (intId - 1) {
                            hasItem = true
                            cntSubs++

                            for mInd, m := range config.Map.Map_id {
                                mIntInd, err := strconv.Atoi(m[3])
                                if err != nil {
                                    log.Printf("ERROR: %v \n %v - %i \n %s", m, m[3], intId, err.Error())
                                }

                                if mIntInd >= intId {
                                    config.Map.Map_id[mInd][3] = strconv.Itoa(mIntInd - cntSubs)
                                }
                            }

                            break
                        }
                    }
                    if !hasItem {
                        sbj = append(sbj, item)
                    }
                }

                config.Templates.Subject = sbj
            case "telnum":
                ids := r.URL.Query()

                var (
                    cntSubs int
                    nums    []string
                )

                for idx, item := range config.Destinations.Tel_num {
                    var (
                        hasItem bool = false
                        intId   int
                    )
                    for _, rid := range ids {
                        intId, err = strconv.Atoi(rid[0])

                        if idx == (intId - 1) {
                            hasItem = true
                            cntSubs++

                            for mInd, m := range config.Map.Map_id {
                                mIntInd, err := strconv.Atoi(m[4])
                                if err != nil {
                                    log.Printf("ERROR: %v \n %v - %i \n %s", m, m[4], intId, err.Error())
                                }

                                if mIntInd >= intId {
                                    config.Map.Map_id[mInd][4] = strconv.Itoa(mIntInd - cntSubs)
                                }
                            }

                            break
                        }
                    }
                    if !hasItem {
                        nums = append(nums, item)
                    }
                }

                config.Destinations.Tel_num = nums
            }

            err = common.SaveConfiguration(config)
            if err != nil {
                return
            }

        default:
            w.WriteHeader(http.StatusMethodNotAllowed)
            err = mux.ErrMethodMismatch
            log.Printf("[WARNING] %s: %s /api/%s", mux.ErrMethodMismatch.Error(), method, module)
        }

    case "system":
        switch method {
        case http.MethodGet:
            system := common.SystemOptions{}

            data, err = json.Marshal(common.Anonymous{
                "data": system.List(),
            })
        case http.MethodPost:
            var (
                Sys struct {
                    Period       string `json:"sys-period"`
                    Logging      string `json:"sys-logging"`
                    Server       string `json:"sys-server"`
                    Login        string `json:"sys-login"`
                    Pass         string `json:"sys-pass"`
                    Folder_check string `json:"sys-folder-check"`
                    Trash        string `json:"sys-trash"`
                    Host         string `json:"sys-ast-host"`
                    Port         string `json:"sys-ast-port"`
                }
            )

            config = common.GetConfiguration()

            r.Body = http.MaxBytesReader(w, r.Body, 1048576)
            decoder := json.NewDecoder(r.Body)
            decoder.DisallowUnknownFields()

            err = decoder.Decode(&Sys)

            if err != nil {
                log.Printf("[FATAL] %s \n %v", err.Error(), err)
                break
            }

            if Sys.Pass == "" {
                Sys.Pass = config.Imap.Pass
            }

            if Sys.Trash == "" {
                Sys.Trash = "0"
            }

            config.System.Period, _ = strconv.Atoi(Sys.Period)
            config.System.Logging = Sys.Logging
            config.Imap.Server = Sys.Server
            config.Imap.Login = Sys.Login
            config.Imap.Pass = Sys.Pass
            config.Imap.Folder_check = Sys.Folder_check
            config.Imap.Trash, _ = strconv.Atoi(Sys.Trash)
            config.Asterisk.Host = Sys.Host
            config.Asterisk.Port, _ = strconv.Atoi(Sys.Port)

            err = common.SaveConfiguration(config)
            if err != nil {
                log.Printf("[FATAL] %s", err.Error())
                break
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

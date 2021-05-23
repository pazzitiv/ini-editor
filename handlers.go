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
		case http.MethodDelete:
			config = common.GetConfiguration()

			switch param {
			case "day":
				ids := r.URL.Query()

				var (
					days [][]string
				)

				for idx, item := range config.Scheduler.Day {
					var (
						hasItem bool = false
						intId int
					)
					for _, rid := range ids {
						intId, err = strconv.Atoi(rid[0]);

						if idx == (intId - 1) {
							hasItem = true
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
					time []string
				)

				for idx, item := range config.Scheduler.Time {
					var (
						hasItem bool = false
						intId int
					)
					for _, rid := range ids {
						intId, err = strconv.Atoi(rid[0]);

						if idx == (intId - 1) {
							hasItem = true
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
					snd []string
				)

				for idx, item := range config.Templates.Sender {
					var (
						hasItem bool = false
						intId int
					)
					for _, rid := range ids {
						intId, err = strconv.Atoi(rid[0]);

						if idx == (intId - 1) {
							hasItem = true
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
					sbj []string
				)

				for idx, item := range config.Templates.Subject {
					var (
						hasItem bool = false
						intId int
					)
					for _, rid := range ids {
						intId, err = strconv.Atoi(rid[0]);

						if idx == (intId - 1) {
							hasItem = true
							break
						}
					}
					if !hasItem {
						sbj = append(sbj, item)
					}
				}

				config.Templates.Subject = sbj
			case "phone":
				ids := r.URL.Query()

				var (
					nums []string
				)

				for idx, item := range config.Destinations.Tel_num {
					var (
						hasItem bool = false
						intId int
					)
					for _, rid := range ids {
						intId, err = strconv.Atoi(rid[0]);

						if idx == (intId - 1) {
							hasItem = true
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

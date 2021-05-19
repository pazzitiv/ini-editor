package main

import (
    "cfgEditor/helpers"
    "fmt"
    "gopkg.in/ini.v1"
    "html/template"
    "io"
    "net/http"
    "strings"
)

func main() {

    http.HandleFunc("/", mainHandler)

    fmt.Println("Server is listening on port 8181...")
    _ = http.ListenAndServe(":8181", nil)
}

func buildTemplate(wr io.Writer, data interface{}, templateFile ...string) (err error) {
    var (
        tpl *template.Template
    )

    tpl, err = tpl.ParseFiles(append(templateFile, "template/header.html", "template/footer.html")...)
    if err != nil {
        return err
    }

    err = tpl.ExecuteTemplate(wr, "header.html", data)

    for _, f := range templateFile {
        err = tpl.ExecuteTemplate(wr, strings.Replace(f, "template/", "", -1), data)
        if err != nil {
            return err
        }
    }

    err = tpl.ExecuteTemplate(wr, "footer.html", data)
    if err != nil {
        return err
    }

    return nil
}

/***
Парсинг ini-файла в Структуру
*/
func parseConfig(cfg *ini.File) (config Configuration, err error) {
    const separator = "|"

    var (
        section   *ini.Section
        value     *ini.Key
    )

    /**
      Section "SHED"
    */
    section, err = cfg.GetSection("SHED")
    if err != nil {
        return config, err
    }

    value, err = section.GetKey("day")
    if err != nil {
        return config, err
    }
    config.Sheduler.Day = helpers.SplitMap(strings.Split(value.String(), separator))

    value, err = section.GetKey("time")
    if err != nil {
        return config, err
    }
    config.Sheduler.Time = strings.Split(value.String(), separator)

    /**
      Section "TMPL"
    */
    section, err = cfg.GetSection("TMPL")
    if err != nil {
        return config, err
    }

    value, err = section.GetKey("sender")
    if err != nil {
        return config, err
    }
    config.Templates.Sender = strings.Split(value.String(), separator)

    value, err = section.GetKey("subject")
    if err != nil {
        return config, err
    }
    config.Templates.Subject = strings.Split(value.String(), separator)

    /**
      Section "DEST"
    */
    section, err = cfg.GetSection("DEST")
    if err != nil {
        return config, err
    }

    value, err = section.GetKey("tel_num")
    if err != nil {
        return config, err
    }
    config.Destionations.Tel_num = strings.Split(value.String(), separator)

    /**
      Section "MAP"
    */
    section, err = cfg.GetSection("MAP")
    if err != nil {
        return config, err
    }

    value, err = section.GetKey("map_id")
    if err != nil {
        return config, err
    }
    config.Map.Map_id = helpers.SplitMap(strings.Split(value.String(), separator))

    /**
      Section "SYS"
    */
    section, err = cfg.GetSection("SYS")
    if err != nil {
        return config, err
    }

    value, err = section.GetKey("period")
    if err != nil {
        return config, err
    }
    config.System.Period = value.MustUint()

    value, err = section.GetKey("logging")
    if err != nil {
        return config, err
    }
    config.System.Logging = value.String()

    /**
      Section "IMAP"
    */
    section, err = cfg.GetSection("IMAP")
    if err != nil {
        return config, err
    }

    value, err = section.GetKey("server")
    if err != nil {
        return config, err
    }
    config.Imap.Server = value.String()

    value, err = section.GetKey("login")
    if err != nil {
        return config, err
    }
    config.Imap.Login = value.String()

    value, err = section.GetKey("pass")
    if err != nil {
        return config, err
    }
    config.Imap.Pass = value.String()

    value, err = section.GetKey("folder_check")
    if err != nil {
        return config, err
    }
    config.Imap.Folder_check = value.String()

    value, err = section.GetKey("trash")
    if err != nil {
        return config, err
    }
    config.Imap.Trash = value.MustUint()
    /**
      Section "ASTERISK"
    */
    section, err = cfg.GetSection("ASTERISK")
    if err != nil {
        return config, err
    }

    value, err = section.GetKey("host")
    if err != nil {
        return config, err
    }
    config.Asterisk.Host = value.String()

    value, err = section.GetKey("port")
    if err != nil {
        return config, err
    }
    config.Asterisk.Port = value.MustUint()

    return
}

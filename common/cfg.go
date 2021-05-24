package common

import (
    "github.com/joho/godotenv"
    "gopkg.in/ini.v1"
    "log"
    "os"
    "strconv"
    "strings"
)

var AppConfig struct{
    ConfigFile string
    Server struct{
        Host string
        Port int
    }
}

func init()  {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    AppConfig.ConfigFile = os.Getenv("CONFIGFILE")
    if AppConfig.ConfigFile == "" {
        log.Fatalf("Empty field CONFIGFILE in .env. Please fill this field.")
    }

    AppConfig.Server.Host = os.Getenv("HOST")
    if AppConfig.Server.Host == "" {
        log.Fatalf("Empty field HOST in .env. Please fill this field.")
    }

    AppConfig.Server.Port, _ = strconv.Atoi(os.Getenv("PORT"))
    if AppConfig.Server.Port == 0 {
        log.Fatalf("Empty field PORT in .env. Please fill this field.")
    }

    log.Printf("Config file: %s", AppConfig.ConfigFile)
}

func Load() (cfg *ini.File, err error) {
    cfg, err = ini.Load(AppConfig.ConfigFile)

    return
}

func Save(config Configuration) (err error) {
    var (
        cfg *ini.File
        value string
        values []string
    )

    cfg, err = ini.Load(AppConfig.ConfigFile)
    if err != nil {
        return err
    }

    /**
      Section "SHED"
    */

    for _, item := range config.Scheduler.Day {
        values = append(values, strings.Join(item, "-"))
    }
    value = strings.Join(values, "|")
    cfg.Section("SHED").Key("day").SetValue(value)
    value = strings.Join(config.Scheduler.Time, "|")
    cfg.Section("SHED").Key("time").SetValue(value)

    /**
      Section "TMPL"
    */
    value = strings.Join(config.Templates.Sender, "|")
    cfg.Section("TMPL").Key("sender").SetValue(value)
    value = strings.Join(config.Templates.Subject, "|")
    cfg.Section("TMPL").Key("subject").SetValue(value)

    /**
      Section "DEST"
    */
    value = strings.Join(config.Destinations.Tel_num, "|")
    cfg.Section("DEST").Key("tel_num").SetValue(value)

    /**
      Section "MAP"
    */
    values = make([]string, 0, 7)
    for _, item := range config.Map.Map_id {
        values = append(values, strings.Join(item, "-"))
    }
    value = strings.Join(values, "|")
    cfg.Section("MAP").Key("map_id").SetValue(value)

    /**
      Section "SYS"
    */
    cfg.Section("SYS").Key("period").SetValue(strconv.Itoa(int(config.System.Period)))
    cfg.Section("SYS").Key("logging").SetValue(config.System.Logging)

    /**
      Section "IMAP"
    */
    cfg.Section("IMAP").Key("server").SetValue(config.Imap.Server)
    cfg.Section("IMAP").Key("login").SetValue(config.Imap.Login)
    cfg.Section("IMAP").Key("pass").SetValue(config.Imap.Pass)
    cfg.Section("IMAP").Key("folder_check").SetValue(config.Imap.Folder_check)
    cfg.Section("IMAP").Key("trash").SetValue(strconv.Itoa(int(config.Imap.Trash)))

    /**
      Section "ASTERISK"
    */
    cfg.Section("ASTERISK").Key("host").SetValue(config.Asterisk.Host)
    cfg.Section("ASTERISK").Key("port").SetValue(strconv.Itoa(int(config.Asterisk.Port)))

    /**
    Save new config
     */
    err = cfg.SaveTo(AppConfig.ConfigFile)
    if err != nil {
        return err
    }

    return
}
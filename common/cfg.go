package common

import (
    "encoding/json"
    "gopkg.in/ini.v1"
)

func Load() (cfg *ini.File, err error) {
    cfg, err = ini.Load("conf.cfg")

    return
}

func Save(config Configuration) (err error) {
    var (
        cfg *ini.File
        value []byte
    )

    cfg, err = ini.Load("conf.cfg")
    if err != nil {
        return err
    }

    /**
      Section "SHED"
    */
    value, err = json.Marshal(config.Scheduler.Day)
    cfg.Section("SHED").Key("day").SetValue(string(value))
    value, err = json.Marshal(config.Scheduler.Time)
    cfg.Section("SHED").Key("time").SetValue(string(value))

    /**
      Section "TMPL"
    */
    value, err = json.Marshal(config.Templates.Sender)
    cfg.Section("TMPL").Key("sender").SetValue(string(value))
    value, err = json.Marshal(config.Templates.Subject)
    cfg.Section("TMPL").Key("subject").SetValue(string(value))

    /**
      Section "DEST"
    */
    value, err = json.Marshal(config.Destinations.Tel_num)
    cfg.Section("DEST").Key("tel_num").SetValue(string(value))

    /**
      Section "MAP"
    */
    value, err = json.Marshal(config.Map.Map_id)
    cfg.Section("MAP").Key("map_id").SetValue(string(value))

    /**
      Section "SYS"
    */
    value, err = json.Marshal(config.System.Period)
    cfg.Section("SYS").Key("period").SetValue(string(value))
    value, err = json.Marshal(config.System.Logging)
    cfg.Section("SYS").Key("logging").SetValue(string(value))

    /**
      Section "IMAP"
    */
    value, err = json.Marshal(config.Imap.Server)
    cfg.Section("IMAP").Key("server").SetValue(string(value))
    value, err = json.Marshal(config.Imap.Login)
    cfg.Section("IMAP").Key("login").SetValue(string(value))
    value, err = json.Marshal(config.Imap.Pass)
    cfg.Section("IMAP").Key("pass").SetValue(string(value))
    value, err = json.Marshal(config.Imap.Folder_check)
    cfg.Section("IMAP").Key("folder_check").SetValue(string(value))
    value, err = json.Marshal(config.Imap.Trash)
    cfg.Section("IMAP").Key("trash").SetValue(string(value))

    /**
      Section "ASTERISK"
    */
    value, err = json.Marshal(config.Asterisk.Host)
    cfg.Section("ASTERISK").Key("host").SetValue(string(value))
    value, err = json.Marshal(config.Asterisk.Port)
    cfg.Section("ASTERISK").Key("port").SetValue(string(value))

    /**
    Save new config
     */
    err = cfg.SaveTo("my.cfg")
    if err != nil {
        return err
    }

    return
}
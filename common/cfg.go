package common

import (
    "gopkg.in/ini.v1"
    "strconv"
    "strings"
)

func Load() (cfg *ini.File, err error) {
    cfg, err = ini.Load("conf.cfg")

    return
}

func Save(config Configuration) (err error) {
    var (
        cfg *ini.File
        value string
        values []string
    )

    cfg, err = ini.Load("conf.cfg")
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
/*    err = cfg.SaveTo("conf.cfg")
    if err != nil {
        return err
    }
*/
    return
}
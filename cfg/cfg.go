package cfg

import (
    "gopkg.in/ini.v1"
)

func Load() (cfg *ini.File, err error) {
    cfg, err = ini.Load("conf.cfg")

    return
}

func Save() (err error) {
    var (
        cfg *ini.File
    )

    cfg, err = ini.Load("conf.conf")
    if err != nil {
        return err
    }

    cfg.Section("").Key("app_mode").SetValue("production")
    err = cfg.SaveTo("my.cfg.local")
    if err != nil {
        return err
    }

    return
}
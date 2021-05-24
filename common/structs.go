package common

import (
    "gopkg.in/ini.v1"
    "strconv"
)

type PageData struct {
    Title string
    Copyrights string
}

type Anonymous map[string]interface{}

func (d *PageData) Reset()  {
    d.Title = "Config Editor"
    d.Copyrights = "pazzitiv@gmail.com"
}

type Configuration struct {
    Dictionaries
    Schedules
    GeneralOptions
}

func (c *Configuration) GetListItem(sl [][]string, ind string) []string {
    var (
        index int
    )

    index, _ = strconv.Atoi(ind)

    return sl[index - 1]
}

func (c *Configuration) GetItem(sl []string, ind string) string {
    var (
        index int
    )

    index, _ = strconv.Atoi(ind)

    if index == 0 {
        return sl[index]
    }

    return sl[index - 1]
}

type Dictionaries struct {
    Scheduler struct {
        Day [][]string
        Time []string
    }
    Templates struct {
        Sender []string
        Subject []string
    }
    Destinations struct {
        Tel_num []string
    }
}

type Dictionary struct {
    config  Configuration
    cfgFile *ini.File
}
func (d *Dictionary) List() Dictionaries {
    d.config = GetConfiguration()
    dictionary := Dictionaries{
        Scheduler: d.config.Scheduler,
        Templates: d.config.Templates,
        Destinations: d.config.Destinations,
    }
    return dictionary
}

type Schedules struct {
    Map struct{
        Map_id [][]string
    }
}

type Scheduler struct {
    config  Configuration
    cfgFile *ini.File
}
func (s *Scheduler) List() Schedules {
    s.config = GetConfiguration()
    return Schedules{
        Map: s.config.Map,
    }
}

type GeneralOptions struct {
    System struct {
        Period int
        Logging string
    }
    Imap struct {
        Server string
        Login string
        Pass string
        Folder_check string
        Trash int
    }
    Asterisk struct {
        Host string
        Port int
    }
}

type SystemOptions struct {
    config  Configuration
    cfgFile *ini.File
}

func (s *SystemOptions) List() GeneralOptions {
    s.config = GetConfiguration()
    return GeneralOptions{
        System: s.config.System,
        Imap: s.config.Imap,
        Asterisk: s.config.Asterisk,
    }
}
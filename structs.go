package main

import "strconv"

type PageData struct {
    Title string
    Copyrights string
}

func (d *PageData) Reset()  {
    d.Title = "Config Editor"
    d.Copyrights = "pazzitiv@gmail.com"
}

type Configuration struct {
    Sheduler struct {
        Day [][]string
        Time []string
    }
    Templates struct {
        Sender []string
        Subject []string
    }
    Destionations struct {
        Tel_num []string
    }
    Map struct {
        Map_id [][]string
    }
    System struct {
        Period uint
        Logging string
    }
    Imap struct {
        Server string
        Login string
        Pass string
        Folder_check string
        Trash uint
    }
    Asterisk struct {
        Host string
        Port uint
    }
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
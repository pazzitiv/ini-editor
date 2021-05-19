package helpers

import (
    "strings"
)

func SplitMap(sliceVar []string) (result [][]string)  {
    if (len(sliceVar) == 0 || cap(sliceVar) == 0) {
        return
    }

    for _, item := range sliceVar {
        result = append(result, strings.Split(item, "-"))
    }

    return
}
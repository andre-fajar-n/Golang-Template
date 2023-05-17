package main

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	if len(os.Args) < 2 {
		panic(`Please input file name`)
	}

	fileName := os.Args[1]
	fileName = strings.ReplaceAll(fileName, ` `, `_`)

	prefixTimeStamp := time.Now().UTC().UnixNano() / int64(time.Millisecond)
	fileName = strconv.Itoa(int(prefixTimeStamp)) + `_` + fileName

	fileNameUp := fileName + `.sql`

	f, err := os.Create("./internal/seed/" + fileNameUp)
	check(err)

	defer f.Close()
}

package main

import (
	"fmt"

	"github.com/oarkflow/date"
)

func main() {
	datetimeStr := "2079/10/14"
	format := "%Y/%m/%d"

	npTime, err := date.ParseNP(datetimeStr, format)
	if err != nil {
		panic(npTime)
	}
	fmt.Println(npTime.GetEnglishTime())
}

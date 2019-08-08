package main

import (
	"fmt"
	"io/ioutil"
	"regexp"

	excelize "github.com/360EntSecGroup-Skylar/excelize/v2"
)

func main() {
	file, err := ioutil.ReadFile(`./222.log`)
	if err == nil {
		mainThreadKeys := matchMainThread(string(file))
		fmt.Println(mainThreadKeys)
	}
}

func matchMainThread(file string) map[string]string {
	reg := regexp.MustCompile(`MAIN_(\S+)\]\s+time:(\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}.\d{3})`)
	// mathed := reg.FindStringSubmatch(text)
	mathes := reg.FindAllStringSubmatch(file, -1)
	result := map[string]string{}
	for _, line := range mathes {
		result[line[1]] = line[2]
		fmt.Println(line[1], line[2])
	}
	return result
}

func createExcel() {
	f := excelize.NewFile()

	f.SetCellValue("Sheet1", "A1", "Hello world.")

	err := f.SaveAs("./log-parsed.xlsx")
	if err != nil {
		fmt.Println(err)
	}
}

package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	excelize "github.com/360EntSecGroup-Skylar/excelize/v2"
)

const SheetName string = "Sheet1"

const ALPHABET string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func main() {
	file, err := ioutil.ReadFile(`./222.log`)
	if err == nil {
		content := string(file)
		mainThread := matchMainThread(content)
		for key, dateTotal := range mainThread {
			fmt.Println(key, dateTotal)
		}

		fmt.Println("-----")

		// subThread := matchSubThread(content)
		// for pid, keyDate := range subThread {
		// 	for key, date := range keyDate {
		// 		fmt.Println(pid, key, date)
		// 	}
		// }

		f := excelize.NewFile()

		// f.SetCellValue("Sheet1", "A1", "Hello world.")
		createMainThreadPart(f, mainThread)

		err := f.SaveAs("./log-parsed.xlsx")
		if err != nil {
			fmt.Println(err)
		}
	}
}

// mainThread { timeKey: { date: date, total: total } }
func matchMainThread(file string) map[string]map[string]string {
	timeKeyReg := regexp.MustCompile(`MAIN_(\w+)\]\s+time:(\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}.\d{3})`)
	timeKeyMatched := timeKeyReg.FindAllStringSubmatch(file, -1)

	timeKeyTotalLineReg := regexp.MustCompile(`.*MAIN_TIME.*\s(\S+)`)
	timeKeyTotalLineMatched := timeKeyTotalLineReg.FindStringSubmatch(file)

	timeKeyTotalReg := regexp.MustCompile(`(\w+):(\d+)`)
	timeKeyTotalMatched := timeKeyTotalReg.FindAllStringSubmatch(timeKeyTotalLineMatched[1], -1)

	mainTread := map[string]map[string]string{}
	for _, keyDate := range timeKeyMatched {
		key, date := keyDate[1], keyDate[2]
		mainTread[key] = map[string]string{"date": date}
		fmt.Println(key, mainTread[key])
	}
	fmt.Println("-----")

	for _, keyTotal := range timeKeyTotalMatched {
		key, total := keyTotal[1], keyTotal[2]
		if _, ok := mainTread[key]; ok {
			mainTread[key]["total"] = total
		} else {
			mainTread[key] = map[string]string{"total": total}
		}
	}

	return mainTread
}

// subThread { pid: { date: date, total: total } }
func matchSubThread(file string) map[string]map[string]map[string]string {
	timeKeyReg := regexp.MustCompile(`SUB_(\w+)_(\w+)\]\s+time:(\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}.\d{3})`)
	timeKeyMatched := timeKeyReg.FindAllStringSubmatch(file, -1)

	subThread := map[string]map[string]map[string]string{}
	for _, pidKeyDate := range timeKeyMatched {
		pid, key, date := pidKeyDate[1], pidKeyDate[2], pidKeyDate[3]
		// fmt.Println(pid, key, date)
		if _, pidOk := subThread[pid]; pidOk {
			if _, keyOk := subThread[pid][key]; !keyOk {
				subThread[pid][key] = map[string]string{"date": date}
			}
		} else {
			subThread[pid] = map[string]map[string]string{key: {"date": date}}
		}
	}

	timeTotalKeyLineReg := regexp.MustCompile(`.*SUB_TIME_(\d+)\S+\s(.*)`)
	timeTotalKeyLineMatched := timeTotalKeyLineReg.FindAllStringSubmatch(file, -1)

	timeKeyTotalReg := regexp.MustCompile(`(\w+):(\d+)`)
	for _, keyTotalLineMatched := range timeTotalKeyLineMatched {
		pid, keyTotalLine := keyTotalLineMatched[1], keyTotalLineMatched[2]
		timeKeyTotalMatched := timeKeyTotalReg.FindAllStringSubmatch(keyTotalLine, -1)
		for _, keyTotal := range timeKeyTotalMatched {
			key, total := keyTotal[1], keyTotal[2]
			if _, ok := subThread[pid][key]; ok {
				subThread[pid][key]["total"] = total
			}
		}
	}

	return subThread
}

func createMainThreadPart(f *excelize.File, mainThread map[string]map[string]string) {

	COL_NAME := strings.Split(ALPHABET, "")

	for row, header := range MainHeaders {
		var date, total string
		date = mainThread[header.Key]["date"]
		if date == "" {
			date = mainThread[header.Key+"Time"]["date"]
		}
		total = mainThread[header.Key]["total"]
		if total == "" {
			total = mainThread[header.Key+"Time"]["total"]
		}
		fmt.Println(row+1, header.Name, total, date)

		var cols = []string{header.Name, total, date}
		for i := 0; i < 3; i++ {
			f.SetCellValue(SheetName, COL_NAME[i]+strconv.Itoa(row+1), cols[i])
		}
	}
}

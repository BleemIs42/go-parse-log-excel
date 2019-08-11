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

		subThread := matchSubThread(content)
		for ruleId, keyDate := range subThread {
			for key, date := range keyDate {
				if ruleId == "67972" {
					fmt.Println(ruleId, key, date)
				}
			}
		}
		createTable(mainThread, subThread)

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

// subThread { ruleId: { timeKey: { date: date, total: total } } }
func matchSubThread(file string) map[string]map[string]map[string]string {
	timeKeyReg := regexp.MustCompile(`SUB_(\w+)_(\w+)\]\s+time:(\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}.\d{3})`)
	timeKeyMatched := timeKeyReg.FindAllStringSubmatch(file, -1)

	subThread := map[string]map[string]map[string]string{}
	for _, ruleIdKeyDate := range timeKeyMatched {
		ruleId, key, date := ruleIdKeyDate[1], ruleIdKeyDate[2], ruleIdKeyDate[3]
		// fmt.Println(ruleId, key, date)
		if _, ruleIdOk := subThread[ruleId]; ruleIdOk {
			if _, keyOk := subThread[ruleId][key]; !keyOk {
				subThread[ruleId][key] = map[string]string{"date": date}
			}
		} else {
			subThread[ruleId] = map[string]map[string]string{key: {"date": date}}
		}
	}

	timeTotalKeyLineReg := regexp.MustCompile(`.*SUB_TIME_(\d+)\S+\s(.*)`)
	timeTotalKeyLineMatched := timeTotalKeyLineReg.FindAllStringSubmatch(file, -1)

	timeKeyTotalReg := regexp.MustCompile(`(\w+):([\d\w]+)`)
	for _, keyTotalLineMatched := range timeTotalKeyLineMatched {
		ruleId, keyTotalLine := keyTotalLineMatched[1], keyTotalLineMatched[2]
		timeKeyTotalMatched := timeKeyTotalReg.FindAllStringSubmatch(keyTotalLine, -1)
		for _, keyTotal := range timeKeyTotalMatched {
			key, total := keyTotal[1], keyTotal[2]
			if _, ok := subThread[ruleId][key]; ok {
				subThread[ruleId][key]["total"] = total
			} else {
				subThread[ruleId][key] = map[string]string{"total": total}
			}
		}
	}

	return subThread
}

func createTable(mainThread map[string]map[string]string, subThread map[string]map[string]map[string]string) {

	f := excelize.NewFile()

	createMainTable(f, mainThread)
	createSubTable(f, subThread)

	err := f.SaveAs("./log-parsed.xlsx")
	if err != nil {
		fmt.Println(err)
	}
}

func createMainTable(f *excelize.File, mainThread map[string]map[string]string) {

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

		var cols = []string{header.Name, total, date}
		for i := 0; i < 3; i++ {
			f.SetCellValue(SheetName, COL_NAME[i]+strconv.Itoa(row+1), cols[i])
		}
	}
}

func createSubTable(f *excelize.File, subThread map[string]map[string]map[string]string) {

	COL_NAME := strings.Split(ALPHABET, "")
	row := 20

	f.SetCellValue(SheetName, COL_NAME[0]+strconv.Itoa(row), "RuleId")
	for idx, col := range SubHeaders {
		f.SetCellValue(SheetName, COL_NAME[idx+1]+strconv.Itoa(row), col)
	}

	for ruleId, _ := range subThread {
		f.SetCellValue(SheetName, COL_NAME[0]+strconv.Itoa(row+1), ruleId)
		for idx, col := range SubHeaders {
			val := subThread[ruleId][col]["total"]
			if val == "" {
				val = subThread[ruleId][col]["date"]
			} else {
				//TODO: max, min, avg, middle, more
			}
			f.SetCellValue(SheetName, COL_NAME[idx+1]+strconv.Itoa(row+1), val)
		}

		row++
	}

}

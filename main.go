package main

import (
	"fmt"

	excelize "github.com/360EntSecGroup-Skylar/excelize/v2"
)

func main() {
	f := excelize.NewFile()

	f.SetCellValue("Sheet1", "A2", "Hello world.")

	err := f.SaveAs("./log-parsed.xlsx")
	if err != nil {
		fmt.Println(err)
	}
}

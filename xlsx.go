package main

import (
	"github.com/360EntSecGroup-Skylar/excelize"
)

func getRows(f *excelize.File, sheetName string) [][]string {
	return f.GetRows(sheetName)
}

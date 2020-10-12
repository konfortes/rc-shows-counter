package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type show struct {
	Name      string
	StartDate time.Time
}

type rcRange struct {
	StartDate    time.Time
	EndDate      time.Time
	CanonicalKey string
	RC           string
}

type rangesByShow map[string][]string

func main() {
	f, err := excelize.OpenFile("a.xlsx")
	if err != nil {
		panic(err)
	}

	showRows := getRows(f, "2021 prime daily")
	rangeRows := getRows(f, "RC")

	shows := parseShows(showRows)
	ranges := parseRanges(rangeRows)

	rangesByShow := countRangesByShow(ranges, shows)

	err = outputResults(rangesByShow, ranges)
	if err != nil {
		panic("unable to write results")
	}
}

func countRangesByShow(ranges []rcRange, shows []show) rangesByShow {
	output := make(map[string][]string)

	for _, show := range shows {
		showRange, found := findShowRange(show, ranges)
		if !found {
			log.Printf("could not find rc date range for %s. show's date: %s", show.Name, show.StartDate)
			continue
		}
		output[show.Name] = append(output[show.Name], showRange)
	}

	return output
}

func findShowRange(show show, ranges []rcRange) (string, bool) {
	for _, r := range ranges {
		if inRange(show.StartDate, r.StartDate, r.EndDate) {
			return r.CanonicalKey, true
		}
	}
	return "", false
}

func inRange(t, t0, t1 time.Time) bool {
	if t.Before(t0) || t.After(t1) {
		return false
	}

	return true
}

func outputResults(rbs rangesByShow, rc []rcRange) error {
	output := prepareOutput(rbs)

	f := excelize.NewFile()
	index := f.NewSheet("results")
	f.SetActiveSheet(index)
	writeRanges(f, rc)
	if err := f.SaveAs("results.xlsx"); err != nil {
		log.Println(err)
	}

	writeOutput(f, output)

	if err := f.SaveAs("results.xlsx"); err != nil {
		log.Println(err)
	}

	log.Println(output)

	return nil
}

func prepareOutput(rbs rangesByShow) map[string]map[string]int {
	output := make(map[string]map[string]int)
	for showName, ranges := range rbs {
		countByRanges := make(map[string]int)
		for _, r := range ranges {
			countByRanges[r]++
		}
		output[showName] = countByRanges
	}
	return output
}

func writeRanges(f *excelize.File, ranges []rcRange) {
	for i, r := range ranges {
		cellCoords := fmt.Sprintf("%q1", rune('A'+i+1))
		// TODO: by order?
		f.SetCellValue("results", cellCoords, r.CanonicalKey)

	}
}

func writeOutput(f *excelize.File, o map[string]map[string]int) {
	columnsByRange := getColumnByRanges(f)

	rowIndex := 2
	for showName, ranges := range o {

		for r, count := range ranges {
			f.SetCellValue("results", fmt.Sprintf("A%s", strconv.Itoa(rowIndex)), showName)
			column := columnsByRange[r]
			if column == "" {
				log.Println("could not find range column for " + r)
				continue
			}
			cellCoords := fmt.Sprintf("%s%s", column, strconv.Itoa(rowIndex))
			f.SetCellValue("results", cellCoords, count)
		}

		rowIndex++
	}
}

func getColumnByRanges(f *excelize.File) map[string]string {
	res := make(map[string]string)
	fileRanges := f.GetRows("results")[0]

	for i, r := range fileRanges {
		if r == "" {
			continue
		}
		column := fmt.Sprintf("%q", rune('A'+i))
		res[r] = strings.Replace(column, "'", "", -1)
	}
	return res
}

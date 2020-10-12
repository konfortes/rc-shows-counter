package main

import (
	"fmt"
	"strings"
	"time"
)

const dateLayout string = "02/01/2006"

// panics on invalid parse
func parseDate(layout, date string) time.Time {
	d, err := time.Parse(layout, date)
	if err != nil {
		panic(err)
	}
	return d
}

func parseShows(rows [][]string) []show {
	shows := []show{}
	for i, row := range rows {
		if i < 1 {
			continue
		}
		name := row[0]
		startDate := parseDate(dateLayout, row[1])
		shows = append(shows, show{name, startDate})

	}
	return shows
}

func parseRanges(rows [][]string) []rcRange {
	rcs := []rcRange{}
	for i, row := range rows {

		if i == 0 || row[2] == "" {
			continue
		}

		startEnd := strings.Split(row[2], "-")
		layout := "02/01/06"
		startDate := parseDate(layout, startEnd[0])
		endDate := parseDate(layout, startEnd[1])

		rc := row[4]

		rcs = append(rcs, rcRange{startDate, endDate, fmt.Sprintf("%s-%s", startEnd[0], startEnd[1]), rc})
	}
	return rcs
}

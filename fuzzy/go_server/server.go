package main

import (
	"fmt"
	"regexp"
	"sort"
	"net/http"
	"net/url"
	"strings"
	"slices"
	"io"
	"encoding/csv"
	"os"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

func main() {

	f, _ := os.Open("../courses.csv")
	csvreader := csv.NewReader(f)
	var names []string
	var pillars []string
	var rows []string
	var instructors []string
	var course_codes []string
	var records [][]string
	headers, _ := csvreader.Read()
	fmt.Println(headers)
	for {
		record, err := csvreader.Read()
		if err == io.EOF {
			break
		}
		records = append(records, record)
		names = append(names, record[6]+" "+record[7])
		instructors = append(instructors , record[7])
		course_codes = append(course_codes, record[5])
		rows = append(rows, strings.Join(record, " "))
		pillars = append(pillars, record[0])
	}
	http.HandleFunc("/search/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("server: %s /\n", r.Method)
		fmt.Printf("server: query id: %s\n", r.URL.Query().Get("id"))
		fmt.Printf("server: content-type: %s\n", r.Header.Get("content-type"))
		fmt.Printf("server: headers:\n")
		for headerName, headerValue := range r.Header {
			fmt.Printf("\t%s = %s\n", headerName, strings.Join(headerValue, ", "))
		}

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("server: could not read request body: %s\n", err)
		}
		fmt.Printf("server: request body: %s\n", reqBody)

		pillar_names := []string {"smt", "hass", "asd", "epd", "esd", "istd", "csd", "dai"}

		msg, _ := strings.CutPrefix(string(reqBody), "search=")
		msg, _ = url.QueryUnescape(msg)
		msg = strings.Join(strings.Fields(msg), " ")

		search_arr := names

		if slices.Contains(pillar_names ,strings.ToLower(msg)) {
			if msg == "csd" {
				msg = "istd"
			}
			search_arr = pillars

		}

		matched, _ := regexp.MatchString(`[0-9][0-9]|[0-9][0-9]\.[0-9][0-9][0-9]*`, msg)
		if matched {
			search_arr = course_codes
		}
		
		found := fuzzy.RankFindFold(msg, search_arr)
		sort.Slice(found, func(i, j int) bool {
			return found[i].Distance < found[j].Distance
		})

		div_html := ""
		for _, elem := range found {
			temp := ""
			//fmt.Println(elem)
			idx := elem.OriginalIndex
			temp = fmt.Sprintf(`
				<div class="course-header">
				<button onclick="stopTimer('%s: %s');" class="enroll-btn">Enroll</button>
				<div class="course-number">%s</div>
				<div class="course-terms">(%s Term %s%s%s)</div>
				<div class="course-name">%s</div>
				<div class="course-instructor">%s</div>
				</div><br>
			` , records[idx][5], records[idx][6], records[idx][5], records[idx][0], records[idx][2], records[idx][3], records[idx][4], records[idx][6], records[idx][7])
			div_html += temp
		}

		//stars := `<img class="star" src="assets/star.svg" />`
		//fmt.Println(div_html)
		fmt.Fprintf(w, `%s`, div_html)
	})


    fs := http.FileServer(http.Dir("../static"))
    http.Handle("/", fs)
	http.ListenAndServe("localhost:8082", nil)
}

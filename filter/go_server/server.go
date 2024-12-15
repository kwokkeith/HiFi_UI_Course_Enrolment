package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"slices"
	"io"
	"encoding/csv"
	"os"
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

		tokens := strings.Split(string(reqBody),"&")
		pillar, _ := url.QueryUnescape(strings.Split(tokens[0],"=")[1])
		instructor, _ := url.QueryUnescape(strings.Split(tokens[1],"=")[1])
		term, _ := url.QueryUnescape(strings.Split(tokens[2],"=")[1])

		fmt.Println(pillar)
		fmt.Println(instructor)
		fmt.Println(term)


		div_html := ""
		for idx, elem := range records {
			currPillar := elem[0]
			currInstructor := elem[7]
			currTerm := elem[2:5]

			if (currPillar == pillar || pillar == "none") && (currInstructor == instructor || instructor == "none") && (slices.Contains(currTerm, term) || slices.Contains(currTerm, ", "+term) || term == "none") {
				temp := ""
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
		}

		//stars := `<img class="star" src="assets/star.svg" />`
		//fmt.Println(div_html)
		fmt.Fprintf(w, `%s`, div_html)
	})


    fs := http.FileServer(http.Dir("../static"))
    http.Handle("/", fs)
	http.ListenAndServe("localhost:8082", nil)
}

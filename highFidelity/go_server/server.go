package main

import (
	"fmt"
	"regexp"
	"sort"
	"net/http"
	"net/url"
	"strings"
	//"slices"
	"io"
	"encoding/csv"
	"encoding/json"
	"os"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

const lipsum = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum sed ante rutrum, tristique quam sit amet, pellentesque risus. Etiam vitae augue arcu. Nunc sed erat ipsum. Proin porttitor purus eget risus convallis faucibus. Aenean turpis turpis, luctus vitae quam ac, rhoncus luctus orci. Aliquam rhoncus ex vitae ornare mollis. Nam id enim molestie, sagittis purus vitae, ornare nibh."

func main() {
	f, _ := os.Open("../courses.csv")
	csvreader := csv.NewReader(f)
	var names []string
	var pillars []string
	var rows []string
	var instructors []string
	var course_codes []string
	var terms []string
	var records [][]string
	headers, _ := csvreader.Read()
	fmt.Println(headers)
	for {
		record, err := csvreader.Read()
		if err == io.EOF {
			break
		}
		records = append(records, record)
		names = append(names, record[3]+" "+record[4])
		instructors = append(instructors , record[4])
		course_codes = append(course_codes, record[2])
		rows = append(rows, strings.Join(record, " "))
		pillars = append(pillars, record[0])
		terms = append(terms, record[1])
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
		//fmt.Printf("server: request body: %s\n", reqBody)

		rawQuery, _ := url.QueryUnescape(string(reqBody))

		tokens := strings.Split(rawQuery,"&")
		//fmt.Println(tokens)
		search, _ := url.QueryUnescape(strings.Split(tokens[0],"=")[1])
		pillar_filter_json, _ := url.QueryUnescape(strings.Split(tokens[1],"=")[1])
		term_filter_json, _ := url.QueryUnescape(strings.Split(tokens[2],"=")[1])

		msg := strings.Join(strings.Fields(search), " ")

		search_arr := names

		matched, _ := regexp.MatchString(`[0-9][0-9]|[0-9][0-9]\.[0-9][0-9][0-9]*`, msg)
		if matched {
			search_arr = course_codes
		}
		
		var filtered_idx []int

		for idx, _ := range records {
			filtered_idx = append(filtered_idx, idx)
		}

		var pillar_filter map[string]interface{}
		var term_filter map[string]interface{}

		json.Unmarshal([]byte(pillar_filter_json), &pillar_filter)
		json.Unmarshal([]byte(term_filter_json), &term_filter)
		//fmt.Println(pillar_filter)
		//fmt.Println(term_filter)

		found := fuzzy.RankFindFold(msg, search_arr)
		sort.Slice(found, func(i, j int) bool {
			return found[i].Distance < found[j].Distance
		})

		div_html := ""
		//div_html += rawQuery + "<br>"

		checkPillar := func(idx int) bool { 
			return pillar_filter[strings.ToLower(records[idx][0])] == 1.0
		}

		checkTerm := func(idx int) bool { 
			result := false
			for _, ch := range records[idx][1] {
				result = result || (term_filter[string(ch)] == 1.0)
			}
			return result
		}

		fmtTerms := func(s string) string { 
			output := ""
			for _, ch := range s {
				output += string(ch) + ", "
			}
			return output[:len(output)-2]
		}

		all_pillars := (sumVals(pillar_filter) == 0)
		all_terms := (sumVals(term_filter) == 0)

		for _, elem := range found {
			temp := ""
			//fmt.Println(elem)
			idx := elem.OriginalIndex
			if (checkPillar(idx) || all_pillars) &&
			(checkTerm(idx) || all_terms) {
				pillar := records[idx][0]
				terms := fmtTerms(records[idx][1])
				course_code := records[idx][2]
				course_name := records[idx][3]
				instructor := records[idx][4]
				description := records[idx][5]
				if description == "" { description = lipsum }
				temp = fmt.Sprintf(`
				<div class="courseBlock">
					<button class="courseBlockButton" type="button">
						<div class="courseBlock1row">
							<span class="courseBlockCourseCode">%s</span>
							<span class="courseBlockTerm"> %s Term %s</span>
						</div>
						<p class="courseBlockCourseTitle">%s</p>
						<div class="courseBlock2row">
							<div class="courseBlock2rowSubContainer">
								<p class="courseBlockInstructor">Instructor: %s</p>
								<div class="courseBlockMessage">
									<img src="./images/courseBlock/starIcon_green.svg">
									<span>This course fulfills your track: Software Development</span>
								</div>
							</div>
							<img src="./images/courseBlock/arrowForwardIcon.svg">
						</div>
					</button>
				</div>
				`, course_code, pillar, terms, course_name, instructor)
				div_html += temp
			}
		}

		//stars := `<img class="star" src="assets/star.svg" />`
		//fmt.Println(div_html)
		fmt.Fprintf(w, `%s`, div_html)
	})


	http.ListenAndServe("localhost:9990", nil)

}

func sumVals(m map[string]interface{}) float64 { 
	sum := 0.0
	for _, v := range m {
		sum += v.(float64)
	}
	return sum
}


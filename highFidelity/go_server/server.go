package main

import (
	"fmt"
	"regexp"
	"sort"
	"net/http"
	"net/url"
	"strings"
	"math/rand"
	"strconv"
	//"slices"
	"io"
	"encoding/csv"
	"encoding/json"
	"os"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

var rating_colour_codes = [6]string {"#FB8787", "#FBAE87", "#FBDA87", "#FBF587", "#BFFB87", "#9AFB87"}


const lipsum = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum sed ante rutrum, tristique quam sit amet, pellentesque risus. Etiam vitae augue arcu. Nunc sed erat ipsum. Proin porttitor purus eget risus convallis faucibus. Aenean turpis turpis, luctus vitae quam ac, rhoncus luctus orci. Aliquam rhoncus ex vitae ornare mollis. Nam id enim molestie, sagittis purus vitae, ornare nibh."

func main() {
	f, _ := os.Open("../courses_alpha.csv")
	csvreader := csv.NewReader(f)

	tracks := make(map[string][]string)
	tracks["istd"] = []string {"Artificial Intelligence", "Data Analytics", "Financial Technology", "IoT and Intelligent Systems", "Security", "Software Engineering", "Visual Analytics and Computing"}
	tracks["epd"] = []string {"Beyond Industry 4.0", "Computer Engineering", "Electrical Engineering", "Healthcare Engineering Design", "Mechanical Engineering", "Robotics"}
	tracks["esd"] = []string {"Avation Systems", "Business Analytics and Operations Research", "Financial Services", "Supply Chain & Logistics"}
	tracks["asd"] = []string {"Architecture Track"}
	tracks["dai"] = []string {"Healthcare Design", "Enterprise Design"}
	tracks["hass"] = []string {"Digital Humanities Minor"}
	tracks["freshmore"] = []string {"Freshmore Minor"}

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
		fmt.Println(tokens)
		if(len(tokens) == 3) {
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
			if msg != "" {
				sort.Slice(found, func(i, j int) bool {
					return found[i].Distance < found[j].Distance
				})
			}

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
					pillar_tracks := tracks[strings.ToLower(pillar)] 
					fulfillment := fmt.Sprintf(`<div class="courseBlockMessage">
					<img src="./images/courseBlock/starIcon_green.svg">
					<span>This course fulfills your track: %s</span>
					</div>`, pillar_tracks[(len(course_name))%len(pillar_tracks)]);
					if description == "" {
						description = lipsum
						fulfillment = ""
					}
					temp = fmt.Sprintf(`
					<!-- Course block 1 -->
					<div class="courseFullContainer">
						<div class="courseBlock">
						<button class="courseBlockButton" type="button" 
						hx-post="/htmx/search/" hx-trigger="click" hx-target="div.rightWindow" hx-vals='{"trigger":1, "index":%d}' hx-include="">
						<div class="courseBlock1row">
						<span class="courseBlockCourseCode">%s</span>
						<span class="courseBlockTerm">%s Term %s</span>
						</div>

						<p class="courseBlockCourseTitle">%s</p>

						<div class="courseBlock2row">
						<div class="courseBlock2rowSubContainer">
						<p class="courseBlockInstructor">Instructor: %s</p>
						%s
						</div>

						<img src="./images/courseBlock/arrowForwardIcon.svg">
						</div>
						</button>
						</div>
					</div>
					`, idx, course_code, pillar, terms, course_name, instructor, fulfillment)
					div_html += temp
				}
			}

			//stars := `<img class="star" src="assets/star.svg" />`
			//fmt.Println(div_html)
			fmt.Fprintf(w, `%s`, div_html)
		} else {
			div_html := ""
			fmt.Println(tokens)
			trigger, _ := url.QueryUnescape(strings.Split(tokens[0],"=")[1])
			idx_str, _ := url.QueryUnescape(strings.Split(tokens[1],"=")[1])
			idx, _ := strconv.Atoi(idx_str)
			if trigger == "1" {
				pillar := records[idx][0]
				course_code := records[idx][2]
				course_name := records[idx][3]
				instructor := records[idx][4]
				description := records[idx][5]
				pillar_tracks := tracks[strings.ToLower(pillar)] 
				ratings := [4]int {rand.Intn(6), rand.Intn(6), rand.Intn(6), rand.Intn(6)}
				rating_colours := [4]string {rating_colour_codes[ratings[0]], rating_colour_codes[ratings[1]], rating_colour_codes[5-ratings[2]], rating_colour_codes[5-ratings[3]]}
				

				buttons := [2]string {fmt.Sprintf(`
					<div class="courseDetailTableButtonContainer">
					<button type="button" class="courseDetailTableAddCartButton" onClick="addCourse()">
					<img src="./images/courseDetail/addShoppingCart_white.svg">
					<span>Add to Cart</span>
					</button>
					<button type="button" class="courseDetailTableEnrollButton" onClick="enrollCourse('%s', '%s')">
					<img src="./images/courseDetail/addIcon.svg">
					<span>Enroll</span>
					</button>
					</div>
				`, course_code, course_name),
				`
					<div class="courseDetailTableButtonContainer">
					<button type="button" class="courseDetailTableAddCartButton" disabled>
					<img src="./images/courseDetail/addShoppingCart_white.svg">
					<span>Add to Cart</span>
					</button>
					<button type="button" class="courseDetailTableEnrollButton" disabled>
					<img src="./images/courseDetail/addIcon.svg">
					<span>Enroll</span>
					</button>
					</div>
				`}


				friends := [2]string {`
					<div class="popup">
					<span class="popuptext popuptext_friend">
					<p><b>Friends taking this class:</b></p>
					<p>Bob L</p>
					<p>Charlie S</p>
					<p>Andy</p>
					</span>
					<img src="./images/courseDetail/3friendIcon.png" class="courseDetailTableFriendIcon">               
					</div>
				`, "None :("}

				statuses := [2]string {`<td style="color: green">Open</td>`, `<td style="color: red">Closed</td>`}

				if course_code == "40.008" {
					buttons[0], buttons[1] = buttons[1], buttons[0]
					friends[0], friends[1] = friends[1], friends[0]
					statuses[0], statuses[1] = statuses[1], statuses[0]
				}

				fulfillment := fmt.Sprintf(`<div class="courseBlockMessage">
				<img src="./images/courseBlock/starIcon_green.svg">
				<span>This course fulfills your track: %s</span>
				</div>`, pillar_tracks[(len(course_name))%len(pillar_tracks)])
				if description == "" {
					description = lipsum
					fulfillment = ""
					friends[0] = "None :("
				}
				div_html += fmt.Sprintf(`
					<script src="hide_unhide.js"></script>
					<div class="courseDetail ">
					<!-- SubContainer contains content except student review button -->
					<div class="courseDetailSubContainer">
						<div class="courseDetailHeader">
						<button type="button" class="courseDetailReturnButton" 
						hx-post="/htmx/search/" hx-trigger="click" hx-target="div.rightWindow" hx-vals='js:{pillars: pillars, terms: terms}' hx-include="[name='search']">
						<img src="./images/courseDetail/undoIcon.svg"></button>
						<span class="courseDetailCourseTitle">%s: %s</span>
						</div>

					<hr>

					<p class="courseDetailCourseDescription">%s</p>

					<div class="courseDetailPrerequisitesContainer">
					<span><b>Prerequisites: </b></span>
					<span class="courseDetailPrerequisitesCourse">10.014 Computational Thinking for Design</span>
					<img src="./images/courseDetail/checkCircleIcon_darkGreen.svg">
					</div>

					<div class="courseDetailRatingContainer_All">

					<div class="courseDetailRatingContainer">
					<span>Content</span>
					<div class="popup">
					<span class="popuptext">The relevance and clarity of materials in supporting course objectives.</span>
					<img src="./images/courseDetail/infoIcon_grey.svg">
					</div>
					<div class="courseDetailRatingBoxContainer" style="background-color: %s">
					<span class="courseDetailRatingBox">%s/5</span>
					</div>
					</div>
					<div class="courseDetailRatingContainer">
					<span>Teaching</span>
					<div class="popup">
					<span class="popuptext">Instructor’s effectiveness in delivering the course material.</span>
					<img src="./images/courseDetail/infoIcon_grey.svg">
					</div>
					<div class="courseDetailRatingBoxContainer" style="background-color: %s">
					<span class="courseDetailRatingBox">%s/5</span>
					</div>
					</div>
					<div class="courseDetailRatingContainer">
					<span>Difficulty</span>
					<div class="popup">
					<span class="popuptext">How challenging students found the course in terms of understanding and completing the material.</span>
					<img src="./images/courseDetail/infoIcon_grey.svg">
					</div>
					<div class="courseDetailRatingBoxContainer" style="background-color: %s">
					<span class="courseDetailRatingBox">%s/5</span>
					</div>
					</div>
					<div class="courseDetailRatingContainer">
					<span>Workload</span>
					<div class="popup">
					<span class="popuptext">The amount of work required for the course.</span>
					<img src="./images/courseDetail/infoIcon_grey.svg">
					</div>
					<div class="courseDetailRatingBoxContainer" style="background-color: %s">
					<span class="courseDetailRatingBox">%s/5</span>
					</div>
					</div>

					</div>

					%s

					<table class="courseDetailTable">
					<tr>
					<th>Class</th>
					<th>Days and Times</th>
					<th>Instructor</th>
					<th>Status</th>
					<th>Friends</th>
					<th>Actions</th>
					</tr>

					<tr>
					<td>CI01</td>
					<td>
					<p>Mon 10:00AM - 11:30AM</p>
					<p>Wed 10:00AM - 11:30AM</p>
					</td>
					<td>%s</td>
					%s
					<td>
					%s
					</td>
					<td class="courseDetailTableActionContainer">
					%s
					</td>
					</tr>

					<tr>
					<td>CI02</td>
					<td>
					<p>Tue 4:00PM - 5:30PM</p>
					<p>Thu 4:00PM - 5:30PM</p>
					</td>
					<td>%s</td>
					%s
					<td>
					%s
					</td>
					<td>
					%s
					</td>
					</tr>

					</table>

					<!-- Student Review Section -->
					<button type="button" class="courseDetailStudentReviewButton collapsibleButton">
					<span>Student Review</span>
					<img src="./images/courseDetail/arrowDropDownIcon.svg">
					</button>
					<div class="courseDetailStudentReviewContainer collapsibleContent hiddenState ">
					<!-- First review -->
					<div class="courseDetailStudentReview">

					<table>
					<tr>
					<th>
					<div class="courseDetailStudentReviewHeader">
					<span>Content</span>
					<div class="courseDetailRatingBoxContainer" style="background-color: %s">
					<span class="courseDetailRatingBox">%s/5</span>
					</div>
					</div>
					</th>
					<th>
					<div class="courseDetailStudentReviewHeader">
					<span>Teaching</span>
					<div class="courseDetailRatingBoxContainer" style="background-color: %s">
					<span class="courseDetailRatingBox">%s/5</span>
					</div>
					</div>
					</th>
					</tr>
					<tr>
					<td>
					<p class="courseDetailStudentReviewDescription">
					Content of the course doesnt make sense at all. It’s pure chaos. no bs lezgoooooo
					</p>
					</td>
					<td>
					<p class="courseDetailStudentReviewDescription">
					Simon best prof 10/10 would recommend -IGN
					</p>
					</td>
					</tr>
					<!-- Second Row -->
					<tr>
					<th>
					<div class="courseDetailStudentReviewHeader">
					<span>Difficulty</span>
					<div class="courseDetailRatingBoxContainer" style="background-color: %s">
					<span class="courseDetailRatingBox">%s/5</span>
					</div>
					</div>
					</th>
					<th>
					<div class="courseDetailStudentReviewHeader">
					<span>Workload</span>
					<div class="courseDetailRatingBoxContainer" style="background-color: %s">
					<span class="courseDetailRatingBox">%s/5</span>
					</div>
					</div>
					</th>
					</tr>
					<tr>
					<td>
					<p class="courseDetailStudentReviewDescription">
					The course is moderately difficult. With enough studying should be not that hard to archive a good grade.
					</p>
					</td>
					<td>
					<p class="courseDetailStudentReviewDescription">
					HAIYAAAAAAAAAAAA so many hw, exam very hard. Overall, I got a A+
					</p>
					</td>
					</tr>
					</table>

					<hr>

					<p>By Bob Lee on Dec 3, 2023</p>

					</div>
					<!-- Second review -->
					<div class="courseDetailStudentReview">

					<table>
					<tr>
					<th>
					<div class="courseDetailStudentReviewHeader">
					<span>Content</span>
					<div class="courseDetailRatingBoxContainer" style="background-color: %s">
					<span class="courseDetailRatingBox">%s/5</span>
					</div>
					</div>
					</th>
					<th>
					<div class="courseDetailStudentReviewHeader">
					<span>Teaching</span>
					<div class="courseDetailRatingBoxContainer" style="background-color: %s">
					<span class="courseDetailRatingBox">%s/5</span>
					</div>
					</div>
					</th>
					</tr>
					<tr>
					<td>
					<p class="courseDetailStudentReviewDescription">
					Some are common sense
					</p>
					</td>
					<td>
					<p class="courseDetailStudentReviewDescription">
					Simon is friendly. He is willing to answer student questions.
					</p>
					</td>
					</tr>
					<!-- Second Row -->
					<tr>
					<th>
					<div class="courseDetailStudentReviewHeader">
					<span>Difficulty</span>
					<div class="courseDetailRatingBoxContainer" style="background-color: %s">
					<span class="courseDetailRatingBox">%s/5</span>
					</div>
					</div>
					</th>
					<th>
					<div class="courseDetailStudentReviewHeader">
					<span>Workload</span>
					<div class="courseDetailRatingBoxContainer" style="background-color: %s">
					<span class="courseDetailRatingBox">%s/5</span>
					</div>
					</div>
					</th>
					</tr>
					<tr>
					<td>
					<p class="courseDetailStudentReviewDescription">
					I didn't study for final but I got an A grade
					</p>
					</td>
					<td>
					<p class="courseDetailStudentReviewDescription">
					No homework. But doing the project will take some time.
					</p>
					</td>
					</tr>
					</table>

					<hr>

					<p>By SuperMan1234 on Apr 1, 2022</p>

					</div>
					</div>  
					</div>

					</div>
				`, course_code, course_name, description,
					rating_colours[0], strconv.Itoa(ratings[0]), 
					rating_colours[1], strconv.Itoa(ratings[1]), 
					rating_colours[2], strconv.Itoa(ratings[2]), 
					rating_colours[3], strconv.Itoa(ratings[3]), 
					fulfillment, 
					instructor, statuses[0], friends[0], buttons[0],
					instructor, statuses[1], friends[1], buttons[1],
					rating_colours[0], strconv.Itoa(ratings[0]), 
					rating_colours[1], strconv.Itoa(ratings[1]), 
					rating_colours[2], strconv.Itoa(ratings[2]), 
					rating_colours[3], strconv.Itoa(ratings[3]), 
					rating_colours[0], strconv.Itoa(ratings[0]), 
					rating_colours[1], strconv.Itoa(ratings[1]), 
					rating_colours[2], strconv.Itoa(ratings[2]), 
					rating_colours[3], strconv.Itoa(ratings[3]))
				fmt.Fprintf(w, `%s`, div_html)
			}
		}
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


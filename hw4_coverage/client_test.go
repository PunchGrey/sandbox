package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

// код писать тут
const (
	xmlFile = "dataset.xml"
)

func getUsers(xmlFile string) []User {
	type XMLRow struct {
		XMLName   xml.Name `xml:"row"`
		ID        int      `xml:"id"`
		GUID      string   `xml:"guid"`
		IsActive  bool     `xml:"isActive"`
		Balance   string   `xml:"balance"`
		Picture   string   `xml:"picture"`
		Age       int      `xml:"age"`
		EyeColor  string   `xml:"eyeColor"`
		FirstName string   `xml:"first_name"`
		LastName  string   `xml:"last_name"`
		Gender    string   `xml:"gender"`
		Company   string   `xml:"company"`
		Email     string   `xml:"email"`
		Phone     string   `xml:"phone"`
		Address   string   `xml:"address"`
		About     string   `xml:"about"`
	}
	type XMLRows struct {
		XMLName xml.Name `xml:"root"`
		Rows    []XMLRow `xml:"row"`
	}

	f, err := os.Open(xmlFile)
	// if we os.Open returns an error then handle it
	if err != nil {
		panic(err)
	}
	defer f.Close()

	byteValue, _ := ioutil.ReadAll(f)
	var xmlRows XMLRows

	err = xml.Unmarshal(byteValue, &xmlRows)
	if err != nil {
		panic(err)
	}

	users := make([]User, len(xmlRows.Rows), len(xmlRows.Rows))

	for i, user := range xmlRows.Rows {
		users[i].Id = user.ID
		users[i].Name = user.FirstName + " " + user.LastName
		users[i].Age = user.Age
		users[i].About = user.About
		users[i].Gender = user.Gender
	}
	return users
}

func selectUsers(users []User, query string) []User {
	if query == "" {
		return users
	}

	outUsers := make([]User, 0, len(users))
	for _, item := range users {
		if strings.Contains(item.Name, query) || strings.Contains(item.About, query) {
			outUsers = append(outUsers, item)
		}
	}
	return outUsers
}

func orderUsers(users []User, orderField string, orderBy int) ([]User, error) {
	if orderBy < -1 || orderBy > 1 {
		return nil, fmt.Errorf("incorrect orderBy")
	}
	if orderBy == 0 {
		return users, nil
	}

	if orderField == "" {
		orderField = "Name"
	}
	switch orderField {
	case "Id":
		if orderBy == 1 {
			sort.Slice(users, func(i, j int) bool {
				return users[i].Id < users[j].Id
			})
		} else {
			sort.Slice(users, func(i, j int) bool {
				return users[i].Id > users[j].Id
			})
		}
		return users, nil
	case "Age":
		if orderBy == 1 {
			sort.Slice(users, func(i, j int) bool {
				return users[i].Age < users[j].Age
			})
		} else {
			sort.Slice(users, func(i, j int) bool {
				return users[i].Age > users[j].Age
			})
		}
		return users, nil
	case "Name":
		if orderBy == 1 {
			sort.Slice(users, func(i, j int) bool {
				return users[i].Name < users[j].Name
			})
		} else {
			sort.Slice(users, func(i, j int) bool {
				return users[i].Name > users[j].Name
			})
		}
		return users, nil
	default:
		return nil, fmt.Errorf("ErrorBadOrderField")

	}
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	accessToken := "7777" //токен для авторизации
	at := r.Header.Get("AccessToken")
	if accessToken != at {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"Error": "Bad AccessToken"}`)
	}

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		io.WriteString(w, `{"Error": "error convert  string  to int for limit"}`)
		return
	}
	offset, err := strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		io.WriteString(w, `{"Error": "error convert  string  to int for offset"}`)
		return
	}
	query := r.FormValue("query")
	if query == "enimForTimeOut" {
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusGatewayTimeout)
		return
	}
	if query == "enumStatusInternalServerError" {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"Error": "internal server error"}`)
		return
	}
	if query == "enumStatusBadRequestCantUnpack" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `[internal server error"]`)
		return
	}
	if query == "enumUnknownBadRequest" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"Error": "unknown bad request"}`)
		return
	}
	if query == "enumCantUnpackJsonArrray" {
		io.WriteString(w, `{"Error": "enumCantUnpackJsonArrray"}`)
		return
	}

	orderField := r.FormValue("order_field")
	orderBy, err := strconv.Atoi(r.FormValue("order_by"))
	if err != nil {
		io.WriteString(w, `{"Error": "error convert  string  to int for orderBy"}`)
		return
	}
	users := getUsers(xmlFile)
	users = selectUsers(users, query)
	users, err = orderUsers(users, orderField, orderBy)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"Error": "`+err.Error()+`"}`)
		return
	}
	if len(users) < offset {
		jsonString, err := json.Marshal([]User{})
		if err != nil {
			io.WriteString(w, `{"Error": "cant pack in json empty array"}`)
		}
		w.Write(jsonString)
		return
	}
	if len(users) > offset+limit {
		jsonString, err := json.Marshal(users[offset : offset+limit])
		if err != nil {
			io.WriteString(w, `{"Error": "cant pack in json"}`)
		}
		w.Write(jsonString)
	} else {
		jsonString, err := json.Marshal(users[offset:len(users)])
		if err != nil {
			io.WriteString(w, `{"Error": "cant pack in json"}`)
		}
		w.Write(jsonString)
	}

}

func TestGetUser(t *testing.T) {
	users := getUsers("dataset.xml")
	user := User{0, "Boyd Wolf", 22, `Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.
`, "male"}
	if !reflect.DeepEqual(users[0], user) {
		t.Errorf("received %v, expected %v", users[0].About, user.About)
	}
}

func TestFindUsers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	sc := &SearchClient{AccessToken: "7777", URL: ts.URL}
	searchReq := SearchRequest{
		Limit:      10,
		Offset:     20,
		Query:      "enim",
		OrderField: "Id",
		OrderBy:    -1,
	}
	searchRes, err := sc.FindUsers(searchReq)
	if err != nil {
		t.Errorf("unexpected error: %#v", err)
		return
	}
	if len(searchRes.Users) != 0 {
		t.Errorf("unexpected value: %#v", searchRes.Users)
	}

}

func TestFindUsersLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	sc := &SearchClient{AccessToken: "7777", URL: ts.URL}
	searchReq := SearchRequest{
		Limit:      -1,
		Offset:     20,
		Query:      "enim",
		OrderField: "Id",
		OrderBy:    -1,
	}
	_, err := sc.FindUsers(searchReq)
	if err == nil {
		t.Errorf("unexpected error %v, expected error %s", err, "limit must be > 0")
	}
	if err.Error() != "limit must be > 0" {
		t.Errorf("unexpected error %s, expected error %s", err.Error(), "limit must be > 0")
	}

	searchReq = SearchRequest{
		Limit:      30,
		Offset:     0,
		Query:      "",
		OrderField: "Id",
		OrderBy:    0,
	}
	searchRes, err := sc.FindUsers(searchReq)
	if len(searchRes.Users) != 25 {
		t.Errorf("unexpected value  len(searchRes.Users) %d", len(searchRes.Users))
	}
}

func TestFindUsersOffset(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	sc := &SearchClient{AccessToken: "7777", URL: ts.URL}

	searchReq := SearchRequest{
		Limit:      10,
		Offset:     -1,
		Query:      "",
		OrderField: "Id",
		OrderBy:    0,
	}
	_, err := sc.FindUsers(searchReq)
	if err == nil {
		t.Errorf("unexpected error %v, expected error %s", err, "offset must be > 0")
	}
	if err.Error() != "offset must be > 0" {
		t.Errorf("unexpected error %s, expected error %s", err.Error(), "offset must be > 0")
	}
}

func TestFindUsersTimeOut(t *testing.T) {
	//handlerFunc := http.HandlerFunc(SearchServer)
	//ts := httptest.NewServer(http.TimeoutHandler(handlerFunc, 1*time.Millisecond, "server timeout"))

	//ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	//sc := &SearchClient{AccessToken: "7777", URL: ts.URL}
	sc := &SearchClient{AccessToken: "7777", URL: "http://127.0.0.1:3560"}

	searchReq := SearchRequest{
		Limit:      25,
		Offset:     0,
		Query:      "enimForTimeOut",
		OrderField: "Name",
		OrderBy:    1,
	}
	_, err := sc.FindUsers(searchReq)
	if !strings.HasPrefix(err.Error(), "unknown error") {
		t.Errorf("unexpected error %s, expected error %s", err.Error(), "unknown error ...")
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	sc = &SearchClient{AccessToken: "7777", URL: ts.URL}
	_, err = sc.FindUsers(searchReq)

	if !strings.HasPrefix(err.Error(), "timeout for") {
		t.Errorf("unexpected error %s, expected error %s", err.Error(), "timeout for...")
	}
}

func TestFindUsersAuth(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	sc := &SearchClient{AccessToken: "77776666", URL: ts.URL}
	searchReq := SearchRequest{
		Limit:      25,
		Offset:     0,
		Query:      "enim",
		OrderField: "Id",
		OrderBy:    -1,
	}
	_, err := sc.FindUsers(searchReq)
	if err.Error() != "Bad AccessToken" {
		t.Errorf("unexpected error %s, expected error %s", err.Error(), "Bad AccessToken")
	}
}

func TestFindUserInternalServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	sc := &SearchClient{AccessToken: "7777", URL: ts.URL}
	searchReq := SearchRequest{
		Limit:      25,
		Offset:     0,
		Query:      "enumStatusInternalServerError",
		OrderField: "Id",
		OrderBy:    -1,
	}
	_, err := sc.FindUsers(searchReq)
	if err.Error() != "SearchServer fatal error" {
		t.Errorf("unexpected error %s, expected error %s", err.Error(), "SearchServer fatal error")
	}
}

func TestBadRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	sc := &SearchClient{AccessToken: "7777", URL: ts.URL}
	searchReq := SearchRequest{
		Limit:      25,
		Offset:     0,
		Query:      "enumStatusBadRequestCantUnpack",
		OrderField: "Id",
		OrderBy:    -1,
	}
	_, err := sc.FindUsers(searchReq)
	if !strings.HasPrefix(err.Error(), "cant unpack error json") {
		t.Errorf("unexpected error %s, expected error %s", err.Error(), "cant unpack error json...")
	}

	searchReq = SearchRequest{
		Limit:      25,
		Offset:     0,
		Query:      "enum",
		OrderField: "Gender",
		OrderBy:    -1,
	}
	_, err = sc.FindUsers(searchReq)
	if !strings.HasPrefix(err.Error(), "OrderFeld") {
		t.Errorf("unexpected error %s, expected error %s", err.Error(), "OrderFeld...")
	}

	searchReq = SearchRequest{
		Limit:      25,
		Offset:     0,
		Query:      "enumUnknownBadRequest",
		OrderField: "Id",
		OrderBy:    -1,
	}
	_, err = sc.FindUsers(searchReq)
	if !strings.HasPrefix(err.Error(), "unknown bad request error") {
		t.Errorf("unexpected error %s, expected error %s", err.Error(), "unknown bad request error...")
	}
}

func TestCantUnpackJsonArrray(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	sc := &SearchClient{AccessToken: "7777", URL: ts.URL}
	searchReq := SearchRequest{
		Limit:      25,
		Offset:     0,
		Query:      "enumCantUnpackJsonArrray",
		OrderField: "Id",
		OrderBy:    -1,
	}
	_, err := sc.FindUsers(searchReq)
	if !strings.HasPrefix(err.Error(), "cant unpack result json:") {
		t.Errorf("unexpected error %s, expected error %s", err.Error(), "cant unpack result json:...")
	}
}

func TestLenDataEqualLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	sc := &SearchClient{AccessToken: "7777", URL: ts.URL}
	searchReq := SearchRequest{
		Limit:      17,
		Offset:     0,
		Query:      "enim",
		OrderField: "Id",
		OrderBy:    -1,
	}
	result, err := sc.FindUsers(searchReq)
	if err != nil {
		t.Errorf("unexpected error %s, expected error nil", err.Error())
	}
	if len(result.Users) != searchReq.Limit {
		t.Errorf("unexpected voluem %d, expected error %d", len(result.Users), searchReq.Limit)
	}

	searchReq = SearchRequest{
		Limit:      20,
		Offset:     0,
		Query:      "enim",
		OrderField: "Id",
		OrderBy:    -1,
	}
	result, err = sc.FindUsers(searchReq)
	if err != nil {
		t.Errorf("unexpected error %s, expected error nil", err.Error())
	}
	if len(result.Users) == searchReq.Limit {
		t.Errorf("unexpected voluem %d equel limit", len(result.Users))
	}
}

/*
    client := http.Client{
        Timeout: time.Duration(5 * time.Millisecond),
    }

    _, err := client.Get("https://golangcode.com")
    if err != nil {
        log.Fatal(err)
	}
*/

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
		if strings.Contains(item.Name, query) && strings.Contains(item.About, query) {
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
				return users[i].Id > users[j].Id
			})
		} else {
			sort.Slice(users, func(i, j int) bool {
				return users[i].Id < users[j].Id
			})
		}
		return users, nil
	case "Age":
		if orderBy == 1 {
			sort.Slice(users, func(i, j int) bool {
				return users[i].Age > users[j].Age
			})
		} else {
			sort.Slice(users, func(i, j int) bool {
				return users[i].Age < users[j].Age
			})
		}
		return users, nil
	case "Name":
		if orderBy == 1 {
			sort.Slice(users, func(i, j int) bool {
				return users[i].Name > users[j].Name
			})
		} else {
			sort.Slice(users, func(i, j int) bool {
				return users[i].Name < users[j].Name
			})
		}
		return users, nil
	default:
		return nil, fmt.Errorf("incorrect orderField")

	}
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	accessToken := "7777" //токен для авторизации
	at := r.Header.Get("AccessToken")
	if accessToken != at {
		io.WriteString(w, `{"Error": "Bad AccessToken"}`)
	}

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		io.WriteString(w, `{"Error": "error convert  string  to int for limit"}`)
		return
	}
	//	offset, err := strconv.Atoi(r.FormValue("offset"))
	//	if err != nil {
	//		io.WriteString(w, `{"Error": "error convert  string  to int for offset"}`)
	//		return
	//	}
	query := r.FormValue("query")
	orderField := r.FormValue("orderField")
	orderBy, err := strconv.Atoi(r.FormValue("orderBy"))
	if err != nil {
		io.WriteString(w, `{"Error": "error convert  string  to int for orderBy"}`)
		return
	}
	users := getUsers(xmlFile)
	users = selectUsers(users, query)
	users, err = orderUsers(users, orderField, orderBy)
	if len(users) > limit {
		jsonString, err := json.Marshal(users[0:limit])
		if err != nil {
			io.WriteString(w, `{"Error": "cant pack in json"}`)
		}
		w.Write(jsonString)
	} else {
		jsonString, err := json.Marshal(users)
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

func TestFindUasers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	sc := &SearchClient{AccessToken: "7777", URL: ts.URL}
	searchReq := SearchRequest{
		Limit:      10,
		Offset:     0,
		Query:      "",
		OrderField: "Name",
		OrderBy:    0,
	}
	searchRes, err := sc.FindUsers(searchReq)
	if err != nil {
		t.Errorf("unexpected error: %#v", err)
		return
	}
	if len(searchRes.Users) < 1 {
		t.Errorf("unexpected value: %#v", searchRes.Users)
	}

}

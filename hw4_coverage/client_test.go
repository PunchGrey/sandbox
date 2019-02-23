package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

// код писать тут

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
		users[i].Name = user.FirstName + user.LastName
		users[i].Age = user.Age
		users[i].About = user.About
		users[i].Gender = user.Gender
	}
	return users
}

func SearchServer(sc *SearchClient, sr *SearchRequest) (*SearchResponse, error) {
	return &SearchResponse{}, nil
}

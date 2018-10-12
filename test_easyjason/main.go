package main

import (
	"fmt"
)

/*
//easyjson:json
type JSONData struct {
	Browsers []string
	Company  string
}
*/
//easyjson -lower_camel_case main.go
//easyjson:json
type JSONData struct {
	Browsers []string
	Company  string
	Country  string
	Email    string
	Job      string
	Name     string
	Phone    string
}

//easyjson:json
//type JSONDataList []JSONData

func main() {

	//	var d JSONData
	dl := JSONData{} //JSONDataList{}
	data := `{
			"browsers":
			[
				"Mozilla/5.0 (Linux; Android 4.4.2; LG-V410 Build/KOT49I.V41010d) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.103 Safari/537.36",
				"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.6; rv:25.0) Gecko/20100101 Firefox/25.0",
				"msnbot-media/1.1 ( http://search.msn.com/msnbot.htm)",
				"Mozilla/4.0 (compatible; MSIE 6.0; Windows CE; IEMobile 8.12; MSIEMobile6.0)"
			],
			"company": "Eazzy",
			"country": "United States Virgin Islands",
			"Email": "DebraReynolds@Kwilith.info",
			"job": "Cost Accountant"	,
			"name": "Martin Wilson",
			"phone": "7-159-571-85-58"
		}`
	err := dl.UnmarshalJSON([]byte(data))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dl)

}

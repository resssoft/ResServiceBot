package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	token  string
	urlTmp = "https://api.apilayer.com/fixer/convert?to=%s&from=%s&amount=%s"
)

type converterResult struct {
	Result  float64 `json:"result,omitempty"`
	Success bool    `json:"success,omitempty"`
}

//rate cache (to 6h)

func configureConverter(val string) {
	token = val
}

func fiat(from, to, amount string) string {
	time.Sleep(time.Millisecond * 100)
	url := fmt.Sprintf(urlTmp, to, from, amount)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", token)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(string(body))
	var resultObj converterResult
	err = json.Unmarshal(body, &resultObj)
	if err != nil {
		fmt.Println(err)
	} else {
		return fmt.Sprintf("%.2f", resultObj.Result)
	}
	return ""
}

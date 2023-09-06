package main

import (
	"net/http"
	"io/ioutil"
)

func getLink(url string, c chan []byte) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	c <- html
}

func RetrieveHTML(metofficeCode string, bbcWeatherCode string) ([]byte, []byte) {
	var metofficeURL string = "https://www.metoffice.gov.uk/weather/forecast/" + metofficeCode
	var bbcWeatherURL string = "https://www.bbc.co.uk/weather/" + bbcWeatherCode
	metofficeChannel := make(chan []byte)
	bbcWeatherChannel := make(chan []byte)
	go getLink(metofficeURL, metofficeChannel)
	go getLink(bbcWeatherURL, bbcWeatherChannel)
	
	metofficeHTML, bbcWeatherHTML := <- metofficeChannel, <- bbcWeatherChannel
	
	return metofficeHTML, bbcWeatherHTML
}

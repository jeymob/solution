package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Currency struct {
	Valute []InfoOfCurrency
}

type InfoOfCurrency struct {
	CharCode string  `xml:"CharCode"`
	Nominal  uint16  `xml:"Nominal"`
	Name     string  `xml:"Name"`
	Value    float64 `xml:"Value"`
}

type CurrencysMap map[string]Currency

func (val *CurrencysMap) GetOneDayCurrency(date string) error {
	resp, err := http.Get(QueryAPI + date)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	xmlList, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	xmlList = bytes.ReplaceAll(xmlList, []byte(","), []byte("."))
	v := Currency{}
	err = xml.Unmarshal(xmlList[OffsetIgnoreEncoding:], &v)
	if err != nil {
		return err
	}
	(*val)[date] = v
	return nil
}

func (val *CurrencysMap) GetAllDaysCurrency() error {

	t := time.Now()
	date := t.Format("02/01/2006")
	for i := 0; i <= int(CountOfDays); i++ {
		err := val.GetOneDayCurrency(date)
		if err != nil {
			return err
		}
		t = time.Unix(t.Unix()-CountSecondOfDay, 0)
		date = t.Format("02/01/2006")
	}
	return nil
}

func (val *CurrencysMap) MinMaxAverage() {

	var count uint32 = 0
	var infoMax, infoMin InfoOfCurrency
	var dayMin, dayMax string
	var min, max, average float64
	for day, arrVal := range *val {
		for _, infoVal := range arrVal.Valute {
			average += infoVal.Value / float64(infoVal.Nominal)
			count++
			if min == 0 || min > infoVal.Value/float64(infoVal.Nominal) {
				min = infoVal.Value / float64(infoVal.Nominal)
				infoMin = infoVal
				dayMin = day
			}
			if max < infoVal.Value/float64(infoVal.Nominal) {
				max = infoVal.Value / float64(infoVal.Nominal)
				infoMax = infoVal
				dayMax = day
			}
		}
	}
	fmt.Println(infoMax.Value, infoMax.CharCode, dayMax)
	fmt.Println(infoMin.Value, infoMin.CharCode, dayMin)
	fmt.Printf("%.4f", average/float64(count))
}

func main() {
	var valutes = make(CurrencysMap)
	err := valutes.GetAllDaysCurrency()
	if err != nil {
		log.Fatal(err)
	}
	valutes.MinMaxAverage()
}

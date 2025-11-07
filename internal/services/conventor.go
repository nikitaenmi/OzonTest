package services

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CBRClient interface {
	GetRatesXML(date time.Time) (string, error)
}

type CurrencyConverter struct {
	cbrClient CBRClient
}

func NewCurrencyConverter(cbrClient CBRClient) CurrencyConverter {
	return CurrencyConverter{
		cbrClient: cbrClient,
	}
}

func (c *CurrencyConverter) ToRUB(amount float64, fromCurrency string, date time.Time) (float64, error) {
	if fromCurrency == "" || fromCurrency == "RUB" {
		return amount, nil
	}

	xmlData, err := c.cbrClient.GetRatesXML(date)
	if err != nil {
		return 0, fmt.Errorf("get rates XML: %v", err)
	}

	rates, err := c.parseRatesXML(xmlData)
	if err != nil {
		return 0, fmt.Errorf("parse rates: %v", err)
	}

	rate, exists := rates[fromCurrency]
	if !exists {
		return 0, fmt.Errorf("currency %s not found", fromCurrency)
	}

	return amount * rate, nil
}

func (c *CurrencyConverter) parseRatesXML(xmlData string) (map[string]float64, error) {
	type Valute struct {
		CharCode  string `xml:"CharCode"`
		VunitRate string `xml:"VunitRate"`
	}

	type ValCurs struct {
		Valutes []Valute `xml:"Valute"`
	}

	var valCurs ValCurs
	if err := xml.Unmarshal([]byte(xmlData), &valCurs); err != nil {
		return nil, fmt.Errorf("parse XML: %v", err)
	}

	rates := make(map[string]float64)
	for _, valute := range valCurs.Valutes {
		valueStr := strings.Replace(valute.VunitRate, ",", ".", -1)
		rate, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return nil, fmt.Errorf("parse rate for %s: %v", valute.CharCode, err)
		}
		rates[valute.CharCode] = rate
	}

	return rates, nil
}

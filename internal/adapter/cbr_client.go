package adapter

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

type MockCBRClient struct {
	baseRates map[string]float64
}

func NewMockCBRClient() *MockCBRClient {
	return &MockCBRClient{
		baseRates: map[string]float64{
			"USD": 75.0,
			"EUR": 90.0,
			"BYN": 25.0,
			"KZT": 0.15,
			"JPY": 0.65,
		},
	}
}

func (m *MockCBRClient) GetRatesXML(date time.Time) (string, error) {
	return m.generateCBRXML(date)
}

// generateCBRXML генерирует XML с курсами валют
//
// Формула вариации курса создает случайное, но детерминированное
// изменение, принимая дату.
//
//	variation = (day * 10 + month * 5 + year) / 100
func (m *MockCBRClient) generateCBRXML(date time.Time) (string, error) {
	variation := float64(date.Day()*10+int(date.Month())*5+date.Year()) / 100.0

	type Valute struct {
		ID        string `xml:"ID,attr"`
		NumCode   string `xml:"NumCode"`
		CharCode  string `xml:"CharCode"`
		Nominal   int    `xml:"Nominal"`
		Name      string `xml:"Name"`
		Value     string `xml:"Value"`
		VunitRate string `xml:"VunitRate"`
	}

	type ValCurs struct {
		XMLName xml.Name `xml:"ValCurs"`
		Date    string   `xml:"Date,attr"`
		Name    string   `xml:"name,attr"`
		Valutes []Valute `xml:"Valute"`
	}

	valCurs := ValCurs{
		Date: date.Format("02/01/2006"),
		Name: "Foreign Currency Market",
	}

	currencies := []string{"USD", "EUR", "BYN", "KZT", "JPY"}
	for i, code := range currencies {
		baseRate := m.baseRates[code]
		rate := baseRate + variation
		valueStr := strings.Replace(fmt.Sprintf("%.4f", rate), ".", ",", -1)

		valCurs.Valutes = append(valCurs.Valutes, Valute{
			ID:        fmt.Sprintf("R%05d", i),
			NumCode:   fmt.Sprintf("%03d", i),
			CharCode:  code,
			Nominal:   1,
			Name:      m.currencyName(code),
			Value:     valueStr,
			VunitRate: valueStr,
		})
	}

	output, err := xml.MarshalIndent(valCurs, "", "  ")
	if err != nil {
		return "", err
	}

	return xml.Header + string(output), nil
}

func (m *MockCBRClient) currencyName(code string) string {
	names := map[string]string{
		"USD": "Доллар США",
		"EUR": "Евро",
		"BYN": "Белорусский рубль",
		"KZT": "Тенге",
		"JPY": "Иена",
	}
	return names[code]
}

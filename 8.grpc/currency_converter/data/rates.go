package data

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"k8s.io/klog/v2"
)

type ExchangeRates struct {
	log   klog.Logger
	rates map[string]float64
}

func NewExchangeRates(logger *klog.Logger) *ExchangeRates {
	er := &ExchangeRates{
		log:   *logger,
		rates: make(map[string]float64),
	}
	er.getRates()
	return er
}

func (e *ExchangeRates) GetRate(base, destination string) (float64, error) {
	br, ok := e.rates[base]
	if !ok {
		return 0, fmt.Errorf("base currency %s not found", base)
	}
	dr, ok := e.rates[destination]
	if !ok {
		return 0, fmt.Errorf("destination currency %s not found", destination)
	}
	return dr / br, nil
}

func (e *ExchangeRates) getRates() error {
	resp, err := http.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		e.log.Error(err, "Cannot get exchange rates from bank")
		return nil
	}
	defer resp.Body.Close()

	// contents, _ := io.ReadAll(resp.Body)
	// err = xml.Unmarshal([]byte(string(contents)), &cube)

	var cube Cubes
	err = xml.NewDecoder(resp.Body).Decode(&cube)
	if err != nil {
		e.log.Error(err, "Cannot decode the xml data")
		return nil
	}
	for _, rate := range cube.CubeData {
		e.rates[rate.Currency], _ = strconv.ParseFloat(rate.Rate, 64)
	}
	e.rates["EUR"] = 1.0
	return nil
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}

// <Envelope>
//   <Cube>
//     <Cube time="2025-09-02">
//       <Cube currency="USD" rate="1.1646"/>
//       <Cube currency="JPY" rate="173.10"/>
//       <Cube currency="BGN" rate="1.9558"/>
//     </Cube>
//   </Cube>
// </Envelope>

// Start at <Cube> (the first one under <Envelope>),

// then go into its child <Cube time="...">,

// then go into its child <Cube currency="..." rate="..."/>.

// type Envelope struct {
// 	XMLName xml.Name `xml:"Envelope"`
// 	Cube    CubeRoot `xml:"Cube"`
// }

// type CubeRoot struct {
// 	Days []Day `xml:"Cube"`
// }

// type Day struct {
// 	Time  string `xml:"time,attr"`
// 	Rates []Rate `xml:"Cube"`
// }

// type Rate struct {
// 	Currency string  `xml:"currency,attr"`
// 	Rate     float64 `xml:"rate,attr"`
// }

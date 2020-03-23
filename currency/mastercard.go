package currency

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const url = "https://www.mastercard.us/settlement/currencyrate/fxDate=%s;transCurr=%s;crdhldBillCurr=USD;bankFee=0;transAmt=1/conversion-rate"

func FetchExchangeToUSD(isocode string) (*CurrencyExchange, error) {
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	nowBrazil := time.Now().In(loc)

	url := fmt.Sprintf(url, now.Format("2006-01-02"), isocode)

	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("accept-language", "en-US,en;q=0.9,pt-BR;q=0.8,pt;q=0.7")
	req.Header.Add("user-agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Mobile Safari/537.36")
	req.Header.Add("referer", "https://www.mastercard.us/en-us/consumers/get-support/convert-currency.html")
	req.Header.Add("authority", "www.mastercard.us")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	response, ok := response["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Failed to read 'data'")
	}

	rate, ok := response["conversionRate"].(float64)
	if !ok {
		return nil, errors.New("Failed to read 'conversionRate'")
	}

	return &CurrencyExchange{rate, now}, nil
}

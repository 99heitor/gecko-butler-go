package currency

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const url = "https://ptax.bcb.gov.br/ptax_internet/consultaBoletim.do?method=gerarCSVFechamentoMoedaNoPeriodo&ChkMoeda=61&DATAINI=%s&DATAFIM=%s"

func FetchYesterdayBRLtoUSD() (*CurrencyExchange, error) {
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	nowBrazil := time.Now().In(loc)

	startDate := nowBrazil.AddDate(0, 0, -4)

	url := fmt.Sprintf(url, startDate.Format("02/01/2006"), nowBrazil.Format("02/01/2006"))

	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	r := csv.NewReader(res.Body)
	r.Comma = ';'

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	latest_row := records[len(records)-1]
	rate, err := strconv.ParseFloat(strings.ReplaceAll(latest_row[5], ",", "."), 64)
	if err != nil {
		return nil, err
	}

	return &CurrencyExchange{rate, nowBrazil}, nil
}

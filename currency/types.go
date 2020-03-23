package currency

import "time"

type CurrencyExchange struct {
	Rate float64
	Time time.Time
}

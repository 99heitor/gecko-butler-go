package commands

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/99heitor/gecko-butler-go/currency"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"time"
	"fmt"
)

const projectID = "geckobutler"

const (
	iof    = 6.38 / 100
	spread = 4.0 / 100
)

func convertToBRL(value float64, exchangeRate float64) {
	return value * exchangeRate * (1 + spread) * (1 + iof)
}

func getCurrencyKey(isocode string) datastore.NameKey {
	return datastore.NameKey("currencyExchange", isocode, nil)
}

func isDateExpired(date time.Time) (bool) {
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	nowBrazil := time.Now().In(loc)
	dateBrazil := date.In(loc)

	return (
		dateBrazil.Year() < nowBrazil.Year() ||
		dateBrazil.YearDay() < nowBrazil.YearDay()
	)
}

func getCurrencyExchange(ctx context.Context, isocode string) (float64, error) {

	var exchangeToUSD currency.CurrencyExchange
	if isocode != "USD" {
		key := getCurrencyKey(isocode)
		err := datastoreClient.Get(ctx, key, &exchangeToUSD)

		// Fetch new exchange
		if err != nil || isDateExpired(exchangeToUSD.Time) {
			exchangeToUSD, err = currency.FetchExchangeToUSD(isocode)
			if err != nil {
				return 0, fmt.Errorf("Failed to fetch exchangeToUSD: %w", err)
			}

			_, err = datastoreClient.Put(ctx, key, &exchangeToUSD)
			if err != nil {
				return 0, fmt.Errorf("Failed to cache exchangeToUSD: %w", err)
			}
		}
	} else {
		exchangeToUSD = CurrencyExchange{1.0, time.Now()}
	}

	var exchangeToBRL currency.CurrencyExchange
	key := getCurrencyKey("USD2BRL")
	err := datastoreClient.Get(ctx, key, &exchangeToBRL)

	// Fetch new exchange
	if err != nil || isDateExpired(exchangeToUSD.Time) {
		exchangeToBRL, err = currency.FetchYesterdayBRLtoUSD()
		if err != nil {
			return 0, fmt.Errorf("Failed to fetch exchangeToBRL: %w", err)
		}

		_, err = datastoreClient.Put(ctx, key, &exchangeToBRL)
		if err != nil {
			return 0, fmt.Errorf("Failed to cache exchangeToBRL: %w", err)
		}
	}

	return exchangeToUSD.Rate * exchangeToBRL.Rate, nil
}

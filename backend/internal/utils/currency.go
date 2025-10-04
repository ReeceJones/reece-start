package utils

import "github.com/stripe/stripe-go/v83"

func GetCurrencyForCountry(country string) stripe.Currency {
	switch country {
	case "US":
		return stripe.CurrencyUSD
	case "CA":
		return stripe.CurrencyCAD
	case "GB":
		return stripe.CurrencyGBP
	case "AU":
		return stripe.CurrencyAUD
	case "NZ":
		return stripe.CurrencyNZD
	case "JP":
		return stripe.CurrencyJPY
	case "KR":
		return stripe.CurrencyKRW
	case "CN":
		return stripe.CurrencyCNY
	case "IN":
		return stripe.CurrencyINR
	case "BR":
		return stripe.CurrencyBRL
	case "MX":
		return stripe.CurrencyMXN
	case "AR":
		return stripe.CurrencyARS
	case "CO":
		return stripe.CurrencyCOP
	case "CL":
		return stripe.CurrencyCLP
	case "PE":
		return stripe.CurrencyPEN
	case "CH":
		return stripe.CurrencyCHF
	case "SE":
		return stripe.CurrencySEK
	case "NO":
		return stripe.CurrencyNOK
	case "DK":
		return stripe.CurrencyDKK
	}

	return stripe.CurrencyEUR
}
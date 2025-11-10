package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stripe/stripe-go/v83"
)

func TestGetCurrencyForCountry(t *testing.T) {
	t.Run("ValidCountry", func(t *testing.T) {
		currency := GetCurrencyForCountry("US")
		assert.Equal(t, stripe.CurrencyUSD, currency)
	})

	t.Run("InvalidCountry", func(t *testing.T) {
		currency := GetCurrencyForCountry("INVALID")
		assert.Equal(t, stripe.CurrencyEUR, currency)
	})
}

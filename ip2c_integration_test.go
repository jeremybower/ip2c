// +build integration

package ip2c

import (
	"testing"
)

func TestLookupIPv4UsingIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	client := NewClient()
	countryInfo, err := client.LookupIPv4("1.1.1.1")
	if err != nil {
		t.Fatal(err)
	}

	if countryInfo == nil {
		t.Error("Nil response from client")
	}
}

func TestLookupDecimalUsingIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	// "1.1.1.1" = 16843009
	client := NewClient()
	countryInfo, err := client.LookupDecimal(16843009)
	if err != nil {
		t.Fatal(err)
	}

	if countryInfo == nil {
		t.Error("Nil response from client")
	}
}

func TestLookupSelfUsingIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	client := NewClient()
	countryInfo, err := client.LookupSelf()
	if err != nil {
		t.Fatal(err)
	}

	if countryInfo == nil {
		t.Error("Nil response from client")
	}
}

func TestErrWrongInputUsingIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	client := NewClient()
	countryInfo, err := client.LookupIPv4("a.b.c.d")
	if err != ErrWrongInput {
		t.Errorf("Expected ErrWrongInput, but found: %s", err)
	}

	if countryInfo != nil {
		t.Error("Expected nil country info for wrong input")
	}
}

package ip2c

import "testing"

func TestSimpleClientForTesting(t *testing.T) {
	client := NewSimpleClientForTesting()

	countryInfo, err := client.LookupDecimal(0)
	if countryInfo == nil || err != nil {
		t.Error("unexpected return values")
	}

	countryInfo, err = client.LookupIPv4("0.0.0.0")
	if countryInfo == nil || err != nil {
		t.Error("unexpected return values")
	}

	countryInfo, err = client.LookupSelf()
	if countryInfo == nil || err != nil {
		t.Error("unexpected return values")
	}
}

func TestErrorClientForTesting(t *testing.T) {
	client := NewErrorClientForTesting()

	countryInfo, err := client.LookupDecimal(0)
	if countryInfo != nil || err != ErrForTesting {
		t.Error("unexpected return values")
	}

	countryInfo, err = client.LookupIPv4("0.0.0.0")
	if countryInfo != nil || err != ErrForTesting {
		t.Error("unexpected return values")
	}

	countryInfo, err = client.LookupSelf()
	if countryInfo != nil || err != ErrForTesting {
		t.Error("unexpected return values")
	}
}

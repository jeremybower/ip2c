package ip2c

import "errors"

type clientForTesting struct {
	countryInfo *CountryInfo
	err         error
}

func (c *clientForTesting) LookupIPv4(ip string) (*CountryInfo, error) {
	return c.countryInfo, c.err
}

func (c *clientForTesting) LookupDecimal(dec int) (*CountryInfo, error) {
	return c.countryInfo, c.err
}

func (c *clientForTesting) LookupSelf() (*CountryInfo, error) {
	return c.countryInfo, c.err
}

// NewSimpleClientForTesting is useful for unit tests where the CountryInfo
// needs to be valid, but not an exact lookup of the IP address. It's a
// reasonable approach because unit tests typically run on localhost or
// in a container.
func NewSimpleClientForTesting() Client {
	return NewSimpleClientForTestingEx(&CountryInfo{
		TwoLetterCode:   "CA",
		ThreeLetterCode: "CAN",
		FullName:        "Canada",
	})
}

// NewSimpleClientForTestingEx provides the same functionality as
// NewSimpleClientForTesting, but allows the CountryInfo to be customized.
func NewSimpleClientForTestingEx(countryInfo *CountryInfo) Client {
	return &clientForTesting{countryInfo: countryInfo}
}

// ErrForTesting is returned by NewErrorClientForTesting to simulate errors
// returned by the Client for unit tests.
var ErrForTesting = errors.New("error from IP2C for testing")

// NewErrorClientForTesting returns ErrForTesting to simulate errors returned
// by the Client for unit tests.
func NewErrorClientForTesting() Client {
	return NewErrorClientForTestingEx(ErrForTesting)
}

// NewErrorClientForTestingEx proveds the same functionality as
// NewErrorClientForTesting, but allows the error to be customized.
func NewErrorClientForTestingEx(err error) Client {
	return &clientForTesting{err: err}
}

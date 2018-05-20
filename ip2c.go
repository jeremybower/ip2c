package ip2c

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	// ErrWrongInput represents an error returned by IP2C in the body of the
	// response when the given IP has an invalid syntax.
	ErrWrongInput = errors.New("ip2c: your request has not been processed due to invalid syntax")

	// ErrUnknown represents an error returned by IP2C in the body of the
	// response when the given IP is properly formatted but unknown.
	ErrUnknown = errors.New("ip2c: given ip/dec not found in database or not yet physically assigned to any country")
)

// CountryInfo is a the country information returned by IP2C
type CountryInfo struct {
	TwoLetterCode   string
	ThreeLetterCode string
	FullName        string
}

// Client is the interface to IP2C
type Client interface {
	LookupIPv4(ip string) (*CountryInfo, error)
	LookupDecimal(dec int) (*CountryInfo, error)
	LookupSelf() (*CountryInfo, error)
}

// Options are the configuration options for the client.
type Options struct {
	BaseURL    string
	HTTPClient *http.Client
	ReaderFunc func(io.Reader) io.Reader
}

// NewOptions will create new options with default values.
func NewOptions() *Options {
	return &Options{
		BaseURL:    "https://ip2c.org",
		HTTPClient: &http.Client{},
		ReaderFunc: func(r io.Reader) io.Reader {
			return r
		},
	}
}

// NewClient will create a new IP2C client with default options.
func NewClient() Client {
	opts := NewOptions()
	return NewClientWithOptions(opts)
}

// NewClientWithOptions will create a new IP2C client with the given options.
func NewClientWithOptions(opts *Options) Client {
	return &clientImpl{
		opts: opts,
	}
}

type clientImpl struct {
	opts *Options
}

func (c *clientImpl) LookupIPv4(ip string) (*CountryInfo, error) {
	return c.lookup(fmt.Sprintf(c.opts.BaseURL+"/?ip=%s", ip))
}

func (c *clientImpl) LookupDecimal(dec int) (*CountryInfo, error) {
	return c.lookup(fmt.Sprintf(c.opts.BaseURL+"/?dec=%d", dec))
}

func (c *clientImpl) LookupSelf() (*CountryInfo, error) {
	return c.lookup(c.opts.BaseURL + "/self")
}

func (c *clientImpl) lookup(url string) (*CountryInfo, error) {
	resp, err := c.opts.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Unexpected response. Expected 200 but found %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(c.opts.ReaderFunc(resp.Body))
	if err != nil {
		return nil, err
	}

	countryInfo, err := newCountryInfoFromString(string(body))
	if err != nil {
		return nil, err
	}

	return countryInfo, nil
}

func newCountryInfoFromString(s string) (*CountryInfo, error) {
	// ex. "1;CA;CAN;Canada"
	segments := strings.Split(s, ";")
	if len(segments) != 4 {
		return nil, fmt.Errorf("Invalid format. Expected 4 segments but found %d: \"%s\"", len(segments), s)
	}

	code := segments[0]
	switch code {
	case "0":
		return nil, ErrWrongInput
	case "1":
		break
	case "2":
		return nil, ErrUnknown
	default:
		return nil, fmt.Errorf("Invalid format. Expected code of 0, 1, or 2 in 1st segment but found %s: \"%s\"", code, s)
	}

	twoLetterCode := segments[1]
	if len(twoLetterCode) != 2 {
		return nil, fmt.Errorf("Invalid format. Expected 2 letter code in 2nd segment but found %d: \"%s\"", len(twoLetterCode), s)
	}

	threeLetterCode := segments[2]
	if len(threeLetterCode) != 3 {
		return nil, fmt.Errorf("Invalid format. Expected 3 letter code in 3rd segment but found %d: \"%s\"", len(threeLetterCode), s)
	}

	fullName := strings.TrimSpace(segments[3])
	if len(fullName) == 0 {
		return nil, fmt.Errorf("Invalid format. Expected full name in 4th segment but found blank: \"%s\"", s)
	}

	return &CountryInfo{
		TwoLetterCode:   twoLetterCode,
		ThreeLetterCode: threeLetterCode,
		FullName:        fullName,
	}, nil
}

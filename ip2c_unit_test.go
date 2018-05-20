// +build unit

package ip2c

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupMockServer() (*http.ServeMux, *httptest.Server, func()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	return mux, server, func() {
		server.Close()
	}
}

func TestTimeoutErrorUsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	ch := make(chan int)
	defer func() { ch <- 0 }()

	mux.HandleFunc("/self", func(w http.ResponseWriter, r *http.Request) {
		_ = <-ch
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "1;CA;CAN;Canada")
	})

	opts := NewOptions()
	opts.BaseURL = server.URL
	opts.HTTPClient = &http.Client{
		Timeout: 10 * time.Millisecond,
	}
	client := NewClientWithOptions(opts)

	countryInfo, err := client.LookupSelf()

	if err == nil {
		t.Error("Expected timeout error")
	}

	if countryInfo != nil {
		t.Error("Expected nil country info")
	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestReadErrorUsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	mux.HandleFunc("/self", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "1;CA;CAN;Canada")
	})

	opts := NewOptions()
	opts.BaseURL = server.URL
	opts.ReaderFunc = func(io.Reader) io.Reader {
		return errReader(0)
	}
	client := NewClientWithOptions(opts)

	countryInfo, err := client.LookupSelf()
	expectedError := "test error"
	if err.Error() != expectedError {
		t.Errorf("Unexpected error: %s", err)
	}

	if countryInfo != nil {
		t.Error("Expected nil country info")
	}
}

func TestUnexpectedStatusCodeUsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	mux.HandleFunc("/self", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "1;CA;CAN;Canada")
	})

	opts := NewOptions()
	opts.BaseURL = server.URL
	client := NewClientWithOptions(opts)

	countryInfo, err := client.LookupSelf()
	expectedError := "Unexpected response. Expected 200 but found 400"
	if err.Error() != expectedError {
		t.Errorf("Unexpected error: %s", err)
	}

	if countryInfo != nil {
		t.Error("Expected nil country info")
	}
}

func TestLookupIPv4UsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "1;AU;AUS;Australia")
	})

	opts := NewOptions()
	opts.BaseURL = server.URL
	client := NewClientWithOptions(opts)

	countryInfo, err := client.LookupIPv4("1.1.1.1")
	if err != nil {
		t.Fatal(err)
	}

	if countryInfo == nil {
		t.Error("Nil response from client")
	}
}

func TestLookupDecimalUsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "1;AU;AUS;Australia")
	})

	opts := NewOptions()
	opts.BaseURL = server.URL
	client := NewClientWithOptions(opts)

	// "1.1.1.1" = 16843009
	countryInfo, err := client.LookupDecimal(16843009)
	if err != nil {
		t.Fatal(err)
	}

	if countryInfo == nil {
		t.Error("Nil response from client")
	}
}

func TestLookupSelfUsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	mux.HandleFunc("/self", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "1;CA;CAN;Canada")
	})

	opts := NewOptions()
	opts.BaseURL = server.URL
	client := NewClientWithOptions(opts)

	countryInfo, err := client.LookupSelf()
	if err != nil {
		t.Fatal(err)
	}

	if countryInfo == nil {
		t.Error("Nil response from client")
	}
}

func TestErrWrongInputUsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "0;;;WRONG INPUT")
	})

	opts := NewOptions()
	opts.BaseURL = server.URL
	client := NewClientWithOptions(opts)

	countryInfo, err := client.LookupIPv4("a.b.c.d")
	if err != ErrWrongInput {
		t.Errorf("Expected ErrWrongInput, but found: %s", err)
	}

	if countryInfo != nil {
		t.Error("Expected nil country info for unknown IP")
	}
}

func TestErrUnknownUsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	mux.HandleFunc("/self", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "2;;;UNKNOWN")
	})

	opts := NewOptions()
	opts.BaseURL = server.URL
	client := NewClientWithOptions(opts)

	countryInfo, err := client.LookupSelf()
	if err != ErrUnknown {
		t.Errorf("Expected ErrUnknown, but found: %s", err)
	}

	if countryInfo != nil {
		t.Error("Expected nil country info for unknown IP")
	}
}

func TestTooManySegmentsUsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	mux.HandleFunc("/self", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "2;;;;Foo")
	})

	opts := NewOptions()
	opts.BaseURL = server.URL
	client := NewClientWithOptions(opts)

	countryInfo, err := client.LookupSelf()
	expectedError := "Invalid format. Expected 4 segments but found 5: \"2;;;;Foo\""
	if err.Error() != expectedError {
		t.Errorf("Unexpected error: %s", err)
	}

	if countryInfo != nil {
		t.Error("Expected nil country info for unknown IP")
	}
}

func TestInvalidCodeUsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	mux.HandleFunc("/self", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "5;CA;CAN;Canada")
	})

	opts := NewOptions()
	opts.BaseURL = server.URL
	client := NewClientWithOptions(opts)

	countryInfo, err := client.LookupSelf()
	expectedError := "Invalid format. Expected code of 0, 1, or 2 in 1st segment but found 5: \"5;CA;CAN;Canada\""
	if err.Error() != expectedError {
		t.Errorf("Unexpected error: %s", err)
	}

	if countryInfo != nil {
		t.Error("Expected nil country info for unknown IP")
	}
}

func TestInvalidTwoLetterCodeUsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	mux.HandleFunc("/self", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "1;CAN;CAN;Canada")
	})

	opts := NewOptions()
	opts.BaseURL = server.URL
	client := NewClientWithOptions(opts)

	countryInfo, err := client.LookupSelf()
	expectedError := "Invalid format. Expected 2 letter code in 2nd segment but found 3: \"1;CAN;CAN;Canada\""
	if err.Error() != expectedError {
		t.Errorf("Unexpected error: %s", err)
	}

	if countryInfo != nil {
		t.Error("Expected nil country info for unknown IP")
	}
}

func TestInvalidThreeLetterCodeUsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	mux.HandleFunc("/self", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "1;CA;CA;Canada")
	})

	opts := NewOptions()
	opts.BaseURL = server.URL
	client := NewClientWithOptions(opts)

	countryInfo, err := client.LookupSelf()
	expectedError := "Invalid format. Expected 3 letter code in 3rd segment but found 2: \"1;CA;CA;Canada\""
	if err.Error() != expectedError {
		t.Errorf("Unexpected error: %s", err)
	}

	if countryInfo != nil {
		t.Error("Expected nil country info for unknown IP")
	}
}

func TestInvalidFullNameUsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	mux.HandleFunc("/self", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "1;CA;CAN;")
	})

	opts := NewOptions()
	opts.BaseURL = server.URL
	client := NewClientWithOptions(opts)

	countryInfo, err := client.LookupSelf()
	expectedError := "Invalid format. Expected full name in 4th segment but found blank: \"1;CA;CAN;\""
	if err.Error() != expectedError {
		t.Errorf("Unexpected error: %s", err)
	}

	if countryInfo != nil {
		t.Error("Expected nil country info for unknown IP")
	}
}

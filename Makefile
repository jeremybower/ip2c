.PHONY: integration unit

integration:
	go test -tags=integration

unit:
	go test -tags=unit

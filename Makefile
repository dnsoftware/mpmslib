PROJECT="go-metrics"

default:
	echo ${PROJECT}

test:
	go test -v -count=1 ./...

.PHONY: default test cover
cover:
	go test -v -coverpkg=./... -coverprofile=coverage.out -covermode=count ./...
	go tool cover -func coverage.out | grep total | awk '{print $3}'

cover2:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	rm coverage.out

cover3:
	go test $(go list ./... | grep -v /proto) -coverprofile=coverage.out -covermode=count ./...
	go tool cover -func coverage.out | grep total | awk '{print $3}'

.PHONY: staticlint
staticlint:
	go build -o ./staticlint ./cmd/staticlint
	go vet -vettool=staticlint ./...


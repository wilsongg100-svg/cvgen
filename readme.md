## windowns

    go run cmd/main.go

## linux or mac

    make dev or make run

## test it:

    curl -X POST http://localhost:8080/api/v1/generate \
    -H "Content-Type: application/json" \
    -d @exemplo-cv.json \
    --output cv.pdf

## To do:

Add tests

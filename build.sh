#!/bin/sh
CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o tictac
docker build -t emilevauge/tictac .
docker push emilevauge/tictac:latest
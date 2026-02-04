#!/bin/bash
set -a
source .env
set +a
go build -o spotify-heardle main.go
exec ./spotify-heardle

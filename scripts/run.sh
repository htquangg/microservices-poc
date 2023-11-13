#!/bin/bash

set -e

service="$1"

cd "./internal/services/$service" && go run "./main.go"



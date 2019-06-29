#!/usr/bin/env bash

docker rm -f pharmeumpaymentapi_pharmeum-payment-api
docker rm -f pharmeumpaymentapi_pharmeum-payment-api-migrator

docker rmi -f pharmeumpaymentapi_pharmeum-payment-api
docker rmi -f pharmeumpaymentapi_pharmeum-payment-api-migrator
#!/bin/bash

VERSION=$(cat version)
docker push localhost:5000/go-user-service:${VERSION}
docker push localhost:5000/go-user-service:latest

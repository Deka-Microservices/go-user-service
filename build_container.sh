#!/bin/bash

VERSION=$(cat version)
docker build . --build-arg VERSION=${VERSION} --tag go-user-service:${VERSION}
docker tag go-user-service:${VERSION} go-user-service:latest
docker tag go-user-service:${VERSION} localhost:5000/go-user-service:${VERSION}
docker tag go-user-service:latest localhost:5000/go-user-service:latest

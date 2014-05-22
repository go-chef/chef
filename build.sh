#!/bin/sh

set -e

# Grab dependencies for coveralls.io integration
go get -u github.com/axw/gocov/gocov
go get -u github.com/mattn/goveralls

# Grab all project dependencies
go get -t -v ./...
go get
go build
go test -v ./...

# Define temp coverage file
t=coverme

# Overwrite the coverage file
go test -coverprofile=coverage

# Exclude our junk dirs
for dir in `find . -not \( -path './.*' -prune \) -not \( -path '*/test' -prune \) -type d`
do 
  go test -coverprofile=$t $dir
  # Get rid of the mode: set from each coverage run
  grep -v 'mode: set' $t >> coverage
done

go tool cover -func=coverage
goveralls -repotoken $COVERALLS_TOKEN -service drone.io -coverprofile=coverage
#!/bin/bash
project=$(go list -m)
echo "Building project $project"
if [ -d "build" ] 
then
    rm -rf build/*
else
	mkdir build/
fi
go build #-ldflags "-s -w"
echo "Copying files"
mv "$project" build/
cp -r data/* build/
cd build || exit 
#stty -icanon
./"$project" -f "$1"

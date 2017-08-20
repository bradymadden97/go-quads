#!/bin/bash

while IFS='' read -r line || [[ -n $line ]]; do
	pip install $line
done < "requirements.txt"

go build ./quads
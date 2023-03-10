#!/bin/sh -l

migrate -source "$1" -database "$2" up >> $GITHUB_OUTPUT
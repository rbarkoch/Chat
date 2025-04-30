#!/bin/bash

platforms=("linux/amd64" "windows/amd64" "darwin/arm64")
workspaceDirectory="/workspaces/Chat"
src="$workspaceDirectory/src"
bin="$workspaceDirectory/bin"

# MARK: LINUX
GOOS=linux 
GOARCH=amd64
go build -C ${src} -o ${bin}/${GOOS}/${GOARCH}/chat .

GOARCH=arm64
go build -C ${src} -o ${bin}/${GOOS}/${GOARCH}/chat .

# MARK: MAC OS
GOOS=darwin 
GOARCH=arm64
go build -C ${src} -o ${bin}/${GOOS}/${GOARCH}/chat .

#MARK: WINDOWS
GOOS=windows 
GOARCH=amd64
go build -C ${src} -o ${bin}/${GOOS}/${GOARCH}/chat.exe .
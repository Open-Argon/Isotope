@echo off
go build -trimpath -ldflags="-s -w" -o bin/isotope.exe ./src
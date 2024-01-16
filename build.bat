
@echo off
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=386
go build  -o TxtReSearch32.exe

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build  -o TxtReSearch.exe 

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=386
go build  -o linuxTxtReSearch32 

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build  -o linuxTxtReSearch 

SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build  -o macTxtReSearch

SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=386
go build  -o macTxtReSearch32
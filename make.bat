set GOPATH  = %MEGAM_GOPATH%:%CD%\..\..\..\..\

:: MEGAM_GOPATH This is where gulp project placed.

set datetimef=%date:~-4%_%date:~3,2%_%date:~0,2%__%time:~0,2%_%time:~3,2%_%time:~6,2%
echo %datetimef%
echo %GOPATH%
rm -rf %MEGAM_GOPATH%
go get %GO_EXTRAFLAGS% -u -d -t -insecure ./...
rm -f gulpd
go build %GO_EXTRAFLAGS% -o gulpd.exe ./cmd/gulpd

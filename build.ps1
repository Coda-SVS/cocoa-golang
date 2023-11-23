$Currentlocation = Get-Location

# $LibPath = Join-Path($Currentlocation) "lib\ffmpeg-master-latest-win64-gpl-shared"

# $Env:Path += ";$LibPath"

& "go" @("build", "-ldflags", "-s -w -H=windowsgui -extldflags=-static", ".\src\cmd\cocoa")
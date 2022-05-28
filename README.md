# dressme-scrapper

## installation

    go mod init dressme-scrapper
    export GOPATH=~/home/olivier/dressme
    go get .
then move code inside `dressme-scrapper/src`

## relative export/import
Start with a capital letter

export:

    package <folder_name>

    func Func() {}
    var Variable = ""

import:

    import (
        "dressme-scrapper/src/../<folder_name>"
    )

    <folder_name>.Func
    <folder_name>.Variable
## run

    go run .

## build

    go build .
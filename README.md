# scrapper

## run

    go run src/main.go

## build

    go build .

## License

must share photos generated with https://creativecommons.org/licenses/by-sa/2.0/

## .env

MONGODB_URI=mongodb://localhost:27017
SCRAPPER_DB=scrapper
UNWANTED_TAGS_COLLECTION=unwantedTags
WANTED_TAGS_COLLECTION=wantedTags
FLICKR_COLLECTION=flickr
FLICKR_PATH=***
PRIVATE_KEY=***
PUBLIC_KEY=***

## installation

    go mod init scrapper
    go get .
then move code inside `scrapper/src`

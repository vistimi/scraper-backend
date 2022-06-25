# scrapper

Online scrapper for building a dataset for ML.
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
    IMAGE_PATH=***
    FLICKR_COLLECTION=flickr
    FLICKR_PRIVATE_KEY=***
    FLICKR_PUBLIC_KEY=***
    UNSPLASH_COLLECTION=unsplash
    UNSPLASH_PRIVATE_KEY=***
    UNSPLASH_PUBLIC_KEY=***

## installation

    go mod init scrapper
    go get .

Then move code inside `scrapper/src`

## linter

https://github.com/mgechev/revive

    revive -config revive.toml

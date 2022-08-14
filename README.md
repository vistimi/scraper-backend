# scraper

Online scraper for building a dataset for ML.


## installation

Install Golang and MongoDB

If pbm with package `<package>: command not found`:

    export GOPATH="$HOME/go"
    PATH="$GOPATH/bin:$PATH"

## Docker

sudo docker build -t scraper-img .
sudo docker run -it --rm --name scraper-run --env-file .env scraper-img

    
## run without docker

    go run src/main.go

## build without docker

    go build -o scraper src/main.go

## License

must share photos generated with https://creativecommons.org/licenses/by-sa/2.0/

## .env

Remove in src/utils/env.go the godotenv part if you run in container.

    MONGODB_URI=mongodb://localhost:27017
    SCRAPER_DB=scraper
    TAGS_UNWANTED_COLLECTION=tagsUnwanted
    TAGS_WANTED_COLLECTION=tagsWanted
    IMAGES_WANTED_COLLECTION=imagesWanted
    IMAGES_PENDING_COLLECTION=imagesPending
    IMAGES_UNWANTED_COLLECTION=imagesUnwanted
    USERS_UNWANTED_COLLECTION=usersUnwanted
    IMAGES_BUCKET=<s3_bucket_name>
    FLICKR_PRIVATE_KEY=***
    FLICKR_PUBLIC_KEY=***
    UNSPLASH_PRIVATE_KEY=***
    UNSPLASH_PUBLIC_KEY=***
    PEXELS_PUBLIC_KEY=***

## linter

https://github.com/mgechev/revive

    revive -config revive.toml



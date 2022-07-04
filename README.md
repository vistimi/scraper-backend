# scraper

Online scraper for building a dataset for ML.


## installation

Install Golang and MongoDB

    git clone git@github.com:KookaS/scraper.git

    
## run

    go run src/main.go

## build

    go build .

## License

must share photos generated with https://creativecommons.org/licenses/by-sa/2.0/

## .env

    MONGODB_URI=mongodb://localhost:27017
    SCRAPER_DB=scraper
    TAGS_WANTED_COLLECTION=tagsUnwanted
    TAGS_WANTED_COLLECTION=tagsWanted
    IMAGES_COLLECTION=images
    IMAGES_UNWANTED_COLLECTION=imagesUnwanted
    USERS_UNWANTED_COLLECTION=usersUnwanted
    IMAGE_PATH=<absolutePath>
    FLICKR_PRIVATE_KEY=***
    FLICKR_PUBLIC_KEY=***
    UNSPLASH_PRIVATE_KEY=***
    UNSPLASH_PUBLIC_KEY=***
    PEXELS_PUBLIC_KEY=***

## linter

https://github.com/mgechev/revive

    revive -config revive.toml

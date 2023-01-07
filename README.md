# scraper

Online scraper for building a dataset for ML.

## installation

VSCode and Docker

## Run

Run in devcontainer the backend and spawn in other terminals the localstack and mongodb.

- ***Create network for docker*** (otherwise it will fail)
- MongoDB
- Localstack
- Backend

```shell
# network
docker network create scraper-net
docker network ls

# docker images
docker run --rm -it --net scraper-net --name scraper-localstack localstack/localstack
docker run --rm -it --net scraper-net --name scraper-mongodb mongo:6.0.1

# To test the connection, should not throw an error
curl --connect-timeout 10 --silent --show-error scraper-mongodb:27017
curl --connect-timeout 10 --silent --show-error scraper-localstack:4566
```

#### Backend with Docker
```shell
sudo docker build -t scraper-img .
sudo docker run --rm -it --net scraper-net --name scraper-run --env-file <state>.env scraper-img
```

#### Backend without docker
    go run src/main.go

### Build without docker
    go build -o scraper src/main.go
    ./scraper

## License

must share photos generated with https://creativecommons.org/licenses/by-sa/2.0/

## Env

Create a local.env file:

    ENV=local
    MONGODB_URI=mongodb://scraper-mongodb:27017
    LOCALSTACK_URI=http://scraper-localstack:4566
    SCRAPER_DB=scraper
    TAGS_UNDESIRED_COLLECTION=tagsUndesired
    TAGS_DESIRED_COLLECTION=tagsDesired
    PRODUCTION=imagesProduction
    PENDING=imagesPending
    UNDESIRED=imagesUndesired
    VALIDATION=imagesValidation
    USERS_UNDESIRED_COLLECTION=usersUndesired
    IMAGES_BUCKET=scraper-backend-test-env
    FLICKR_PRIVATE_KEY=***
    FLICKR_PUBLIC_KEY=***
    UNSPLASH_PRIVATE_KEY=***
    UNSPLASH_PUBLIC_KEY=***
    PEXELS_PUBLIC_KEY=***

ENV is either `production`, `staging`, `development` or `local`

## linter

https://github.com/mgechev/revive

    revive -config revive.toml

## Dependencies

    go mod tidy


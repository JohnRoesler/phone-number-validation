# phone number validation

todo: what it does

## usage

1. Start the server 

  - option 1: go run

    ```
    go run main.go
    ```

  - option 2: go build

    ```
    go build main.go
    ./main
    ```

  - option 3: docker

    ```
    docker build -t phone-number-validation .
    docker run -p 8080:8080 phone-number-validation
    ```

2. Hit the api

```
curl "localhost:8080/v1/phone-numbers?phoneNumber=%2B12125690123"
# success

curl "localhost:8080/v1/phone-numbers?phoneNumber=631%20311%208150"
# error

```

## local development

*requirements*

- go >= 1.19 installed - [documentation](https://go.dev/doc/install)

### run tests

```
go test -v ./...
```

## Additional considerations

- Why go? Go is simplistic in nature making it a good language for readability and maintenance. 
- Why [nyaruka/phonenumbers](github.com/nyaruka/phonenumbers)? The library provided the requirements for validating phone numbers and country codes and has a reasonable number of stars as well as recent commits.
- How would this deploy to production? Assuming a Kubernetes environment, the proper config files would be added either raw, or using a tool like helm to managed the deployment. Ideally there would also be CI such as GitHub Actions to run tests automatically and any other necessary validations prior to deployment to a development environment. After validation in dev (whether automated or manual) the application would be released into production via some trigger e.g. a GitHub release/tag.
- Any assumptions? That it was ok to ignore extra white space in between numbers since the library used handled that without issue.
- What improvements would I make? Metrics and logging. Graceful shutdown of the HTTP server. OpenAPI documentation (could consider using [goa.design](goa.design)). If action is required for invalid phone numbers, perhaps a mechanism to place those in a DB or message queue for later investigation.

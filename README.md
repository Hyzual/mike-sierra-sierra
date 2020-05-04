# mike-sierra-sierra

More soon!

### Contributing

#### Start the dev Docker container

```sh
$ make start
# Then, access http://localhost:8080
```

#### Go commands

```sh
# Run unit tests
$ make test
# Run golint
$ make lint
# Generate html test coverage
$ make coverhtml
```

#### Build the production Docker image

```sh
# Build the image
$ make build-docker-image
# Run goss tests on the built image to ensure everything keeps working
$ make dgoss-ci
```

#### How to edit the goss.yaml file interactively

```sh
# It will run the container with goss and goss.yaml inside
$ dgoss edit --user=$(id -u) hyzual/mike-sierra-sierra
# Once in the container, you can run goss commands
[container]$ goss autoadd
```

#### Run stylelint

```sh
$ npm run stylelint -- ./assets
# To automatically fix problems
$ npm run stylelint -- --fix ./assets
```

#### Run prettier

```sh
# To automatically format HTML templates and CSS assets
$ npm run prettier -- --write ./templates ./assets
```

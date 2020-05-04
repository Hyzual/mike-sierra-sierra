# mike-sierra-sierra

More soon!

### Contributing

#### Go commands

```sh
$ make build
# Start the server
$ make start
# Remove the binary
$ make clean
# Run unit tests
$ make test
# Run golint
$ make lint
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
container$ goss autoadd
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

# mike-sierra-sierra

More soon!

### Contributing

#### Enable HTTPS support

Download and install [mkcert](https://github.com/FiloSottile/mkcert#installation)
Then, run the following commands:

```sh
$ mkcert -install
$ mkcert localhost ::1
$ mv localhost+1.pem ./certs/cert.pem
$ mv localhost+1-key.pem ./certs/key.pem
```

> Warning: the rootCA-key.pem file that mkcert automatically generates gives complete power to intercept secure requests from your machine. Do not share it.

#### Start the dev Docker container

```sh
$ make start
# Then, access https://localhost:8443
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

#### Database

Mike-sierra-sierra uses [SQLite](https://www.sqlite.org) to persist data.
The database file is located at ./database/file/mike.db
If it does not exist (it is gitignored), the server will create it on first run.

#### Add music files to the dev Docker image

```sh
# Run an alpine container with bind-mount to a folder where you have music (on your host)
$ docker run -it --rm -v /home/<my-user>/Music:/source -v mike_music:/dest alpine ash
[container]$ cp /source/*.mp3 dest
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
# It will run the container with goss and goss.yaml inside. --user allows you to edit the goss.yaml file.
# The env variable disables HTTPS (otherwise you need to provide valid cert and key)
$ dgoss edit -e MIKE_DISABLE_HTTPS=1 --user=$(id -u) hyzual/mike-sierra-sierra
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

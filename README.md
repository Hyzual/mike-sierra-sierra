# mike-sierra-sierra

More soon!

### Contributing

#### Install Nix

Only two dev tools are required: `git` to clone the sources from GitHub and `nix`. Please see the appropriate installation instructions from <https://nixos.org/download.html#nix-quick-install>.

Then, run the following command that will download all dev tools (the first time) and open a shell with all dev tools installed:

```sh
$ nix-shell
```

#### Enable HTTPS support

Run the following commands:

```sh
$ mkcert -install
$ mkcert localhost ::1
$ mv localhost+1.pem ./secrets/cert.pem
$ mv localhost+1-key.pem ./secrets/key.pem
```

> Warning: `~/.local/share/mkcert/rootCA-key.pem` file that mkcert automatically generates gives complete power to intercept secure requests from your machine. Do not share it.

#### Start the dev Docker container

```sh
# build and watch the Go sources
$ make start
# Build and watch the CSS and Javascript
$ make watch
# Then, access https://localhost:8443
```

#### Database

Mike-sierra-sierra uses [SQLite](https://www.sqlite.org) to persist data.
The database file is located at `./database/file/mike.db`
You must create the database file yourself and run the [`./database/schema/install.sql`](file://./database/schema/install.sql) script on it in order to create the DB tables.

#### First-time registration

Upon first starting the container, go to https://localhost:8443/first-time-registration to register your admin user.

#### Go commands

```sh
# Run unit tests
$ make test-go
# Run golint
$ make lint-go
# Generate html test coverage
$ make coverage-go-html
```

#### NPM commands

```sh
# Build and minify the CSS and Javascript (for production)
$ make build-assets
# Run unit tests
$ make test-jest
# Display text test coverage
$ make coverage-jest
# Build and watch (for development)
$ make watch
```
s
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
$ npm run stylelint -- ./styles
# To automatically fix problems
$ npm run stylelint -- --fix ./styles
```

#### Run eslint

```sh
$ npm run eslint -- .
# To automatically fix problems
$ npm run eslint -- --fix .
```

#### Run prettier

```sh
# To automatically format HTML templates, Typescript files and CSS assets
$ npm run prettier -- --write ./templates ./scripts ./styles
```

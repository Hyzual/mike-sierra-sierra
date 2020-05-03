# mike-sierra-sierra

More soon!

### Contributing

#### Build the production Docker image

```sh
# Build the image
$ docker build -t hyzual/mike-sierra-sierra .
# Run goss tests on the built image to ensure everything keeps working
$ dgoss run hyzual/mike-sierra-sierra
```

#### How to edit the goss.yaml file interactively

```sh
$ dgoss edit --user=$(id -u) hyzual/mike-sierra-sierra
# It will run the container with goss and goss.yaml inside
# Once in the container, you can run goss commands
container$ goss autoadd
```

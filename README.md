# lokahi

Lokahi is a HTTP uptime and response time monitoring service. Each check has a webhook URL that will have check information POST-ed to it when the service status changes.

## Building

### Binaries

```console
$ mage build
```

Binaries for the current `GOARCH` and `GOOS` combination will be put into 
`./bin`.

### Docker

```console
$ mage docker
```

The docker image `xena/lokahi` will be built.

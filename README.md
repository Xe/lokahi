# lokahi

Lokahi is a HTTP uptime and response time monitoring service that also keeps 
track of response time for each check in a histogram, allowing for more 
complicated analysis of data.

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

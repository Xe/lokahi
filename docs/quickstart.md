# Quickstart

The quickest way to get lokahi set up is to install the development tools and
then launch lokahi in docker using docker compose.

## Development

Install the following things on your computer:
 - Docker
 - Go (any version, preferably the version listed in `mage.go`)
 - [Docker Compose](https://docs.docker.com/compose/install/)
 - [retool](https://github.com/twitchtv/retool)

Build the generators and tool dependencies:

```console
$ retool build
```

Build all of the command line tools:

```console
$ retool do mage generate build
```

## Running

```console
$ retool do mage run
```

```console
$ ./bin/lokahictl
See https://github.com/Xe/lokahi for more information

Usage:
  lokahictl [command]

Available Commands:
  create      creates a check
  delete      deletes a check
  get         dumps information about a check
  help        Help about any command
  list        lists all checks that you have permission to access
  put         puts updates to a check
  run         runs a check
  runstats    gets performance information

Flags:
  -h, --help            help for lokahictl
      --server string   http url of the lokahid instance (default "http://AzureDiamond:hunter2@127.0.0.1:24253")

Use "lokahictl [command] --help" for more information about a command.
```

### create a check

```console
$ ./bin/lokahictl create -e 60 -u http://duke:9001 \
  -w http://samplehook:9001/twirp/github.xe.lokahi.Webhook/Handle
{
  "id": "dc30dc6d-fcb1-49d1-add6-c6d63e337c56",
  "url": "http://duke:9001",
  "webhook_url": "http://samplehook:9001/twirp/github.xe.lokahi.Webhook/Handle",
  "every": 60
}
```

### get a check

```console
$ ./bin/lokahictl get dc30dc6d-fcb1-49d1-add6-c6d63e337c56
{
  "id": "dc30dc6d-fcb1-49d1-add6-c6d63e337c56",
  "url": "http://duke:9001",
  "webhook_url": "http://samplehook:9001/twirp/github.xe.lokahi.Webhook/Handle",
  "every": 60,
  "state": 2
}
```

### list checks

```console
$ ./bin/lokahictl list
{
  "results": [
    {
      "check": {
        "id": "dc30dc6d-fcb1-49d1-add6-c6d63e337c56",
        "url": "http://duke:9001",
        "webhook_url": "http://samplehook:9001/twirp/github.xe.lokahi.Webhook/Handle",
        "every": 60,
        "state": 2
      }
    }
  ]
}
```

## Integrations

Each check has a webhook. The webhook URL given will be HTTP post-ed with a
protobuf-encoded [CheckStatus](https://github.com/Xe/lokahi/blob/master/rpc/lokahi/lokahi.proto#L83)
with the following additional headers:
- `Content-Type: application/protobuf`
- `Accept: application/protobuf`
- `User-Agent: lokahi/dev (+https://github.com/Xe/lokahi)`

See [samplehook](https://github.com/Xe/lokahi/blob/master/cmd/sample_hook/main.go) 
for a simple example on how to receive this in the form of a twirp service.

Any data that the webhook handler returns will be ignored. Return non-2xx if
the service is not healthy.

Webhooks will only be POST-ed to when the state of a service changes.

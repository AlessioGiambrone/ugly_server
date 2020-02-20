# Ugly server

[![Build Status](https://travis-ci.com/AlessioGiambrone/ugly_server.svg?branch=master)](https://travis-ci.com/AlessioGiambrone/ugly_server)

`ugly_server` is a reverse proxy that permits some basic manipulation of
request parameters.

Its aim is to permit the proxied server's caching to work efficiently.

## Example

```bash
$ CONFIG=config.yaml  ugly_server
2020/02/17 21:34:57 Serving http://localhost:8000 at port :7072
```

## Configuration

Configuration is done via YAML file.

By default, `ugly_server` search for `config.yaml` in its own directory or
in `./conf`, but can be specified another path via `CONFIG` environment
variable.

Example:

```yaml
port: 7072                       # the listening port for ugly_server
proxiedService: "localhost:5000" # the remote server we're proxying
constraints:                     # a list of parameters that we want to manipulate
  "lat":                         # for example, we put some constraint to "lat"
    round: 5                     # round at the 5th decimal
    max: 47                      # topping at 47
    min: -45                     # minimum value -45
```

## Docker

You can try `ugly_server` using docker:

```bash
docker run --rm -v $(pwd)/conf/:/conf/ --net host ugly_server
```

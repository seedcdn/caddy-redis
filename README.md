## Caddy Redis Adapter
This is a [config adapter](https://caddyserver.com/docs/config-adapters) for Caddy which uses Redis to store and updates configuration based on event or periodically.

### Install
First, ensure your GOROOT and GOPATH environment variables are correct for your environment.

Then, follow the xcaddy install process [here](https://github.com/caddyserver/xcaddy#install).

Then, build Caddy with this Go module plugged in. For example:

```bash
$ xcaddy build --with github.com/seedcdn/caddy-redis
```

### Use
Using this config adapter is the same as all the other config adapters.

* [Learn about config adapters in the Caddy docs](https://caddyserver.com/docs/config-adapters)
* You can adapt your config with the [`adapt` command](https://caddyserver.com/docs/command-line#caddy-adapt)

### Disclaimer
This project is not affiliated with Caddy or Redis.

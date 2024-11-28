# Janus Cloudflare TURN Proxy

A proxy to transform the response of a Cloudflare TURN request into a format Janus understands.

## The problem

[Janus](https://github.com/meetecho/janus-gateway) is a WebRTC gateway. It can dynamically
connect to TURN servers by calling a REST endpoint. Janus expects a json response as defined
in their [Internet Draft](https://datatracker.ietf.org/doc/html/draft-uberti-behave-turn-rest-00#section-2.2):

```json
{
    "username" : "12334939:mbzrxpgjys",
    "password" : "adfsaflsjfldssia",
    "ttl" : 86400,
    "uris" : [
        "turn:1.2.3.4:9991?transport=udp",
        "turn:1.2.3.4:9992?transport=tcp",
        "turns:1.2.3.4:443?transport=tcp"
    ]
}
```

Cloudflare, however, generates a json response similar to the following:

```json
{
  "iceServers": {
    "urls": [
      "stun:stun.cloudflare.com:3478",
      "turn:turn.cloudflare.com:3478?transport=udp",
      "turn:turn.cloudflare.com:3478?transport=tcp",
      "turns:turn.cloudflare.com:5349?transport=tcp"
    ],
    "username": "username",
    "credential": "password"
  }
}
```

This project aims to fill the gap by acting as a proxy. Janus calls the proxy endpoint and receives the response from
Cloudflare transformed as expected.

## Usage

Recommended deployment strategy is to use the provided Docker image:

```yml
services:
  turn-proxy:
    image: ghcr.io/actcoding/janus-cloudflare-turn-proxy:latest
    environment:
      CF_TURN_KEY_ID: <CLOUDFLARE_TURN_SERVER_SECRET>
```

Replace `<CLOUDFLARE_TURN_SERVER_SECRET>` with the secret generated in the Cloudflare Dashboard. Detailed documentation
can be found over there: https://developers.cloudflare.com/calls/turn/generate-credentials/

The proxy exposes an HTTP endpoint on port `1337` which Janus can be pointed to.

## Configuration

The proxy is configured using environment variables. Below you can find an overview
of all available options.

| Environment Variable | Default? | Description |
| --- | --- | --- |
| `SERVER_ADDRESS` | `0.0.0.0` | The address the server is listening on. |
| `SERVER_PORT` | `1337` | The port the server is listening on. |
| `CF_TTL` | `86400` | Credentials lifetime in seconds. Defaults to 24 hours. |
| `CF_TURN_KEY_ID` | **Required** | The Cloudflare TURN key id. Not to be confused with the secret mentioned earlier. |

## License

[MIT](LICENSE)

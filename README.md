# Project Zomboid Info API

An HTTP API that fetches and displays information from a Project Zomboid server using RCON (remote console). The API provides details such as the server's online status, players connected, max players allowed, and more. This project includes caching mechanisms for faster responses and Docker support for easy deployment.

## Table of Contents

- [API Endpoint](#api-endpoint)
- [Environment Variables](#environment-variables)
- [Docker Usage](#docker-usage)
- [Caching Explained](#caching-explained)
- [Local Development](#local-development)
- [License](#license)

## API Endpoint

- **GET /**: Returns a JSON with server status and details. Example response:

  ```json
  {
    "status": "Online",
    "players": ["player1", "player2"],
    "onlinePlayerCount": 2,
    "maxPlayers": 10,
    "playersString": "2/10",
    "publicName": "My Project Zomboid Server",
    "publicDescription": "This is a public description."
  }
  ```

* If something goes wrong or the server is offline, the response will be `{ "status": "Offline" }` and other fields may be empty.

## Environment Variables

Create a `.env` file (or set these in your environment directly) with the following variables:

| Variable            | Description                                    | Default    |
| ------------------- | ---------------------------------------------- | ---------- |
| `RCON_SOCKET`       | The IP and port for the RCON connection.       | (Required) |
| `RCON_PASSWORD`     | RCON password for the server.                  | (Required) |
| `CACHE_TTL_SECONDS` | How long (in seconds) the cache remains valid. | `5`        |

**Note**: If using Docker Compose, you can reference these with an `env_file` directive or an `environment` directive in the `compose.yml` file.

## Docker Usage

* The image is available on Docker Hub: [beyenilmez/pz-info-api](https://hub.docker.com/r/beyenilmez/pz-info-api)

* It is also available on GHCR: [ghcr.io/beyenilmez/pz-info-api](https://ghcr.io/beyenilmez/pz-info-api)

Using docker compose:

```bash
services:
  pz-info-api:
    image: ghcr.io/beyenilmez/pz-info-api:latest
    container_name: pz-info-api
    environment:
      - RCON_SOCKET=127.0.0.1:27015
      - RCON_PASSWORD=mysecret
      - CACHE_TTL_SECONDS=5
    ports:
      - "8080:8080"
    restart: unless-stopped
```

or docker run:

```bash
docker run -d --name pz-info-api \
  -e RCON_SOCKET=127.0.0.1:27015 \
  -e RCON_PASSWORD=mysecret \
  -e CACHE_TTL_SECONDS=5 \
  -p 8080:8080 \
  --restart unless-stopped \
  ghcr.io/beyenilmez/pz-info-api:latest
```

## Caching Explained
* The response is cached for a duration configured by `CACHE_TTL_SECONDS`.
* During that time, repeated requests to / serve the same JSON without hitting RCON again.
* Once the cache expires, the API makes new RCON calls to update the response.

## Local Development

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/beyenilmez/pz-info-api.git
   cd pz-info-api
   ```
2. **Set Up `.env`**:
   ```bash
   cp .env.example .env
   ```
3. **Install Dependencies**:
   ```bash
   go mod download
   ```
4. **Run the Server**:
   ```bash
   go run ./cmd/server
   ```
5. **Access the API**:
   ```bash
   curl http://localhost:8080
   ```

## License

Distributed under the MIT License. See [LICENSE](https://github.com/beyenilmez/pz-info-api/blob/main/LICENSE) for more information.
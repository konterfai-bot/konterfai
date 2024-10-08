[![Matrix](https://img.shields.io/matrix/konterfai%3Amatrix.org?logo=matrix)](https://app.element.io/#/room/#konterfai:matrix.org)
[![License](https://img.shields.io/badge/License-AGPL-v3.svg)](LICENSE)
[![status-badge](https://ci.codeberg.org/api/badges/13605/status.svg)](https://ci.codeberg.org/repos/13605)
[![Go Report Card](https://goreportcard.com/badge/codeberg.org/konterfai/konterfai)](https://goreportcard.com/report/codeberg.org/konterfai/konterfai)
[![codecov](https://codecov.io/github/konterfai-bot/konterfai/graph/badge.svg?token=25T5Y462Q1)](https://codecov.io/github/konterfai-bot/konterfai)
[![Codeberg Issues](https://img.shields.io/gitea/issues/open/konterfai/konterfai?gitea_url=https%3A%2F%2Fcodeberg.org)](https://codeberg.org/konterfai/konterfai)
[![Codeberg Release](https://img.shields.io/gitea/v/release/konterfai/konterfai?gitea_url=https%3A%2F%2Fcodeberg.org&sort=semver)](https://codeberg.org/konterfai/konterfai/releases)

# konterfAI

(c) 2024 konterfAI

konterfAI is a proof-of-concept for a model-poisoner for LLM (Large Language Models) to generate nonsense("bullshit")
content suitable to degenerate these models.

Although it's still work in progress and not yet ready for production, it already shows the concept of fighting fire
with fire: The backend queries a tiny
LLM running in [ollama](https://ollama.com/) with a high ai-temperature setting to generate hallucinatory content.
If you wonder how this looks like, check out
the [example-hallucination.md](https://codeberg.org/konterfai/konterfai/src/branch/main/docs/example-hallucination.md)
file.

**NOTE:** The developers created konterfAI not as an offensive (hacking) tool, but a countermeasure against AI-crawlers
that ignore robots.txt and other rules. The Tool was inspired by reports of web admins suffering from TeraByte of Data
caused by AI crawlers - cost that can be avoided.

## License

konterfAI is licensed under the AGPL (GNU AFFERO GENERAL PUBLIC LICENSE). See [LICENSE](LICENSE) for the full license
text.

## Get in touch

Join the [Matrix-Chat](https://app.element.io/#/room/#konterfai:matrix.org) to get in touch.

## Contributing

see [CONTRIBUTING](https://codeberg.org/konterfai/konterfai/src/branch/main/docs/contributing.md).

## FAQ (Frequently Asked Questions)

see [FAQ](https://codeberg.org/konterfai/konterfai/src/branch/main/docs/faq.md).

## How does it work?

konterfAI is supposed to run behind a reverse-proxy, like nginx or traefik.
The reverse proxy needs the ability to detect the user-agent of the incoming request and filter it by a given list.
If there is a match the crawler will not be presented with the original content, but with the poisoned content.
The poisoned content is also cluttered with randomized self-references to catch the crawlers in some kind of tar-pit.

![A diagram showing the basic concept of konterfAI](https://codeberg.org/konterfai/konterfai/raw/branch/main/docs/img/how_it_works.png)

**Note:** Those are examples and not intended for copy & paste usage. Make sure to read them carefully and adjust them
to your needs.

## What ollama models does konterfAI ship?

None, konterfAI does not ship any models. A default model will be downloaded upon ollama start.
If you want to use a different model, you can pick one from the [ollama-models](https://ollama.com/models) page and
adapt your configuration accordingly.

## Building

```bash
$> make build
```

For a full list of build targets see [Makefile](https://codeberg.org/konterfai/konterfai/src/branch/main/Makefile).

## How to run it?

### Production deployment

If you are really brave and want to try konterfAI in a production environment, see there are two examples for
[nginx](https://codeberg.org/konterfai/konterfai/src/branch/main/deployments/nginx)
and [traefik](https://codeberg.org/konterfai/konterfai/src/branch/main/deployments/traefik) in the deployment-folder.

**Note:** These examples are not intended for copy & paste usage.
Make sure to read them carefully and adjust them to your needs.

**WARNING:** IMPROPER CONFIGURATION WILL HAVE NEGATIVE EFFECTS ON YOUR SEO

### Development

**Note:** `-gpu` is optional, if you do not have an ollama-capable GPU, you can omit it.

```bash
$> make start-ollama[-gpu]
$> make run
```

### Tracing

see [Tracing](https://codeberg.org/konterfai/konterfai/src/branch/main/docs/tracing.md).

### Default Ports

konterfAI will start two webservers, one is the service itself, listening on port 8080.
The other is the statistics server, listening on port 8081. If you are running this locally from source,
you can access both servers via [http://localhost:8080](http://localhost:8080) and [http://localhost:8081](http://localhost:8081).
These ports can be changed via the `--port` and `--statistics-port` flags.

### Prometheus Metrics

konterfAI exposes prometheus metrics on the `/metrics` endpoint from the statistics server.
You can access them via [http://localhost:8081/metrics](http://localhost:8081/metrics).

### Docker

**Start:**

```bash
$> make start-ollama[-gpu]
$> make docker-build
$> make docker-run
```

**Stop:**

```bash
$> make docker-stop
```

### Docker-Compose

**Start:**

```bash
$> make docker-compose-up
```

**Stop:**

```bash
$> make docker-compose-down
```

For more complex examples
edit [docker-compose.yml](https://codeberg.org/konterfai/konterfai/src/branch/main/docker-compose-dev.yml) to suit your
needs.

#### Pre-built Docker-Image

You can also use the pre-built docker-image from [Docker-Hub](https://hub.docker.com/r/konterfai/konterfai)
or [Quay.io](https://quay.io/repository/konterfai/konterfai).

#### Pre-built Docker-Images tags

| Tag           | Description                           |
|---------------|---------------------------------------|
| `latest`      | The latest stable release             |
| `v*.*.*`      | A specific version  (e.g. 0.1.0, 0.1) |
| `latest-main` | The latest main branch build          |

## Configuration

konterfAI is configured via cli-flags. For a full list of supported flags run (after [building](#building)):

```bash
$> ./bin/konterfai --help
```

The docker-image is configured via environment-variables.
For a full list of supported variables
see [docker-compose.yml](https://codeberg.org/konterfai/konterfai/src/branch/main/docker-compose-dev.yml)
and [entrypoint.sh](https://codeberg.org/konterfai/konterfai/src/branch/main/entrypoint.sh).

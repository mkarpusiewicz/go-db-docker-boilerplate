# Go + DB + Docker Boilerplate

## Configuration

Configuration is done via environemnt variables.

For local dev they are stored in `.env` file in repo root, for local docker build in `build/local/.env`

### Variables

- `PROJECT` - project name, used for docker
- `SERVER_PORT` - rest api http server port

For local dev also:

- `APP_ENV=local` - so that app not runs in `production` config, e.g. removed debug logs, endpoints, etc.

For integration tests:

- `SERVER_URL` - url of rest api http server to be tested

## Docker files

- `src.Dockerfile` - creates base image with src and dependencies, used by `build` and `int`
- `build.Dockerfile` - builds app image
- `int.Dockerfile` - builds integration tests image

## Deployment

Docker compose files are located in `build` directory:

- `docker-compose.yml` - contains the application
- `docker-compose.dependencies.yml` - contains local dependencies to be substitued by infrastructure on provider
- `docker-compose.integration.yml` - contains integration tests docker container

## Requirements

- docker, docker-compose
- make
- git
- curl (optional)

### Windows

- [https://docs.docker.com/docker-for-windows/install/](https://docs.docker.com/docker-for-windows/install/)
- [http://gnuwin32.sourceforge.net/packages/make.htm](http://gnuwin32.sourceforge.net/packages/make.htm)
- [https://git-scm.com/](https://git-scm.com/)
- [https://curl.haxx.se/windows/](https://curl.haxx.se/windows/)

or using [chocolatey](https://chocolatey.org/):

```cmd
choco install docker-desktop make git curl
```

or using [scoop](https://scoop.sh/)

```cmd
scoop install docker make git curl
```

## Building and running

To install dependencies for local dev: `make deps`

To build and run locally:

```sh
make dev
```

To run locally with hot reload:

```sh
make dev-hot
```

To build docker image:
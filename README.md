# GOLocator

This repository contains system made with GO, using two independent microservices, for managing client location.

<hr>

## Preparing for Development

1. Ensure `go` 1.22 or above is installed.
2. Ensure `git` is installed.
3. Clone repository: `git clone https://github.com/JustFiesta/GOLocator`
4. Create containers for each service using provided dockerfiles: `docker build usr_service`, `docker build location_service`
5. Build cli:

```shell
cd cmd/goloc
go build -o goloc
```

6. Install program: `go install`

OR

Run `setup.sh` script after cloning reposiory

<hr>

### Usage

* Update current user location - pass the username and location

```shell
goloc update -u <user_name> -c <location>
```

* Search for users in some location within the provided radius

```shell
goloc search -c <location> -r <radius>
```

* Return distance traveled by a user within some date/time range. Time range defaults to 1 day.

```shell
goloc travel -t <YYYY-MM-DDTHH:MM:SS+UTC>
```

<hr>

#### Running Tests

```shell
go test
```

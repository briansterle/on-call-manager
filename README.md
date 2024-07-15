# On Call Manager

The On Call manager is implemented using Go, [htmx](https://htmx.org/) and [bootstrap](https://getbootstrap.com/)

## Run the app

### pre-reqs
- [docker](https://docs.docker.com/engine/install/)
- [go](https://go.dev/doc/install)

```bash
# start postgres
docker run -p 5432:5432 -e POSTGRES_PASSWORD=foobar -d postgres 

# start server
./run.sh
```

webserver runs on localhost:8080

## Install PWA on iOS
- navigate to localhost:8080 in safari
- open share sheet
- click "Add to Home Screen"

# Dev Guide

## No DB mode

When you are working on dex it's convenient to use the `--no-db` flag. This starts up dex in a mode which uses an in-memory datastore for persistence. It also does not rotate keys, so no overlord is required.

In this mode you provide the binary with paths to files for connectors, users, and emailer. There are example files you can use inside of `static/fixtures` named *"connectors.json.sample"*, *"users.json.sample"*, and *"emailer.json.sample"*, respectively.

You can rename these to the equivalent without the *".sample"* suffix since the defaults point to those locations:

```console
mv static/fixtures/connectors.json.sample static/fixtures/connectors.json
mv static/fixtures/users.json.sample static/fixtures/users.json
mv static/fixtures/emailer.json.sample static/fixtures/emailer.json
```

Starting dex is then as simple as:

```console
bin/dex-worker  --no-db
```

***Do not use this flag in production*** - it's not thread safe and data is destroyed when the process dies. In addition, there is no key rotation.

Note: If you want to test out the registration flow, you need to enable that feature by passing `--enable-registration=true` as well.

## Building

To build using the go binary on your host, use the `./build` script.

You can also use a copy of `go` hosted inside a Docker container if you prefix your command with `go-docker`, as in: `./go-docker ./build`

## Docker Build and Push

Once binaries are compiled you can build and push a dex image to quay.io. Before doing this step binaries must be built above using one of the build tools.

```console
export DOCKER_USER=<<your user>>
export DOCKER_PASSWORD=<<your password>>
./build-docker-push
```

By default the script pushes to `quay.io/coreos/dex`; if you want to push to a different repository, override the `DOCKER_REGISTRY` and `DOCKER_REPO` environment variables.

## Rebuild API from JSON schema

Go API bindings are generated from a JSON Discovery file.
To regenerate run:

```console
schema/generator
```

For updating generator dependencies see docs in: `schema/generator_import.go`.

## Running Tests

To run all tests (except functional) use the `./test` script;

If you want to test a single package only, use `PKG=<pkgname> ./test`

The functional tests require a database; create a database (eg. `createdb dex_func_test`) and then pass it as an environment variable to the functional test script, eg.  `DEX_TEST_DSN=postgres://localhost/dex_func_test?sslmode=disable ./test-functional`

To run these tests with Docker is a little trickier; you need to have a container running Postgres, and then you need to link that container to the container running your tests:


```console
# Run the Postgres docker container, which creates a db called "postgres"
docker run --name dex_postgres -d postgres

# The host name in the DSN is "postgres"; that works because that is what we
# will alias the link as, which causes Docker to modify /etc/hosts with a "postgres"
# entry.
export DEX_TEST_DSN=postgres://postgres@postgres/postgres?sslmode=disable

# Run the test container, linking it to the Postgres container.
DOCKER_LINKS=dex_postgres:postgres DOCKER_ENV=DEX_TEST_DSN ./go-docker ./test-functional docker

# Remove the container after the tests are run.
rm -f dex_postgres
```

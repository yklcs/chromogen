# panchro

Self-hosted image gallery.

Some key features:

- Static site generator and server modes
- EXIF metadata support
- Themable and extendable
- REST API
- Own your own images
- Minimal JavaScript, super light
- RSS/Atom feed

## Directory layout

- build: files for `panchro build`
- internal
- serve: files for `panchro serve`
- web

## Design

Images are ordered by file name (descending).

## Get started

Download the latest release, or build and install via the Go toolchain:

```shell
$ go install github.com/yklcs/panchro@latest
$ panchro serve
```

## Usage

Building/serving related config is done through CLI flags, and site related config is done through [a JSON file](panchro.example.json).

```shell
# Build site from images/
$ panchro build images

# Build site from images/ and add new.jpg to site
$ panchro build images
$ panchro build new.jpg

# Build site from images/ with config from config.json
$ panchro build -c=config.json images
```

### Auth

On startup, a random token is generated and output to stdout.
This token is used for API calls and the admin page.
Currently, tokens live as long as the process.

## API

Get all photos as JSON

```http
GET /photos
```

Get single photo as JSON

```http
GET /photos/{id}
```

Delete single photo

```http
DELETE /photos/{id}
Authorization: Bearer {token}
```

Upload photo (via form data)

```http
POST /photos
Content-Type: multipart/form-data
Authorization: Bearer {token}
```

## Theming

Theming is performed through [Go templates](https://pkg.go.dev/html/template) and static files.

Look at [web/](web/) for an example of the default theme. Mandatory template files are:

- index.tmpl
- photo.tmpl
- panchro.tmpl
- auth.tmpl

Theme-specific config should go in `"theme_config"` of the config file.

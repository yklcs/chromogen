# panchro

Self-hosted image gallery.

Some key features:

- Own your own images
- EXIF metadata support
- Themable and extendable
- Static site generator and server modes
- Minimal JS, progressively enhanced (used for thumbs + `/panchro` admin page only)

## Auth

On startup, a random token is generated and output to stdout. This token should be used for API calls in the form of:

```http
Authorization: Bearer tokengoeshere
```

Or, in the admin page.

Currently, tokens live as long as the process.

## Usage

Building/serving related config is done through CLI flags, and site related config is done through a JSON file.

```shell
# build static site (input from ./images, output to ./dist)
$ panchro build images

# serve (save photos and db to ./panchro, listen on :8000)
$ panchro serve
```

## API

```http
# Get all photos as JSON
GET /photos
```

```http
# Get single photo as JSON
GET /photos/{id}
```

```http
# Delete single photo
DELETE /photos/{id}
Authorization: Bearer {token}
```

```http
# Upload photo (via form data)
POST /photos
Content-Type multipart/form-data
Authorization: Bearer {token}
```

## Theming

Theming is performed through [Go templates](https://pkg.go.dev/html/template) and static files.

Look at [web/](web/) for an example of the default theme. Mandatory template files are:

- index.tmpl
- photo.tmpl
- panchro.tmpl

Theme-specific config should go in `"theme_config"` of the config file.

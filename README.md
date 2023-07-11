# panchro

Self-hosted image gallery.

Some key features:

- Own your own images
- EXIF metadata support
- Themable and extendable
- Static site generator and server modes

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

## Theming

Theming is performed through [Go templates](https://pkg.go.dev/html/template) and static files.

Look at [web/](web/) for an example of the default theme. Mandatory template files are:

- index.tmpl
- photo.tmpl
- panchro.tmpl

Theme-specific config should go in `"theme_config"` of the config file.

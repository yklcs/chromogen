# panchro

Self-hosted image gallery.

Why would anyone need this?

- Own your own images
- S3 integration
- EXIF/metadata support

## Usage

```
# static site
$ panchro build -d panchro -c config.json -i file://./images

# serve static site
$ panchro serve -d panchro -c config.json -i file://./images
```

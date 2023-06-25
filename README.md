# panchro

Self-hosted image gallery.

Why would anyone need this?

- Own your own images
- S3 integration
- EXIF/metadata support

## Usage examples

```shell
# build static site (input from ./images, output to ./dist)
$ panchro build images

# build static site (input from s3://input, output to s3://output)
$ panchro build -o s3://output s3://input

# serve cms (save photos to ./dist)
$ panchro serve

# serve cms (save photos to s3://photos)
$ panchro serve -storage s3://photos
```

## Backends

Local backend and S3 backends are supported. In both backends, original photos are left untouched.

On the local backend, copies of the original photos and their compressed versions are created in the output directory.
On the S3 backend, files are downloaded from S3 and copies of the original photos and their compressed versions are created in the output directory.

Uploads will not update the backend source directory or S3 bucket.

The `-urlprefix` flag can be used to use a custom URL prefix (CDN, S3, etc.) for images.
In this case, image URLs will use the URL prefix.

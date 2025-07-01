# cloudinary

A cloudinary driver for `facades.Storage()` of Goravel.

## Version

| goravel/cloudinary | goravel/framework |
|--------------------|-------------------|
| v1.4.*             | v1.16.*           |
| v1.3.*             | v1.15.*           |
| v1.2.*             | v1.14.*           |
| v1.1.*             | v1.13.*           |

## Install

Run the command below in your project to install the package automatically:

```
./artisan package:install github.com/goravel/cloudinary
```

Or check [the setup file](./setup/setup.go) to install the package manually.

## Testing

Run command below to run test:

```
CLOUDINARY_ACCESS_KEY_ID= CLOUDINARY_ACCESS_KEY_SECRET= CLOUDINARY_CLOUD= go test ./...
```

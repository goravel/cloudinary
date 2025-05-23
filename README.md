# cloudinary

A cloudinary driver for `facades.Storage()` of Goravel.

## Version

| goravel/cloudinary | goravel/framework |
|--------------------|-------------------|
| v1.3.*             | v1.15.*           |
| v1.2.*             | v1.14.*           |
| v1.1.*             | v1.13.*           |

## Install

1. Add package

```bash
go get github.com/goravel/cloudinary
```

2. Register service provider

```
// config/app.go
import "github.com/goravel/cloudinary"

"providers": []foundation.ServiceProvider{
    ...
    &cloudinary.ServiceProvider{},
}
```

3. Add cloudinary disk to `config/filesystems.go` file

```go
// config/filesystems.go
...
   
import (
      cloudinaryfacades "github.com/goravel/cloudinary/facades"
      "github.com/goravel/framework/contracts/filesystem"
)
   
"disks": map[string]filesystem.Disk{
      ...
      "cloudinary": map[string]any{
            "driver": "custom",
            "cloud":  config.Env("CLOUDINARY_CLOUD"),
            "key":    config.Env("CLOUDINARY_ACCESS_KEY_ID"), 
            "secret": config.Env("CLOUDINARY_ACCESS_KEY_SECRET"),
            "via": func()(filesystem.Driver, error) {
                  return cloudinaryfacades.Cloudinary("cloudinary") // The `cloudinary` value is the `disks` key
            },
      }
}
```
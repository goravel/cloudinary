# cloudinary
A cloudinary driver for `facades.Storage()` of Goravel.

## Version
| goravel/cloudinary | cloudinary | goravel/framework |
|--------------------|------------|-------------------|
| v1.0.0             | v2.3.0     | v1.12.0           |

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
            "github.com/goravel/framework/filesystem"
      )
   
      "disks": map[string]filesystem.Disk{
            ...
            "cloudinary": map[string]any{
                  "driver": "custom",
                  "cloud":  config.Env("CLOUDINARY_CLOUD"),
                  "key":    config.Env("CLOUDINARY_ACCESS_KEY_ID"), 
                  "secret": config.Env("CLOUDINARY_ACCESS_KEY_SECRET"),
                  "resource_types": map[string][]string{
                       "image": {"png"},
                       "video": {},
                       "raw":   {"txt", "pdf"},
                  },
                  "via": func()(filestystem.Disk, error) {
                        return cloudinaryfacades.Cloudinary("cloudinary"), nil // The `cloudinary` value is the `disks` key
                  },
            }
      }
      ```
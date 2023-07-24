# cloudinary
A cloudinary driver for facades.Storage of Goravel.

## Version
| goravel/cloudinary | cloudinary | goravel/framework |
| --- | --- |-------------------|
| 1.0.* | 1.13.* | v1.12.0           |

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
   import (
         cloudinaryFacades "github.com/goravel/cloudinary"
         "github.com/goravel/framework/filesystem"
   )
   
   "disks": map[string]filesystem.Disk{
         ...
         "cloudinary": map[string]any{
               "driver": "custom",
               "cloud": "your_cloud_name",
               "key": "your_api_key", 
               "secret": "your_api_secret",
               "via": func()(filestystem.Disk, error) {
                     return cloudinaryFacades.Cloudinary("cloudinary"), nil
               },
         }
   }
   ```
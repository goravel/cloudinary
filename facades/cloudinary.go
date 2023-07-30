package facades

import (
	"log"

	"github.com/goravel/framework/contracts/filesystem"

	"github.com/goravel/cloudinary"
)

func Cloudinary(disk string) filesystem.Driver {
	instance, err := cloudinary.App.MakeWith(cloudinary.Binding, map[string]any{"disk": disk})
	if err != nil {
		log.Println(err)
		return nil
	}

	return instance.(*cloudinary.Cloudinary)
}

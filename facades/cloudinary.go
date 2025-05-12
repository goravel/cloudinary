package facades

import (
	"github.com/goravel/framework/contracts/filesystem"

	"github.com/goravel/cloudinary"
)

func Cloudinary(disk string) (filesystem.Driver, error) {
	instance, err := cloudinary.App.MakeWith(cloudinary.Binding, map[string]any{"disk": disk})
	if err != nil {
		return nil, err
	}

	return instance.(*cloudinary.Cloudinary), nil
}

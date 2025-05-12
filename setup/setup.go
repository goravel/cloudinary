package main

import (
	"os"

	"github.com/goravel/framework/packages"
	"github.com/goravel/framework/packages/match"
	"github.com/goravel/framework/packages/modify"
	"github.com/goravel/framework/support/path"
)

var config = `map[string]any{
            "driver": "custom",
            "cloud":  config.Env("CLOUDINARY_CLOUD"),
            "key":    config.Env("CLOUDINARY_ACCESS_KEY_ID"), 
            "secret": config.Env("CLOUDINARY_ACCESS_KEY_SECRET"),
            "via": func()(filesystem.Driver, error) {
                  return cloudinaryfacades.Cloudinary("cloudinary") // The ` + "`cloudinary`" + ` value is the ` + "`disks`" + ` key
            },
      }`

func main() {
	packages.Setup(os.Args).
		Install(
			modify.GoFile(path.Config("app.go")).
				Find(match.Imports()).Modify(modify.AddImport(packages.GetModulePath())).
				Find(match.Providers()).Modify(modify.Register("&cloudinary.ServiceProvider{}")),
			modify.GoFile(path.Config("filesystems.go")).
				Find(match.Imports()).Modify(modify.AddImport("github.com/goravel/framework/contracts/filesystem"), modify.AddImport("github.com/goravel/cloudinary/facades", "cloudinaryfacades")).
				Find(match.Config("filesystems.disks")).Modify(modify.AddConfig("cloudinary", config)),
		).
		Uninstall(
			modify.GoFile(path.Config("app.go")).
				Find(match.Imports()).Modify(modify.RemoveImport(packages.GetModulePath())).
				Find(match.Providers()).Modify(modify.Unregister("&cloudinary.ServiceProvider{}")),
			modify.GoFile(path.Config("filesystems.go")).
				Find(match.Config("filesystems.disks")).Modify(modify.RemoveConfig("cloudinary")).
				Find(match.Imports()).Modify(modify.RemoveImport("github.com/goravel/framework/contracts/filesystem"), modify.RemoveImport("github.com/goravel/cloudinary/facades", "cloudinaryfacades")),
		).
		Execute()
}

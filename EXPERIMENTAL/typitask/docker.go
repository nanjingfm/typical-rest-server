package typitask

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"gopkg.in/yaml.v2"
)

// GenerateDockerCompose to generate docker-compose.yaml
func GenerateDockerCompose(ctx *typictx.ActionContext) (err error) {
	mainDocker := ctx.DockerCompose
	if mainDocker == nil {
		log.Info("No Docker Compose defined in Typical Context")
		return
	}
	for _, module := range ctx.Modules {
		moduleDocker := module.DockerCompose
		if moduleDocker == nil {
			continue
		}
		for name, service := range moduleDocker.Services {
			mainDocker.RegisterService(name, service)
		}
		for name, network := range moduleDocker.Networks {
			mainDocker.RegisterNetwork(name, network)
		}
		for name, volume := range moduleDocker.Volumes {
			mainDocker.RegisterVolume(name, volume)
		}
	}
	d1, _ := yaml.Marshal(mainDocker)
	return ioutil.WriteFile("docker-compose.yml", d1, 0644)
}

// DockerUp to create and start containers
func DockerUp(ctx *typictx.ActionContext) (err error) {
	// TODO:
	return
}

// DockerDown to stop and remove containers, networks, images, and volumes
func DockerDown(ctx *typictx.ActionContext) (err error) {
	// TODO:
	return
}
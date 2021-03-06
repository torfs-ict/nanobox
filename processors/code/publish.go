package code

import (
	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	container_generator "github.com/nanobox-io/nanobox/generators/containers"
	"github.com/nanobox-io/nanobox/generators/hooks/build"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/hookit"
)

// Publish ...
func Publish(envModel *models.Env, WarehouseConfig WarehouseConfig) error {
	display.OpenContext("Deploying app")
	defer display.CloseContext()

	// initialize the docker client
	// init docker client
	if err := provider.Init(); err != nil {
		return util.ErrorAppend(err, "failed to init docker client")
	}

	// pull the latest build image
	buildImage, err := pullBuildImage()
	if err != nil {
		return util.ErrorAppend(err, "failed to pull the build image")
	}

	display.StartTask("Starting docker container")

	// if a publish container was leftover from a previous publish, let's remove it
	docker.ContainerRemove(container_generator.PublishName())

	// start the container
	config := container_generator.PublishConfig(buildImage)
	container, err := docker.CreateContainer(config)
	if err != nil {
		lumber.Error("code:Build:docker.CreateContainer(%+v): %s", config, err.Error())
		display.ErrorTask()
		return util.ErrorAppend(err, "failed to start docker container")
	}
	// ensure we stop the container when we're done
	defer docker.ContainerRemove(container.ID)

	display.StopTask()

	lumber.Prefix("code:Publish")
	defer lumber.Prefix("")

	display.StartTask("Uploading")

	// run user hook
	// TODO: display output from hooks
	payload := build.UserPayload()
	if err != nil {
		lumber.Error("code:Publish:build.UserPayload()")
		return util.ErrorAppend(err, "unable to retrieve user payload")
	}
	if _, err := hookit.DebugExec(container.ID, "user", payload, "info"); err != nil {
		return err
	}

	buildWarehouseConfig := build.WarehouseConfig{
		BuildID:        WarehouseConfig.BuildID,
		WarehouseURL:   WarehouseConfig.WarehouseURL,
		WarehouseToken: WarehouseConfig.WarehouseToken,
		PreviousBuild:  WarehouseConfig.PreviousBuild,
	}

	payload = build.PublishPayload(envModel, buildWarehouseConfig)
	if err != nil {
		lumber.Error("code:Publish:build.UserPayload()")
		display.ErrorTask()
		return util.ErrorAppend(err, "unable to retrieve user payload")
	}
	if _, err := hookit.DebugExec(container.ID, "publish", payload, "info"); err != nil {
		return util.ErrorAppend(err, "failed to run publish hook")
	}

	display.StopTask()

	return nil
}

package client

import (
	"fmt"
	"strings"

	"golang.org/x/net/context"

	Cli "github.com/docker/docker/cli"
	flag "github.com/docker/docker/pkg/mflag"
	"github.com/docker/engine-api/types"
)

// CmdRm removes one or more containers.
//
// Usage: docker rm [OPTIONS] CONTAINER [CONTAINER...]
func (cli *DockerCli) CmdRm(args ...string) error {
	cmd := Cli.Subcmd("rm", []string{"CONTAINER [CONTAINER...]"}, Cli.DockerCommands["rm"].Description, true)
	v := cmd.Bool([]string{"v", "-volumes"}, false, "Remove the volumes associated with the container")
	link := cmd.Bool([]string{"l", "-link"}, false, "Remove the specified link")
	force := cmd.Bool([]string{"f", "-force"}, false, "Force the removal of a running container (uses SIGKILL)")
	cmd.Require(flag.Min, 1)

	cmd.ParseFlags(args, true)

	ctx := context.Background()

	var errs []string
	for _, name := range cmd.Args() {
		if name == "" {
			return fmt.Errorf("Container name cannot be empty")
		}
		name = strings.Trim(name, "/")

		if err := cli.RemoveContainer(ctx, name, *v, *link, *force); err != nil {
			errs = append(errs, err.Error())
		} else {
			fmt.Fprintf(cli.out, "%s\n", name)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "\n"))
	}
	return nil
}

// RemoveContainer removes a container
// TODO: this can be unexported again once all container commands are under
// api/client/container
func (cli *DockerCli) RemoveContainer(ctx context.Context, container string, removeVolumes, removeLinks, force bool) error {
	options := types.ContainerRemoveOptions{
		RemoveVolumes: removeVolumes,
		RemoveLinks:   removeLinks,
		Force:         force,
	}
	if err := cli.client.ContainerRemove(ctx, container, options); err != nil {
		return err
	}
	return nil
}

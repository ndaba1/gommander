package main

import "github.com/ndaba1/gommander"

func main() {

	// DISCLAIMER: NOT A FULL CLONE OF DOCKER FUNCTIONALITY, JUST A SIMPLE EXAMPLE
	// Run this example with --help to see how it is printed out

	app := gommander.App().
		Version("0.1.0").
		Help("A simple docker example").
		Name("docker")

	image := app.SubCommand("image").
		Alias("i").
		Help("Manage images").
		AddToGroup("Management Commands")

	// docker image subcommands
	image.SubCommand("ls").
		Help("List images").
		Flag("-a --all", "Show all images (default hides intermediate images)").
		Flag("--digests", "Show all images (default hides intermediate images)").
		Option("-f --filter filter", "Filter output based on conditions provided").
		Alias("l")

	image.SubCommand("build").
		Alias("b").
		Help("Build an image from a Dockerfile").
		Argument("<path>", "The path to the build context").
		Flag("--add-host list", "Add a custom host-to-IP mapping (host:ip)").
		Option(" -f --file <string>", "Name of the Dockerfile (Default is 'PATH/Dockerfile')")

	container := app.SubCommand("container").
		Alias("cont").
		Help("Manage containers").
		AddToGroup("Management Commands")

	// docker container subcommands
	container.SubCommand("prune").
		Help("Remove all stopped containers").
		Flag("-f --force", "Do not prompt for confirmation")

	container.SubCommand("start").
		Help("Start one or more stopped containers").
		Flag("-i --interactive", "Attach container's STDIN")

	container.SubCommand("create").
		Help("Create a new container").
		Flag("--add-host list", "Add a custom host-to-IP mapping (host:ip)")

	// Other app subcommands
	app.SubCommand("network").
		Alias("n").
		Help("Manage networks").
		AddToGroup("Management Commands")

	app.SubCommand("builder").
		Help("Manage builds").
		AddToGroup("Management Commands")

	app.SubCommand("compose").
		Help("Docker Compose (Docker Inc., v2.4.1)").
		AddToGroup("Management Commands")

	// A new set of command groups
	attach := app.
		SubCommand("attach").
		Help("Attach local standard input, output, and error streams to a running container")

	build := app.
		SubCommand("build").
		Help("Build an image from a Dockerfile")

	commit := app.
		SubCommand("commit").
		Help("Create a new image from a container's changes")

	cp := app.
		SubCommand("cp").
		Help("Copy files/folders between a container and the local filesystem")

	create := app.
		SubCommand("create").
		Help("Create a new container")

	diff := app.
		SubCommand("diff").
		Help("Inspect changes to files or directories on a container's filesystem")

	// An alternative way for creating subcommand groups
	app.SubCommandGroup("Commands", []*gommander.Command{attach, build, commit, cp, create, diff})

	// The real docker-cli has no color printing
	app.Set(gommander.DisableColor, true)

	app.Parse()
}

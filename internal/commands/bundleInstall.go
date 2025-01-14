package commands

import (
	"github.com/jlentink/aem/internal/aem"
	"github.com/jlentink/aem/internal/aem/bundle"
	"github.com/jlentink/aem/internal/cli/project"
	"github.com/jlentink/aem/internal/output"
	"github.com/spf13/cobra"
	"os"
)

type commandBundleInstall struct {
	verbose       bool
	instanceName  string
	instanceGroup string
	bundle        string
	startLevel    string
}

func (c *commandBundleInstall) setup() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bundle-install",
		Short:   "Install bundle",
		Aliases: []string{"binstall"},
		PreRun:  c.preRun,
		Run:     c.run,
	}
	cmd.Flags().StringVarP(&c.instanceName, "name", "n", aem.GetDefaultInstanceName(), "Instance to install bundle on")
	cmd.Flags().StringVarP(&c.instanceGroup, "group", "g", ``, "Instance group to install bundle on")
	cmd.Flags().StringVarP(&c.bundle, "bundle", "b", ``, "Instance group to install bundle on")
	cmd.Flags().StringVarP(&c.startLevel, "level", "l", "20", "Bundle start level")
	return cmd
}

func (c *commandBundleInstall) preRun(cmd *cobra.Command, args []string) {
	c.verbose, _ = cmd.Flags().GetBool("verbose")
	output.SetVerbose(c.verbose)
}

func (c *commandBundleInstall) run(cmd *cobra.Command, args []string) {
	_, is, errorString, err := getConfigAndInstanceOrGroupWithRoles(c.instanceName, c.instanceGroup, []string{aem.RoleAuthor, aem.RolePublisher})
	if err != nil {
		output.Printf(output.NORMAL, errorString, err.Error())
		os.Exit(EXIT_ERROR)
	}

	if !project.Exists(c.bundle) {
		output.Printf(output.NORMAL, "Could not find bundle at: %s", c.bundle)
		os.Exit(EXIT_ERROR)
	}

	for _, i := range is {
		err := bundle.Install(&i, c.bundle, c.startLevel)
		if err != nil {
			output.Printf(output.NORMAL, "Could not install bundle %s", err.Error())
			os.Exit(EXIT_ERROR)
		}
	}
}

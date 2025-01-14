package commands

import (
	"github.com/jlentink/aem/internal/aem"
	"github.com/jlentink/aem/internal/aem/bundle"
	"github.com/jlentink/aem/internal/output"
	"github.com/spf13/cobra"
	"os"
)

type commandBundleStart struct {
	verbose       bool
	instanceName  string
	instanceGroup string
	bundle        string
}

func (c *commandBundleStart) setup() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bundle-start",
		Short:   "Start bundle",
		Aliases: []string{"bstart"},
		PreRun:  c.preRun,
		Run:     c.run,
	}
	cmd.Flags().StringVarP(&c.instanceName, "name", "n", aem.GetDefaultInstanceName(), "Instance to install bundle on")
	cmd.Flags().StringVarP(&c.instanceGroup, "group", "g", ``, "Instance group to install bundle on")
	cmd.Flags().StringVarP(&c.bundle, "bundle", "b", ``, "Instance group to install bundle on")
	return cmd
}

func (c *commandBundleStart) preRun(cmd *cobra.Command, args []string) {
	c.verbose, _ = cmd.Flags().GetBool("verbose")
	output.SetVerbose(c.verbose)
}

func (c *commandBundleStart) run(cmd *cobra.Command, args []string) {
	_, is, errorString, err := getConfigAndInstanceOrGroupWithRoles(c.instanceName, c.instanceGroup, []string{aem.RoleAuthor, aem.RolePublisher})
	if err != nil {
		output.Printf(output.NORMAL, errorString, err.Error())
		os.Exit(EXIT_ERROR)
	}

	for idx, i := range is {
		if idx == 0 && c.bundle == "" {
			bundleObject, err := bundle.BundleSearch(&i, "Starting")
			if err != nil {
				output.Printf(output.NORMAL, "Could not list bundles: %s", err.Error())
				os.Exit(EXIT_ERROR)
			}
			c.bundle = bundleObject.SymbolicName
		}

		bndl, err := bundle.Get(&i, c.bundle)
		if err != nil {
			output.Printf(output.NORMAL, "Could not find bundle on: %s", i.Name)
			os.Exit(EXIT_ERROR)
		}

		b, err := bundle.Start(&i, bndl)
		if err != nil {
			output.Printf(output.NORMAL, "Could not start bundle %s", err.Error())
			os.Exit(EXIT_ERROR)
		}

		if b.StateRaw == 32 {
			output.Printf(output.NORMAL, "\U00002705 Bundle %s - %s\n", b.SymbolicName, bundle.BundleRawState[b.StateRaw])
			os.Exit(EXIT_NORMAL)
		}
		output.Printf(output.NORMAL, "\U0001F631 Bundle %s - %s\n", b.SymbolicName, bundle.BundleRawState[b.StateRaw])
		os.Exit(EXIT_ERROR)

	}
}

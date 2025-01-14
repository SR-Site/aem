package commands

import (
	"bufio"
	"fmt"
	"github.com/jlentink/aem/internal/aem/objects"
	"github.com/jlentink/aem/internal/output"
	"github.com/spf13/cobra"
	"os"
)

type commandPassword struct {
	verbose      bool
	instanceName string
}

func (c *commandPassword) setup() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "passwords",
		Aliases: []string{"password", "passwd"},
		Short:   "Set passwords into your keychain",
		PreRun:  c.preRun,
		Run:     c.run,
	}
	cmd.Flags().StringVarP(&c.instanceName, "name", "n", ``, "Update specific instance")
	return cmd
}

func (c *commandPassword) preRun(cmd *cobra.Command, args []string) {
	c.verbose, _ = cmd.Flags().GetBool("verbose")
	output.SetVerbose(c.verbose)
}

func (c *commandPassword) run(cmd *cobra.Command, args []string) {
	cnf, err := getConfig()
	if err != nil {
		output.Printf(output.NORMAL, "Could not get config file: %s", err.Error())
		os.Exit(EXIT_ERROR)
	}

	if len(c.instanceName) > 0 {
		_, i, errorString, err := getConfigAndInstance(c.instanceName)
		if err != nil {
			output.Printf(output.NORMAL, errorString, err.Error())
			os.Exit(EXIT_ERROR)
		}
		cnf.Instances = []objects.Instance{*i}
	}

	if !cnf.KeyRing {
		output.Printf(output.NORMAL, "keyring is disabled. use passwords from the aem.toml file.")
		os.Exit(EXIT_ERROR)
	}

	for _, i := range cnf.Instances {
		output.Printf(output.NORMAL, "\U0001F5A5 Instance: %s\n", i.Name)
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter password: ")
		pw, _ := reader.ReadString('\n')
		err := i.SetPassword(pw[:len(pw)-1])
		if err != nil {
			output.Printf(output.NORMAL, "Could not update password: %s", err.Error())
			os.Exit(EXIT_ERROR)
		}
	}
}

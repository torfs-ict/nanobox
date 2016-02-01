// package box ...
package box

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/vagrant"
)

var (
	BoxCmd = &cobra.Command{
		Use:   "box",
		Short: "Subcommands for managing the nanobox/boot2docker.box",
		Long:  ``,

		// ensure all dependencies are satisfied before running box commands
		PersistentPreRun: func(ccmd *cobra.Command, args []string) {

			// ensure vagrant exists
			if exists := vagrant.Exists(); !exists {
				fmt.Println("Missing dependency 'Vagrant'. Please download and install it to continue (https://www.vagrantup.com/).")
				os.Exit(1)
			}

			// ensure virtualbox exists
			if exists := util.VboxExists(); !exists {
				fmt.Println("Missing dependency 'Virtualbox'. Please download and install it to continue (https://www.virtualbox.org/wiki/Downloads).")
				os.Exit(1)
			}
		},
	}
)

//
func init() {
	BoxCmd.AddCommand(installCmd)
	BoxCmd.AddCommand(updateCmd)
}

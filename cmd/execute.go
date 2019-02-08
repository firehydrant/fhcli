package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	executeCmd.Flags().StringVarP(&environment, "environment", "e", "", "Environmnent UUID")
	executeCmd.Flags().StringVarP(&service, "service", "s", "", "Service UUID")
	executeCmd.Flags().StringVarP(&payload, "payload", "p", "", "Any additional metadata to include in the change execute")
	executeCmd.Flags().StringVarP(&identities, "identities", "i", "", "Change identities to allow for explicitly linking change executes, one or more k=v, comma separated")
	executeCmd.Flags().StringVarP(&revision, "revision", "r", "", "Git revision")
}

var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Execute command and submit change data",
	Long:  `Executes a command, monitors its duration and submits the change event to the FireHydrant API.`,
	RunE:  execute,
}

func execute(cmd *cobra.Command, args []string) error {
	setupClient()
	ce := NewChangeEvent()
	command := strings.Join(args, " ")

	ce.Summary = command
	ce.Environment = environment
	ce.Service = service
	ce.RawIdentities = identities

	fmt.Println(fmt.Sprintf("Executing command %s", command))
	cmdExec := exec.Command(command)

	// We don't actually want to capture this; pass it through to the calling shell
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	err := cmdExec.Run()
	ce.EndsAt = time.Now()

	// getting the exit code is a pita so i'm not doing it yet
	if err != nil {
		ce.Identities["error"] = err.Error()
		fmt.Println(fmt.Sprintf("Executed command, duration %dms, error: %s", 0, err))
	} else {
		fmt.Println(fmt.Sprintf("Executed command, duration %dms", 0))
	}

	id, err := ce.Submit()
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Created event %s", id))

	return nil
}

func init() {
	rootCmd.AddCommand(executeCmd)
}

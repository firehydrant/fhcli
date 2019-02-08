package cmd

import (
	"errors"
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
	Use:   "execute [flags] [command]",
	Short: "Execute command and submit change data",
	Long:  `Executes a command, monitors its duration and submits the change event to the FireHydrant API.`,
	RunE:  execute,
}

func execute(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("No command passed")
	}

	setupClient()

	ce := NewChangeEvent()
	ce.Summary = strings.Join(args, " ")
	ce.Environment = environment
	ce.Service = service
	ce.RawIdentities = identities

	command, flags := args[0], args[1:]

	fmt.Println(fmt.Sprintf("Executing command %s with args %s", command, strings.Join(flags, " ")))
	cmdExec := exec.Command(command, flags...)

	// We don't actually want to capture this; pass it through to the calling shell
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	ce.StartsAt = time.Now()
	err := cmdExec.Run()
	duration := time.Since(ce.StartsAt) / time.Millisecond
	ce.EndsAt = time.Now()

	// getting the exit code is a pita so i'm not doing it yet
	if err != nil {
		ce.Identities["error"] = err.Error()
		fmt.Println(fmt.Sprintf("Executed command, duration %dms, error: %s", duration, err))
	} else {
		fmt.Println(fmt.Sprintf("Executed command, duration %dms", duration))
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

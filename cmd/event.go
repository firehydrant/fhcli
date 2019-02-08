package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	eventCmd.Flags().StringVarP(&environment, "environment", "e", "", "Environmnent UUID")
	eventCmd.Flags().StringVarP(&service, "service", "s", "", "Service UUID")
	eventCmd.Flags().StringVarP(&payload, "payload", "p", "", "Any additional metadata to include in the change event")
	eventCmd.Flags().StringVarP(&identities, "identities", "i", "", "Change identities to allow for explicitly linking change events, one or more k=v, comma separated")
	eventCmd.Flags().StringVarP(&revision, "revision", "r", "", "Git revision")
}

var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "Submit change event data",
	Long:  `Submit change events and other metadata to the FireHydrant API.`,
	RunE:  event,
}

func event(cmd *cobra.Command, args []string) error {
	setupClient()
	ce := NewChangeEvent()

	ce.Summary = strings.Join(args, " ")
	ce.Environment = environment
	ce.Service = service
	ce.RawIdentities = identities

	id, err := ce.Submit()
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Created event %s", id))

	return nil
}

func init() {
	rootCmd.AddCommand(eventCmd)
}

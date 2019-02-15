package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/firehydrant/api-client-go/client/changes"
	"github.com/firehydrant/api-client-go/fhclient"
	"github.com/firehydrant/api-client-go/models"
	"github.com/go-openapi/strfmt"
	"github.com/urfave/cli"
)

func executeCmd(c *cli.Context) error {
	client, err := NewApiClient(c)
	if err != nil {
		return err
	}

	if len(c.Args()) < 1 {
		return errors.New("No command passed")
	}

	params := changes.NewPostV1ChangesEventsParams()
	command, flags := c.Args()[0], c.Args()[1:]

	identities := fhclient.ParseKV(c.String("identities"))

	params.V1ChangesEvents = &models.PostV1ChangesEvents{
		Environments: fhclient.ParamToList(c.String("environment")),
		Services:     fhclient.ParamToList(c.String("service")),
		StartsAt:     strfmt.DateTime(time.Now()),
		Summary:      &command,
		Labels:       fhclient.ParseKV(c.String("labels")),
	}

	fmt.Println(fmt.Sprintf("Executing command %s with args %s", command, strings.Join(flags, " ")))
	cmdExec := exec.Command(command, flags...)

	// We don't actually want to capture this; pass it through to the calling shell
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	start := time.Now()

	cErr := cmdExec.Run()

	duration := time.Since(start) / time.Millisecond
	params.V1ChangesEvents.EndsAt = strfmt.DateTime(time.Now())
	params.V1ChangesEvents.StartsAt = strfmt.DateTime(start)

	// getting the exit code is a pita so i'm not doing it yet
	if cErr != nil {
		identities["error"] = err.Error()
		fmt.Println(fmt.Sprintf("Executed command, duration %dms, error: %s", duration, cErr))
	} else {
		fmt.Println(fmt.Sprintf("Executed command, duration %dms", duration))
	}

	params.V1ChangesEvents.ChangeIdentities = fhclient.APIKVtoChangeIdentities(fhclient.MapToAPIKV(identities))
	resp, err := client.Client.Changes.PostV1ChangesEvents(params, client.Auth)

	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Created change event %s", resp.Payload.ID))

	if cErr != nil {
		return cErr
	}

	return nil
}
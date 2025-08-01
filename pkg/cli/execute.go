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
	"gopkg.in/urfave/cli.v1"
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
	summary := strings.Join(c.Args(), " ")

	identities := fhclient.ParseKV(c.String("identities"))

	params.PostV1ChangesEvents = &models.PostV1ChangesEvents{
		Environments: fhclient.ParamToList(c.String("environment")),
		Services:     fhclient.ParamToList(c.String("service")),
		StartsAt:     strfmt.DateTime(time.Now()),
		Summary:      &summary,
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
	params.PostV1ChangesEvents.EndsAt = strfmt.DateTime(time.Now())
	params.PostV1ChangesEvents.StartsAt = strfmt.DateTime(start)

	// getting the exit code is a pita so i'm not doing it yet
	if cErr != nil {
		identities["error"] = cErr.Error()
		fmt.Println(fmt.Sprintf("Executed command, duration %dms, error: %s", duration, cErr))
	} else {
		fmt.Println(fmt.Sprintf("Executed command, duration %dms", duration))
	}

	params.PostV1ChangesEvents.ChangeIdentities = fhclient.APIKVtoChangeIdentities(fhclient.MapToAPIKV(identities))
	resp, err := client.Client.Changes.PostV1ChangesEvents(params, client.Auth)

	if err != nil {
		return handleError(c, err)
	}

	fmt.Println(fmt.Sprintf("Created change event %s", resp.Payload.ID))

	if cErr != nil {
		return handleError(c, cErr)
	}

	return nil
}

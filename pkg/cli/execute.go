package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	changeevents "github.com/firehydrant/fhcli/pkg/change_events"
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

	ce := changeevents.NewChangeEvent()

	ce.Summary = strings.Join(c.Args(), " ")
	ce.Environment = c.String("environment")
	ce.Service = c.String("service")
	ce.RawIdentities = c.String("identities")
	ce.RawLabels = c.String("labels")

	command, flags := c.Args()[0], c.Args()[1:]

	fmt.Println(fmt.Sprintf("Executing command %s with args %s", command, strings.Join(flags, " ")))
	cmdExec := exec.Command(command, flags...)

	// We don't actually want to capture this; pass it through to the calling shell
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	ce.StartsAt = time.Now()
	err = cmdExec.Run()
	duration := time.Since(ce.StartsAt) / time.Millisecond
	ce.EndsAt = time.Now()

	// getting the exit code is a pita so i'm not doing it yet
	if err != nil {
		ce.Identities["error"] = err.Error()
		fmt.Println(fmt.Sprintf("Executed command, duration %dms, error: %s", duration, err))
	} else {
		fmt.Println(fmt.Sprintf("Executed command, duration %dms", duration))
	}

	id, err := ce.Submit(client)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Created event %s", id))

	return nil
}

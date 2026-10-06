package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/urfave/cli/v3"
)

func newExecuteCmd() *cli.Command {
	return &cli.Command{
		Name:   "execute",
		Usage:  "execute a command and submit metrics to the FireHydrant API",
		Action: executeCmd,
		Flags:  sharedFlags,
	}
}

func executeCmd(ctx context.Context, cmd *cli.Command) error {
	if cmd.NArg() < 1 {
		return handleError(cmd, errors.New("no command passed"))
	}
	// Without credentials the command still runs; only the change event is
	// skipped. Skipping a deploy because a login expired would be worse.
	client, clientErr := getClient(ctx)

	command := cmd.Args().First()
	flags := cmd.Args().Tail()
	summary := strings.Join(cmd.Args().Slice(), " ")

	identities := parseKV(cmd.String("identities"))

	if cmd.Bool("interactive") {
		fmt.Fprintf(os.Stderr, "Executing command %s with args %s\n", command, strings.Join(flags, " "))
	}

	cmdExec := exec.Command(command, flags...)
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	start := time.Now()

	cErr := cmdExec.Run()

	duration := time.Since(start) / time.Millisecond
	now := time.Now()

	if cErr != nil {
		identities["error"] = cErr.Error()
		if cmd.Bool("interactive") {
			fmt.Fprintf(os.Stderr, "Executed command, duration %dms, error: %s\n", duration, cErr)
		}
	} else if cmd.Bool("interactive") {
		fmt.Fprintf(os.Stderr, "Executed command, duration %dms\n", duration)
	}

	if clientErr != nil {
		if err := handleError(cmd, clientErr); err != nil {
			return err
		}
		if cErr != nil {
			return handleError(cmd, cErr)
		}
		return nil
	}

	changeIdentities := make([]components.ChangeIdentity, 0, len(identities))
	for k, v := range identities {
		changeIdentities = append(changeIdentities, components.ChangeIdentity{
			Type:  k,
			Value: v,
		})
	}

	req := components.CreateChangeEvent{
		Summary:          summary,
		Environments:     paramToList(cmd.String("environment")),
		Services:         paramToList(cmd.String("service")),
		StartsAt:         &start,
		EndsAt:           &now,
		Labels:           parseKV(cmd.String("labels")),
		ChangeIdentities: changeIdentities,
	}

	resp, err := client.Changes.CreateChangeEvent(ctx, req)
	if err != nil {
		return handleError(cmd, err)
	}

	if err := outputResult(cmd, resp); err != nil {
		return err
	}

	if cErr != nil {
		return handleError(cmd, cErr)
	}

	return nil
}

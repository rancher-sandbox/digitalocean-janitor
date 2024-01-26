package action

import (
	"context"
	"fmt"
	"strings"

	"github.com/digitalocean/godo"
)

type DoJanitorAction interface {
	Cleanup(ctx context.Context, input *Input) error
}

func New(commit bool) DoJanitorAction {
	return &action{
		commit: commit,
	}
}

type action struct {
	commit bool
}

func (a *action) Cleanup(ctx context.Context, input *Input) error {

	//NOTE: ordering matters here!
	cleanupFuncs := map[string]CleanupFunc{
		"droplets": a.cleanDroplets,
	}

	client := godo.NewFromToken(input.Token)

	regions, _, err := client.Regions.List(ctx, &godo.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed getting list of regions: %w", err)
	}

	regionsToClean := []string{}

	if input.AllRegions() {
		for _, region := range regions {
			regionsToClean = append(regionsToClean, region.Slug)
		}
	} else {
		regionsToClean = strings.Split(input.Regions, ",")
		//TODO: validate the ones passed in
	}

	for name, cleanupFunc := range cleanupFuncs {

		scope := &CleanupScope{
			TTL:       input.TTL,
			Client:    client,
			Commit:    input.Commit,
			Regions:   regionsToClean,
			IgnoreTag: input.IgnoreTag,
		}

		Log("Cleaning up %s", name)
		if err := cleanupFunc(ctx, scope); err != nil {
			return fmt.Errorf("failed running cleanup for %s: %w", name, err)
		}
	}

	return nil
}

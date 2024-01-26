package action

import (
	"context"
	"fmt"
	"time"

	"github.com/digitalocean/godo"
)

func (a *action) cleanDroplets(ctx context.Context, input *CleanupScope) error {
	LogDebug("Cleanup up droplets")

	list := []godo.Droplet{}
	opt := &godo.ListOptions{}

	for {
		droplets, resp, err := input.Client.Droplets.List(ctx, opt)
		if err != nil {
			return fmt.Errorf("failed listing droplets: %w", err)
		}

		list = append(list, droplets...)

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return fmt.Errorf("failed getting current droplets page: %w", err)
		}

		opt.Page = page + 1
	}

	if len(list) == 0 {
		Log("no droplets to delete")
		return nil
	}

	dropletsToDelete := []godo.Droplet{}
	for _, droplet := range list {
		createdAt, err := time.Parse(time.RFC3339, droplet.Created)
		if err != nil {
			return fmt.Errorf("failed to parse time of droplet: %w", err)
		}

		maxAge := createdAt.Add(input.TTL)

		if time.Now().Before(maxAge) {
			LogDebug("droplet %s has max age greater than now, skipping cleanup", droplet.Name)
			continue
		}

		if dropletHastag(&droplet, input.IgnoreTag) {
			LogDebug("droplet %s has ignore tag %s, skipping cleanup", droplet.Name, input.IgnoreTag)
			continue
		}

		if !dropletInRegion(&droplet, input.Regions) {
			LogDebug("droplet %s not in the specified regions, skipping cleanup", droplet.Name)
		}

		dropletsToDelete = append(dropletsToDelete, droplet)

	}

	if len(dropletsToDelete) == 0 {
		Log("no droplets to delete")
		return nil
	}

	for _, droplet := range dropletsToDelete {
		if !a.commit {
			Log("[dryrun] would delete droplet %s", droplet.Name)
			continue
		}

		if err := a.deleteDroplet(ctx, &droplet, input.Client); err != nil {
			LogError("failed to delete droplet %s: %s", droplet.Name, err)
		}
	}

	return nil
}

func (a *action) deleteDroplet(ctx context.Context, droplet *godo.Droplet, client *godo.Client) error {
	Log("deleting droplet %s", droplet.Name)

	if _, err := client.Droplets.Delete(ctx, droplet.ID); err != nil {
		return fmt.Errorf("failed to delete droplet %s(%d): %w", droplet.Name, droplet.ID, err)
	}

	return nil
}

func dropletInRegion(droplet *godo.Droplet, regions []string) bool {
	for _, region := range regions {
		if droplet.Region.Slug == region {
			return true
		}
	}

	return false
}

func dropletHastag(droplet *godo.Droplet, tagName string) bool {
	if tagName == "" {
		return false
	}

	if len(droplet.Tags) == 0 {
		return false
	}

	for _, tag := range droplet.Tags {
		if tag == tagName {
			return true
		}
	}

	return false
}

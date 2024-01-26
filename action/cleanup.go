package action

import (
	"context"
	"time"

	"github.com/digitalocean/godo"
)

type CleanupScope struct {
	Client    *godo.Client
	TTL       time.Duration
	Commit    bool
	Regions   []string
	IgnoreTag string
}

type CleanupFunc func(ctx context.Context, input *CleanupScope) error

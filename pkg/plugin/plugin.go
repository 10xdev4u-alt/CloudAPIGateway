package plugin

import (
	"context"
	"net/http"
)

type Plugin interface {
	Name() string
	OnRequest(ctx context.Context, r *http.Request) error
	OnResponse(ctx context.Context, res *http.Response) error
}

type Metadata struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

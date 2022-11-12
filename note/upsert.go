package note

import (
	"context"
	"fmt"
	"os"
)

// Upsert creates/updates a file under path and creates appropriate index entries.
func Upsert(ctx context.Context, dir string, ar Note, update bool) (*Note, error) {
	if err := ar.Sanitize(); err != nil {
		return nil, err
	}

	md, err := ar.Markdown()
	if err != nil {
		return nil, err
	}

	dir = cleanPath(dir, ar.Slug)

	if !update {
		if _, err := os.Stat(dir); err == nil {
			return nil, fmt.Errorf("file '%s' already exists, but this is not an update", dir)
		}
	}

	fh, err := os.OpenFile(dir, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	if _, err := fh.WriteString(md); err != nil {
		return nil, err
	}

	return &ar, nil
}

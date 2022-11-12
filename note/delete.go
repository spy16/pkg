package note

import (
	"context"
	"os"
	"path/filepath"
)

// Delete removes an note file and index by its unique identifier.
func Delete(ctx context.Context, dir, title string) error {
	dir = cleanPath(dir, title)
	return os.Remove(dir)
}

// DeleteAll provides a function to batch delete notes matching the given
// query. If the underlying store does not implemented batchDeleter, search
// will be used to get all notes and then Delete will be used.
func DeleteAll(ctx context.Context, dir string, query Query) error {
	notes, err := Search(ctx, dir, query)
	if err != nil {
		return err
	}

	if cancelled(ctx) {
		return ctx.Err()
	}

	for _, ar := range notes {
		if err := Delete(ctx, dir, ar.Title); err != nil {
			return err
		}
	}

	return nil
}

func cancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true

	default:
		return false
	}
}

func cleanPath(base string, p string) string {
	// TODO: Sanitize file name
	return filepath.Join(base, p+".md")
}

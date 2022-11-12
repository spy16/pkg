package note

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Search searches in the path for files that match the given query criteria and
// returns all the notes matching the query.
func Search(ctx context.Context, dir string, query Query) ([]Note, error) {
	var results []Note

	err := filepath.Walk(dir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			// walkErr or not a markdown file
			return walkErr
		}

		ar, err := readNoteFile(path, query)
		if err != nil {
			if err == ErrNotMatch {
				log.Printf("no match: %s", path)
				return nil
			}

			return err
		}

		results = append(results, *ar)
		return nil
	})

	return results, err
}

func readNoteFile(path string, query Query) (*Note, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ar := Note{}
	if err := ar.FromMarkdown(data); err != nil {
		return nil, err
	}

	if err := query.Match(ar); err != nil {
		return nil, err
	}

	return &ar, nil
}

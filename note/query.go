package note

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
)

// ErrNotMatch is returned when a query does not match an note value.
var ErrNotMatch = errors.New("not match")

// Query represents selection criteria/search query. Zero value of query should
// match any note.
type Query struct {
	// Slug should be absolute Slug of the target note. If Slug is set, all other
	// criteria will be ignored.
	Slug string `json:"slug,omitempty"`

	// Title can be a regular expression for selecting the note by title.
	// If left empty, it will match any title (equivalent to setting '.*')
	Title string `json:"title,omitempty"`

	// Tags can be provided to select notes with matching tags. If left nil
	// any tag will match.
	Tags []string `json:"tags,omitempty"`

	// Meta can be set to select notes by meta fields.
	Meta map[string]string `json:"meta,omitempty"`

	// Limit sets the max number of notes to be returned. Search stops when
	// the limit is reached.
	Limit int `json:"limit,omitempty"`
}

// Match returns true if the query matches the given note based on all the
// specified criteria in the query.
func (q Query) Match(ar Note) error {
	if id := strings.TrimSpace(q.Slug); id != "" && id == ar.Slug {
		return nil
	}

	isMatch, err := q.matchTitle(ar.Title)
	if err != nil {
		return err
	}

	isMatch = isMatch && q.matchTags(ar.Tags) && q.matchMeta(ar.Meta)
	if !isMatch {
		return ErrNotMatch
	}

	return nil
}

func (q Query) matchTags(to []string) bool {
	tagSet := map[string]struct{}{}
	for _, tagQ := range to {
		tagSet[strings.TrimSpace(tagQ)] = struct{}{}
	}

	for _, tag := range q.Tags {
		_, found := tagSet[tag]
		if !found {
			return false
		}
	}

	return true
}

func (q Query) matchMeta(to map[string]string) bool {
	if q.Meta == nil || len(q.Meta) == 0 {
		return true
	}

	q.Meta = cleanMeta(q.Meta)

	return reflect.DeepEqual(q.Meta, to)
}

func (q Query) matchTitle(title string) (bool, error) {
	if strings.TrimSpace(q.Title) == "" {
		return true, nil
	}

	pattern, err := regexp.Compile(q.Title)
	if err != nil {
		return false, err
	}

	if !pattern.MatchString(title) {
		return false, nil
	}

	return true, nil
}

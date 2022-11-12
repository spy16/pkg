package note

import (
	"bufio"
	"bytes"
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/unicode/norm"
	"gopkg.in/yaml.v2"
)

const startMeta = "---"

var (
	replacement = '_'

	// The "safe" set of characters.
	alphanum = &unicode.RangeTable{
		R16: []unicode.Range16{
			{0x0030, 0x0039, 1}, // 0-9
			{0x0041, 0x005A, 1}, // A-Z
			{0x0061, 0x007A, 1}, // a-z
		},
	}
	// Characters in these ranges will be ignored.
	nop = []*unicode.RangeTable{
		unicode.Mark,
		unicode.Sk, // Symbol - modifier
		unicode.Lm, // Letter - modifier
		unicode.Cc, // Other - control
		unicode.Cf, // Other - format
	}

	alphaNumeric = regexp.MustCompile("[^a-zA-Z0-9]+")
)

// Note represents a post or a link or an note etc. in dave.
type Note struct {
	Slug      string            `yaml:"slug" json:"slug"`
	Title     string            `yaml:"title" json:"title"`
	Tags      []string          `yaml:"tags" json:"tags"`
	Meta      map[string]string `yaml:"meta" json:"meta"`
	CreatedAt time.Time         `yaml:"created_at" json:"create_at"`
	UpdatedAt time.Time         `yaml:"updated_at" json:"updated_at"`
	Content   string            `yaml:"-" json:"-"`
}

// Sanitize performs basic validations, sanitizes and sets some defaults.
func (ar *Note) Sanitize() error {
	ar.Slug = strings.TrimSpace(ar.Slug)
	ar.Title = strings.TrimSpace(ar.Title)
	ar.Tags = cleanTags(ar.Tags)
	ar.Meta = cleanMeta(ar.Meta)

	if ar.Title == "" {
		return errors.New("title cannot be empty")
	}

	if ar.Slug == "" {
		ar.Slug = makeSlug(ar.Title)
	}

	if ar.CreatedAt.IsZero() {
		ar.CreatedAt = time.Now()
	}

	if ar.UpdatedAt.IsZero() {
		ar.UpdatedAt = time.Now()
	}

	return nil
}

// FromMarkdown reads data as markdown and populates the note.
func (ar *Note) FromMarkdown(data []byte) error {
	*ar = Note{}

	read, err := readMarkdownMeta(data, ar)
	if err != nil {
		return err
	}

	readMarkdown(data[read:], ar)
	ar.Content = string(data[read:])

	return nil
}

// Markdown converts the note into markdown format. Content will not
// be formatted and will be written as is.
func (ar *Note) Markdown() (string, error) {
	content := ar.Content
	ar.Content = ""

	meta, err := yaml.Marshal(ar)
	if err != nil {
		return "", err
	}

	var s strings.Builder

	s.WriteString(startMeta + "\n")
	_, _ = s.Write(meta)
	s.WriteString(startMeta + "\n\n")

	s.WriteString(content)
	return s.String(), nil
}

func readMarkdown(data []byte, ar *Note) {
	sc := bufio.NewScanner(bytes.NewReader(data))

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "# ") {
			ar.Title = line[2:]
			break
		} else if strings.HasPrefix(line, "## ") {
			ar.Title = line[3:]
			break
		}
	}
}

func readMarkdownMeta(data []byte, ar *Note) (int, error) {
	var metaYAML []byte
	sc := bufio.NewScanner(bytes.NewReader(data))

	inMeta := false
	count := 0

	for sc.Scan() {
		line := sc.Text()
		count += len(line) + 1

		if strings.TrimSpace(line) == "" {
			continue
		}

		if !inMeta {
			if line == startMeta {
				inMeta = true
			}

			continue
		} else if line == startMeta {
			break
		}

		metaYAML = append(metaYAML, []byte(line+"\n")...)
	}

	if err := yaml.Unmarshal(metaYAML, ar); err != nil {
		return 0, err
	}

	return count, nil
}

func cleanTags(tags []string) []string {
	var clean []string
	for _, t := range tags {
		t = strings.TrimSpace(t)

		if t != "" {
			clean = append(clean, t)
		}
	}

	return clean
}

func cleanMeta(meta map[string]string) map[string]string {
	if meta == nil {
		return nil
	}

	clean := map[string]string{}
	for k, v := range meta {
		k = strings.TrimSpace(k)
		if k != "" {
			clean[k] = v
		}
	}

	return clean
}

func makeSlug(s string) string {
	buf := make([]rune, 0, len(s))
	replace := false

	for _, r := range norm.NFKD.String(s) {
		switch {
		case unicode.In(r, alphanum):
			buf = append(buf, unicode.ToLower(r))
			replace = true
		case unicode.IsOneOf(nop, r):
			// skip
		case replace:
			buf = append(buf, replacement)
			replace = false
		}
	}

	// Strip trailing Replacement byte
	if i := len(buf) - 1; i >= 0 && buf[i] == replacement {
		buf = buf[:i]
	}

	return string(buf)
}

package pindb

import (
	"errors"
	"sort"
	"strings"
)

type Tags []Tag

func (t Tags) populate(tags ...string) Tags {
	for _, s := range tags {
		t.add(NewTag(s))
	}
	return t
}

func (t Tags) index(tag Tag) int {
	f := sort.SearchStrings(t.Strings(), tag.String())
	if f >= 0 && f < len(t) && t[f].String() == tag.String() {
		return f
	}
	return -1
}

func (t Tags) has(tag Tag) bool {
	if len(t) > 0 &&
		strings.TrimSpace(tag.String()) != "" {
		return t.index(tag) >= 0
	}

	return false
}

func (t *Tags) add(tags ...Tag) {
	for _, tag := range tags {
		if strings.TrimSpace(tag.String()) != "" && !t.has(tag) {
			*t = append(*t, tag)
		}
	}
	if !sort.IsSorted(t) {
		sort.Sort(t)
	}
}

func (t *Tags) remove(tags ...Tag) {
	for _, tag := range tags {
		if strings.TrimSpace(tag.String()) != "" && t.has(tag) {
			index := t.index(tag)
			*t = append(*t, (*t)[:index]...)
			*t = append(*t, (*t)[index+1:]...)
		}
	}
}

func (t Tags) Strings() []string {
	n := []string{}
	for _, i := range t {
		n = append(n, i.String())
	}
	return n
}

func (t Tags) Len() int {
	return len(t)
}

func (t Tags) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t Tags) Less(i, j int) bool {
	return t[i].String() < t[j].String()
}

func (t Tags) record() []byte {
	return []byte(strings.Join(t.Strings(), "\u2064"))
}

func (t Tags) parse(record []byte) Tags {
	n := Tags{}
	p := strings.Split(string(record), "\u2064")
	for _, i := range p {
		n.add(NewTag(i))
	}
	return n
}

func newTags() Tags {
	return Tags{}
}

type Tag string

func (t Tag) String() string {
	return string(t)
}

func (t Tag) Validate() (bool, error) {
	text := t.String()

	if strings.Contains(text, "\u2063") {
		return false, errors.New("tag cannot contain invisible separator (U+2063)")
	}

	if strings.Contains(text, "\u2064") {
		return false, errors.New("tag cannot contain invisible plus (U+2064)")
	}

	if strings.Contains(text, " ") {
		return false, errors.New("tag cannot contain spaces")
	}

	if strings.Contains(text, ",") {
		return false, errors.New("tag cannot contain commas")
	}

	return true, nil
}

func (t Tag) Is(tag Tag) bool {
	return t.String() == tag.String()
}

func NewTag(text string) Tag {
	return Tag(text)
}

func formatTag(text string) Tag {
	return Tag(strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ToLower(text),
			",",
			"٬"),
		" ",
		"-"))
}

package pindb

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/tmstn/pinboard"
)

type links map[uuid.UUID]*Link

func (l *links) get(uuid uuid.UUID) (*Link, error) {
	v, ok := (*l)[uuid]
	if ok {
		return v, nil
	}

	return nil, errors.New("link does not exist")
}

func (l *links) has(uuid uuid.UUID) bool {
	_, ok := (*l)[uuid]
	return ok
}

func (l *links) set(value *Link) {
	(*l)[value.uuid] = value
}

func (l *links) unset(link *Link) error {
	if !l.has(link.uuid) {
		return errors.New("link does not exist")
	}
	delete(*l, link.uuid)
	return nil
}

func (l *links) writeBytes() []byte {
	var f bytes.Buffer
	for _, v := range *l {
		fmt.Fprintf(&f, "%s\n", v.record())
	}
	return f.Bytes()
}

func (l *links) list() []*Link {
	links := []*Link{}
	for _, v := range *l {
		links = append(links, v)
	}
	return links
}

func (l *links) json() LinksJSON {
	j := LinksJSON{}
	for _, v := range l.list() {
		j = append(j, v.JSON())
	}
	return j
}

func newLinks() *links {
	return &links{}
}

type Link struct {
	bucket      *Bucket
	uuid        uuid.UUID
	title       string
	description []byte
	url         *url.URL
	group       Tag
	tags        Tags
	warnings    Warnings
}

func (l *Link) UUID() uuid.UUID {
	return l.uuid
}

func (l *Link) Title() string {
	return l.title
}

func (l *Link) Group() Tag {
	return l.group
}

func (l *Link) SetGroup(group Tag) (*Link, error) {
	if !strings.HasPrefix(group.String(), fmt.Sprintf("/pindb/bucket:\"%s\"/group:\"", l.bucket.uuid.String())) {
		group = Tag(fmt.Sprintf("/pindb/bucket:\"%s\"/group:\"%s\"", l.bucket.uuid.String(), group.String()))
	}

	l.tags.remove(l.group)
	l.group = group
	err := l.bucket.store.pb.Posts.Add(l.Options(true))
	if err != nil {
		return l, err
	}
	l.description = l.record()
	return l, nil
}

func (l *Link) UnsetGroup() (*Link, error) {
	l.tags.remove(l.group)
	l.group = NewTag("")
	err := l.bucket.store.pb.Posts.Add(l.Options(true))
	if err != nil {
		return l, err
	}
	l.description = l.record()
	return l, nil
}

func (l *Link) URL(detail bool) *url.URL {
	if detail {
		return l.url
	} else {
		r := *l.url
		r.Query().Del("pindbuuid")
		r.RawQuery = r.Query().Encode()
		return &r
	}
}

func (l *Link) Tags(detail bool) Tags {
	if detail {
		return l.tags
	} else {
		rgs := regexp.MustCompile(`^/pindb/store:\"[0-9a-f\-]+\"$`)
		rgb := regexp.MustCompile(`^/pindb/store:\"[0-9a-f\-]+\"/bucket:\"[0-9a-f\-]+\"$`)
		rgg := regexp.MustCompile(`^/pindb/bucket:\"([0-9a-f\-]+)\"/group:\"([0-9a-f\-]+)\"$`)

		s := newTags()
		for _, t := range l.tags {
			if rgs.MatchString(t.String()) {
				continue
			} else if rgb.MatchString(t.String()) {
				continue
			} else if rgg.MatchString(t.String()) {
				continue
			}

			s = append(s, t)
		}

		return s
	}
}

func (l *Link) Warnings() Warnings {
	return l.warnings
}

func (l *Link) Remove() error {
	err := l.bucket.store.user.pb.Posts.Delete(l.url.String())
	if err != nil {
		return err
	}
	l.bucket.links.unset(l)
	return nil
}

func (l *Link) Validate() bool {
	warnings := newWarnings()

	if string(l.description) != string(l.record()) {
		warnings = append(warnings, NewWarning(MismatchRecordWarning, nil))
	}
	q := l.url.Query().Get("pindbuuid")
	if strings.TrimSpace(q) == "" {
		warnings = append(warnings, NewWarning(NoUUIDWarning, nil))
	} else if q != l.uuid.String() {
		warnings = append(warnings, NewWarning(MismatchUUIDWarning, nil))
	}

	rgs := regexp.MustCompile(`^/pindb/store:\"[0-9a-f\-]+\"$`)
	rgb := regexp.MustCompile(`^/pindb/store:\"[0-9a-f\-]+\"/bucket:\"[0-9a-f\-]+\"$`)
	rgg := regexp.MustCompile(`^/pindb/bucket:\"([0-9a-f\-]+)\"/group:\"([0-9a-f\-]+)\"$`)
	gt := newTags()
	ugt := newTags()
	for _, t := range l.tags {
		st := NewTag(fmt.Sprintf("/pindb/store:\"%s\"", l.bucket.store.uuid.String()))
		bt := NewTag(fmt.Sprintf("/pindb/store:\"%s\"/bucket:\"%s\"", l.bucket.store.uuid.String(), l.bucket.uuid.String()))
		if rgs.MatchString(t.String()) && !t.Is(st) {
			warnings = append(warnings, NewWarning(UnrelatedPinDBStoreTagWarning, &t))
		}
		if rgb.MatchString(t.String()) && !t.Is(bt) {
			warnings = append(warnings, NewWarning(UnrelatedPinDBBucketTagWarning, &t))
		}
		if rgg.MatchString(t.String()) {
			buid := rgg.FindStringSubmatch(t.String())[1]
			if buid == l.bucket.uuid.String() {
				gt = append(gt, t)
			} else {
				ugt = append(ugt, t)
			}
		}
	}
	if len(gt) > 1 {
		for _, t := range gt {
			warnings = append(warnings, NewWarning(MultiplePinDBGroupTagWarning, &t))
		}
	}
	if len(ugt) > 01 {
		for _, t := range gt {
			warnings = append(warnings, NewWarning(UnrelatedPinDBGroupTagWarning, &t))
		}
	}
	l.warnings = warnings
	return len(l.warnings) == 0
}

func (l *Link) Options(replace bool) *pinboard.PostsAddOptions {
	l.tags.add(l.bucket.Tag())
	l.tags.add(l.bucket.store.Tag())
	l.description = l.record()
	if strings.TrimSpace(l.group.String()) != "" {
		l.tags.add(l.group)
	}

	opts := &pinboard.PostsAddOptions{}
	opts.Description = l.title
	opts.Extended = l.record()
	opts.Replace = replace
	opts.Shared = false
	opts.Toread = false
	opts.URL = l.url.String()
	opts.Extended = l.record()
	opts.Tags = l.tags.Strings()

	return opts
}

func (l *Link) Fix(warning Warning) (*Link, error) {
	return nil, nil
}

func (l *Link) JSON() LinkJSON {
	j := LinkJSON{}
	j.Description = string(l.description)
	j.Group = l.group.String()
	j.Tags = l.tags.Strings()
	j.Title = l.title
	j.UUID = l.uuid.String()
	j.Url = l.url.String()
	j.Warnings = l.warnings.json()
	return j
}

func (l *Link) record() []byte {
	return []byte(fmt.Sprintf(
		"L\u2063%s\u2063%s\u2063%s\u2063%s\u2063%s\u2063%s",
		l.uuid.String(),
		l.bucket.uuid.String(),
		l.url.String(),
		l.title,
		l.group.String(),
		l.tags.record()))
}

func newLink(bucket *Bucket, title string, url *url.URL, group Tag, tags ...Tag) (*Link, error) {
	if strings.Contains(title, "\u2063") {
		return nil, errors.New("title cannot contain invisible separator (U+2063)")
	}

	if strings.Contains(url.String(), "\u2063") {
		return nil, errors.New("url cannot contain invisible separator (U+2063)")
	}

	_, err := group.Validate()
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		_, err = tag.Validate()
		if err != nil {
			return nil, err
		}
	}

	l := &Link{
		uuid:     uuid.New(),
		title:    title,
		url:      url,
		group:    group,
		tags:     tags,
		warnings: newWarnings(),
		bucket:   bucket,
	}

	l.description = l.record()
	l.Validate()

	return l, nil
}

type LinksJSON []LinkJSON

type LinkJSON struct {
	UUID        string       `json:"uuid,omitempty"`
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	Url         string       `json:"url,omitempty"`
	Group       string       `json:"group,omitempty"`
	Tags        []string     `json:"tags,omitempty"`
	Warnings    WarningsJSON `json:"warnings,omitempty"`
}

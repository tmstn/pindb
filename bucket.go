package pindb

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tmstn/pinboard"
)

type buckets map[uuid.UUID]*Bucket

func (b *buckets) get(uuid uuid.UUID) (*Bucket, error) {
	v, ok := (*b)[uuid]
	if ok {
		return v, nil
	}

	return nil, errors.New("bucket does not exist")
}

func (b *buckets) has(uuid uuid.UUID) bool {
	_, ok := (*b)[uuid]
	return ok
}

func (b *buckets) set(bucket *Bucket) {
	(*b)[bucket.uuid] = bucket
}

func (b *buckets) unset(bucket *Bucket) error {
	if !b.has(bucket.uuid) {
		return errors.New("bucket does not exist")
	}
	delete(*b, bucket.uuid)
	return nil
}

func (b *buckets) writeBytes() []byte {
	var f bytes.Buffer
	for _, v := range *b {
		fmt.Fprintf(&f, "%s\n", v.record())
	}
	return f.Bytes()
}

func (b *buckets) list() []*Bucket {
	buckets := []*Bucket{}
	for _, v := range *b {
		buckets = append(buckets, v)
	}
	return buckets
}

func (b *buckets) json() BucketsJSON {
	j := BucketsJSON{}
	for _, v := range b.list() {
		j = append(j, v.JSON())
	}
	return j
}

func newBuckets() *buckets {
	return &buckets{}
}

type Bucket struct {
	refreshedAt *time.Time
	uuid        uuid.UUID
	name        string
	links       *links
	store       *Store
}

func (b *Bucket) Links() []*Link {
	return b.links.list()
}

func (b *Bucket) RefreshedAt() *time.Time {
	return b.refreshedAt
}

func (b *Bucket) UUID() uuid.UUID {
	return b.uuid
}

func (b *Bucket) Name() string {
	return b.name
}

func (b *Bucket) Tag() Tag {
	return NewTag(fmt.Sprintf("/pindb/store:\"%s\"/bucket:\"%s\"", b.store.uuid.String(), b.uuid.String()))
}

func (b *Bucket) Link(key uuid.UUID) (*Link, error) {
	return b.links.get(key)
}

func (b *Bucket) Has(key uuid.UUID) bool {
	return b.links.has(key)
}

func (b *Bucket) Add(title string, url *url.URL, group Tag, tags ...Tag) (*Link, error) {
	l, err := newLink(b, title, url, group, tags...)
	if err != nil {
		return nil, err
	}

	q := l.url.Query()
	q.Set("pindbuuid", l.uuid.String())
	l.url.RawQuery = q.Encode()

	err = b.store.user.pb.Posts.Add(l.Options(false))
	if err != nil {
		return nil, err
	}

	l.description = l.record()
	l.Validate()
	b.links.set(l)

	return l, nil
}

func (b *Bucket) Rename(name string) *Bucket {
	b.name = name
	return b
}

func (b *Bucket) Remove(removeLinks, removeTags bool) (*Bucket, error) {
	if removeLinks || removeTags {
		for _, l := range *b.links {
			var err error
			if removeLinks {
				err = b.store.pb.Posts.Delete(l.url.String())
			} else if removeTags {
				tag := fmt.Sprintf("/pindb/store:\"%s\"/bucket:\"%s\"", b.store.uuid.String(), b.uuid.String())
				err = b.store.pb.Tags.Delete(tag)
			}
			if err != nil {
				return b, err
			}
		}
	}
	b.store.buckets.unset(b)
	return nil, nil
}

func (b *Bucket) Refresh(force bool) (*Bucket, error) {
	if !force {
		refreshed, err := b.Updated()
		if err != nil {
			return b, err
		}

		if !refreshed {
			return b, err
		}
	}

	posts, err := b.store.pb.Posts.All(&pinboard.PostsAllOptions{
		Tag: []string{
			fmt.Sprintf(
				"/pindb/store:\"%s\"/bucket:\"%s\"",
				b.store.uuid.String(),
				b.uuid.String())},
	})

	if err != nil {
		return b, err
	}

	links := newLinks()
	for _, post := range posts {
		link, err := newLink(
			b,
			post.Description,
			post.Href,
			NewTag(""),
			newTags().populate(post.Tags...)...,
		)

		if err != nil {
			return b, err
		}

		uid := link.url.Query().Get("pindbuuid")
		if strings.TrimSpace(uid) == "" {
			q := link.url.Query()
			q.Set("pindbuuid", link.uuid.String())
			link.url.RawQuery = q.Encode()
		} else {
			puid, err := uuid.Parse(uid)
			if err == nil {
				link.uuid = puid
			}
		}

		rg := regexp.MustCompile(fmt.Sprintf(`^/pindb/bucket:\"%s\"/group:\"([0-9a-f\-]+)\"$`, link.bucket.uuid.String()))
		for _, t := range link.tags {
			if rg.MatchString(t.String()) {
				link.group = NewTag(rg.FindStringSubmatch(t.String())[1])
			}
		}

		link.description = link.record()
		link.Validate()
		links.set(link)
	}

	b.links = links
	ut := time.Now()
	b.refreshedAt = &ut
	return b, nil
}

func (b *Bucket) Updated() (bool, error) {
	if b.refreshedAt == nil {
		return true, nil
	}

	t, err := b.store.pb.Posts.Update()
	if err != nil {
		return false, err
	}

	return t.After(*b.refreshedAt), nil
}

func (b *Bucket) JSON() BucketJSON {
	j := BucketJSON{}
	j.Name = b.name
	j.UUID = b.uuid.String()
	if b.refreshedAt != nil {
		j.RefreshedAt = b.refreshedAt.Format(time.RFC3339)
	}
	j.Links = b.links.json()
	return j
}

func (b *Bucket) record() []byte {
	t := ""
	if b.refreshedAt != nil {
		t = b.refreshedAt.Format(time.RFC3339)
	}
	return []byte(fmt.Sprintf(
		"B\u2063%s\u2063%s\u2063%s\n%s",
		b.uuid.String(),
		t,
		b.name,
		b.links.writeBytes()))
}

func newBucket(store *Store, name string) (*Bucket, error) {
	if strings.Contains(name, "\u2063") {
		return nil, errors.New("name cannot contain invisible separator (U+2063)")
	}

	return &Bucket{
		uuid:  uuid.New(),
		name:  name,
		links: newLinks(),
		store: store,
	}, nil
}

type BucketsJSON []BucketJSON

type BucketJSON struct {
	RefreshedAt string    `json:"refreshed_at,omitempty"`
	UUID        string    `json:"uuid,omitempty"`
	Name        string    `json:"name,omitempty"`
	Links       LinksJSON `json:"links,omitempty"`
}

package pindb

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tmstn/pinboard"
)

type stores map[uuid.UUID]*Store

func (s *stores) get(uuid uuid.UUID) (*Store, error) {
	v, ok := (*s)[uuid]
	if ok {
		return v, nil
	}

	return nil, errors.New("store does not exist")
}

func (s *stores) has(uuid uuid.UUID) bool {
	_, ok := (*s)[uuid]
	return ok
}

func (s *stores) set(store *Store) {
	m := *s
	m[store.uuid] = store
	s = &m
}

func (s *stores) unset(store *Store) error {
	if !s.has(store.uuid) {
		return errors.New("store does not exist")
	}
	m := *s
	delete(m, store.uuid)
	s = &m
	return nil
}

func (s *stores) list() []*Store {
	stores := []*Store{}
	for _, v := range *s {
		stores = append(stores, v)
	}
	return stores
}

func (s *stores) json() StoresJSON {
	j := StoresJSON{}
	for _, v := range s.list() {
		j = append(j, v.JSON())
	}
	return j
}

func (s *stores) read(path string) (*Store, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, fmt.Errorf("%s is not a file", path)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return s.readBytes(b)
}

func (s *stores) readEncrypted(path, token string) (*Store, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, fmt.Errorf("%s is not a file", path)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	b, err = decrypt(token, b)
	if err != nil {
		return nil, err
	}

	return s.readBytes(b)
}

func (s *stores) readBase64(data string) (*Store, error) {
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return s.readBytes(b)
}

func (s *stores) readBytes(data []byte) (*Store, error) {
	f := string(data)
	if !strings.HasPrefix(f, "PINDBSTORE:\n") {
		return nil, errors.New("invalid pindb store file")
	}

	f = strings.TrimPrefix(f, "PINDBSTORE:\n")
	v := &Store{
		uuid:    uuid.New(),
		buckets: &buckets{},
	}

	for i, l := range strings.Split(f, "\n") {
		switch true {
		case strings.HasPrefix(l, "UA\u2063"):
			ts := strings.Split(l, "\u2063")[1]
			if strings.TrimSpace(ts) != "" {
				u, err := time.Parse(time.RFC3339, strings.Split(l, "\u2063")[1])
				if err != nil {
					return nil, fmt.Errorf("line %d: %s", i+2, err.Error())
				}
				v.refreshedAt = &u
			}
		case strings.HasPrefix(l, "SN\u2063"):
			v.name = strings.Split(l, "\u2063")[1]
		case strings.HasPrefix(l, "SU\u2063"):
			t := strings.Split(l, "\u2063")[1]
			u, err := newUser(t)
			if err != nil {
				return nil, fmt.Errorf("line %d: %s", i+2, err.Error())
			}

			u.pb = pinboard.New(t)
			err = u.Authenticate()
			if err != nil {
				return nil, fmt.Errorf("line %d: %s", i+2, err.Error())
			}

			v.pb = u.pb
			v.user = u
		case strings.HasPrefix(l, "SI\u2063"):
			uid, err := uuid.Parse(strings.Split(l, "\u2063")[1])
			if err != nil {
				return nil, fmt.Errorf("line %d: %s", i+2, err.Error())
			}

			v.uuid = uid
		case strings.HasPrefix(l, "B\u2063"):
			l = strings.TrimPrefix(l, "B\u2063")
			parts := strings.Split(l, "\u2063")
			if len(parts) != 3 {
				return nil, fmt.Errorf("line %d: invalid bucket record", i+2)
			}

			uid, err := uuid.Parse(parts[0])
			if err != nil {
				return nil, fmt.Errorf("line %d: %s", i+2, err.Error())
			}

			var t *time.Time
			if strings.TrimSpace(parts[1]) != "" {
				u, err := time.Parse(time.RFC3339, parts[1])
				if err != nil {
					return nil, fmt.Errorf("line %d: %s", i+2, err.Error())
				}
				t = &u
			}

			b, err := newBucket(v, parts[2])
			if err != nil {
				return nil, fmt.Errorf("line %d: %s", i+2, err.Error())
			}

			b.uuid = uid
			b.refreshedAt = t
			v.buckets.set(b)
		case strings.HasPrefix(l, "L\u2063"):
			l = strings.TrimPrefix(l, "L\u2063")
			parts := strings.Split(l, "\u2063")
			if len(parts) != 6 {
				return nil, fmt.Errorf("line %d: invalid link record", i+2)
			}

			uid, err := uuid.Parse(parts[0])
			if err != nil {
				return nil, fmt.Errorf("line %d: %s", i+2, err.Error())
			}

			buid, err := uuid.Parse(parts[1])
			if err != nil {
				return nil, fmt.Errorf("line %d: %s", i+2, err.Error())
			}

			b, err := v.Bucket(buid)
			if err != nil {
				return nil, fmt.Errorf("line %d: %s", i+2, err.Error())
			}

			u, err := url.Parse(parts[2])
			if err != nil {
				return nil, fmt.Errorf("line %d: %s", i+2, err.Error())
			}

			tags := newTags().parse([]byte(parts[5]))
			for _, t := range tags {
				_, err := t.Validate()
				if err != nil {
					return nil, fmt.Errorf("line %d: %s", i+2, err.Error())
				}

				tags.add(t)
			}

			n, err := newLink(b, parts[3], u, formatTag(parts[4]), tags...)
			if err != nil {
				return nil, fmt.Errorf("line %d: %s", i+2, err.Error())
			}

			n.uuid = uid
			n.description = n.record()
			n.Validate()
			b.links.set(n)
		}
	}

	s.set(v)
	return v, nil
}

func newStores() *stores {
	return &stores{}
}

type Store struct {
	refreshedAt *time.Time
	user        *user
	name        string
	uuid        uuid.UUID
	buckets     *buckets
	client      *Client
	pb          *pinboard.Client
}

func (s *Store) Buckets() []*Bucket {
	return s.buckets.list()
}

func (s *Store) RefreshedAt() *time.Time {
	return s.refreshedAt
}

func (s *Store) Name() string {
	return s.name
}

func (s *Store) UUID() uuid.UUID {
	return s.uuid
}

func (s *Store) Tag() Tag {
	return NewTag(fmt.Sprintf("/pindb/store:\"%s\"", s.uuid.String()))
}

func (s *Store) Bucket(uuid uuid.UUID) (*Bucket, error) {
	return s.buckets.get(uuid)
}

func (s *Store) Has(uuid uuid.UUID) bool {
	return s.buckets.has(uuid)
}

func (s *Store) Add(name string) (*Bucket, error) {
	b, err := newBucket(s, name)
	if err != nil {
		return nil, err
	}

	s.buckets.set(b)
	return b, nil
}

func (s *Store) Rename(name string) *Store {
	s.name = name
	return s
}

func (s *Store) Remove(removeLinks, removeTags bool) (*Store, error) {
	if removeLinks || removeTags {
		for _, b := range *s.buckets {
			_, err := b.Remove(removeLinks, removeLinks)
			if err != nil {
				return s, err
			}
		}
		if removeTags {
			tag := fmt.Sprintf("/pindb/store:\"%s\"", s.uuid.String())
			err := s.pb.Tags.Delete(tag)
			if err != nil {
				return s, err
			}
		}
	}
	s.client.stores.unset(s)
	return nil, nil
}

func (s *Store) Write(path string) error {
	info, err := os.Stat(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err == nil && info.IsDir() {
		return fmt.Errorf("%s is not a file", path)
	}

	err = os.WriteFile(
		path,
		s.WriteBytes(),
		0644,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) WriteEncrypted(path string, passphrase string) error {
	info, err := os.Stat(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err == nil && info.IsDir() {
		return fmt.Errorf("%s is not a file", path)
	}

	b := s.WriteBytes()
	enc, err := encrypt(passphrase, b)
	if err != nil {
		return err
	}

	err = os.WriteFile(
		path,
		enc,
		0644,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) WriteBytes() []byte {
	var b bytes.Buffer
	fmt.Fprint(&b, "PINDBSTORE:\n")
	if s.refreshedAt != nil {
		fmt.Fprintf(&b, "UA\u2063%s\n", s.refreshedAt.Format(time.RFC3339))
	}
	fmt.Fprintf(&b, "SN\u2063%s\n", s.name)
	fmt.Fprintf(&b, "SU\u2063%s\n", s.user.token)
	fmt.Fprintf(&b, "SI\u2063%s\n", s.uuid)
	b.Write(s.buckets.writeBytes())
	return b.Bytes()
}

func (s *Store) WriteBase64() string {
	b := s.WriteBytes()
	return base64.RawStdEncoding.EncodeToString(b)
}

func (s *Store) authenticate() error {
	err := s.user.Authenticate()
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Refresh(force bool) (*Store, error) {
	if !force {
		refreshed, err := s.Updated()
		if err != nil {
			return s, err
		}

		if !refreshed {
			return s, err
		}
	}

	posts, err := s.pb.Posts.All(&pinboard.PostsAllOptions{
		Tag: []string{
			fmt.Sprintf(
				"/pindb/store:\"%s\"",
				s.uuid.String())},
	})

	if err != nil {
		return s, err
	}

	for _, bucket := range *s.buckets {
		bucket.links = newLinks()
	}

	for _, post := range posts {
		buckets := []*Bucket{}
		tags := newTags().populate(post.Tags...)
		rg := regexp.MustCompile(`^/pindb/store:\"\"/bucket:\"([0-9a-f\-]+)\"$`)
		for _, tag := range tags {
			if rg.MatchString(tag.String()) {
				uid := rg.FindStringSubmatch(tag.String())[1]
				puid, err := uuid.Parse(uid)
				if err == nil {
					bucket, err := s.buckets.get(puid)
					if err == nil {
						buckets = append(buckets, bucket)
					}
				}
			}
		}
		for _, bucket := range buckets {
			link, err := newLink(
				bucket,
				post.Description,
				post.Href,
				NewTag(""),
				newTags().populate(post.Tags...)...,
			)

			if err != nil {
				return s, err
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
			bucket.links.set(link)
		}
	}

	ut := time.Now()
	s.refreshedAt = &ut
	return s, nil
}

func (s *Store) Updated() (bool, error) {
	if s.refreshedAt == nil {
		return true, nil
	}

	t, err := s.pb.Posts.Update()
	if err != nil {
		return false, err
	}

	return t.After(*s.refreshedAt), nil
}

func (s *Store) JSON() StoreJSON {
	j := StoreJSON{}
	if s.refreshedAt != nil {
		j.RefreshedAt = s.refreshedAt.Format(time.RFC3339)
	}
	j.User = s.user.JSON()
	j.Name = s.name
	j.UUID = s.uuid.String()
	j.Buckets = s.buckets.json()
	return j
}

func newStore(client *Client, token, name string) (*Store, error) {
	if strings.Contains(token, "\u2063") {
		return nil, errors.New("token cannot contain invisible separator (U+2063)")
	}

	if strings.Contains(name, "\u2063") {
		return nil, errors.New("name cannot contain invisible separator (U+2063)")
	}

	user, err := newUser(token)
	if err != nil {
		return nil, err
	}

	user.pb = pinboard.New(token)

	s := &Store{
		name:    name,
		user:    user,
		uuid:    uuid.New(),
		buckets: newBuckets(),
		client:  client,
	}

	err = s.authenticate()
	if err != nil {
		return nil, err
	}

	return s, nil
}

type StoresJSON []StoreJSON

type StoreJSON struct {
	RefreshedAt string      `json:"refreshed_at,omitempty"`
	User        UserJSON    `json:"user,omitempty"`
	Name        string      `json:"name,omitempty"`
	UUID        string      `json:"uuid,omitempty"`
	Buckets     BucketsJSON `json:"buckets,omitempty"`
}

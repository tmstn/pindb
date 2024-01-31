package pindb

import (
	"errors"
	"strings"

	"github.com/tmstn/pinboard"
)

type user struct {
	token    string
	username string
	key      string
	pb       *pinboard.Client
}

func (u *user) Token() string {
	return u.token
}

func (u *user) Username() string {
	return u.username
}

func (u *user) Key() string {
	return u.key
}

func (u *user) Authenticate() error {
	_, err := u.pb.User.Secret()
	if err != nil {
		return err
	}
	return nil
}

func (u *user) JSON() UserJSON {
	j := UserJSON{}
	j.Key = u.key
	j.Token = u.token
	j.Username = u.username
	return j
}

func newUser(token string) (*user, error) {
	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		return nil, errors.New("token format must be in the username:key format")
	}

	return &user{
		token:    token,
		username: parts[0],
		key:      parts[1],
	}, nil
}

type UserJSON struct {
	Token    string `json:"token,omitempty"`
	Username string `json:"username,omitempty"`
	Key      string `json:"key,omitempty"`
}

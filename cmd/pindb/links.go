package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/tmstn/pindb"
	"github.com/urfave/cli/v2"
)

func listLinks(cCtx *cli.Context) error {
	pdb := pindb.New()
	path := cCtx.String("path")
	passphrase := cCtx.String("passphrase")

	var store *pindb.Store
	var err error
	if strings.TrimSpace(passphrase) == "" {
		store, err = pdb.Read(path)
	} else {
		store, err = pdb.ReadEncrypted(path, passphrase)
	}

	if err != nil {
		return err
	}

	bid, err := uuid.Parse(cCtx.String("bucket"))
	if err != nil {
		return err
	}

	b, err := store.Bucket(bid)
	if err != nil {
		return err
	}

	printLinks(b.Links())

	return nil
}

func readLink(cCtx *cli.Context) error {
	pdb := pindb.New()
	path := cCtx.String("path")
	passphrase := cCtx.String("passphrase")

	var store *pindb.Store
	var err error
	if strings.TrimSpace(passphrase) == "" {
		store, err = pdb.Read(path)
	} else {
		store, err = pdb.ReadEncrypted(path, passphrase)
	}

	if err != nil {
		return err
	}

	bid, err := uuid.Parse(cCtx.String("bucket"))
	if err != nil {
		return err
	}

	b, err := store.Bucket(bid)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(cCtx.String("uuid"))
	if err != nil {
		return err
	}

	l, err := b.Link(id)
	if err != nil {
		return err
	}

	printLink(l)

	return nil
}

func addLink(cCtx *cli.Context) error {
	pdb := pindb.New()
	path := cCtx.String("path")
	passphrase := cCtx.String("passphrase")

	var store *pindb.Store
	var err error
	if strings.TrimSpace(passphrase) == "" {
		store, err = pdb.Read(path)
	} else {
		store, err = pdb.ReadEncrypted(path, passphrase)
	}

	if err != nil {
		return err
	}

	bid, err := uuid.Parse(cCtx.String("bucket"))
	if err != nil {
		return err
	}

	b, err := store.Bucket(bid)
	if err != nil {
		return err
	}

	title := cCtx.String("title")
	urls := cCtx.String("url")
	group := pindb.NewTag(cCtx.String("group"))
	tags := []pindb.Tag{}

	u, err := url.Parse(urls)
	if err != nil {
		return err
	}

	l, err := b.Add(title, u, group, tags...)
	if err != nil {
		return err
	}

	if strings.TrimSpace(passphrase) == "" {
		err = store.Write(path)
	} else {
		err = store.WriteEncrypted(path, passphrase)
	}

	if err != nil {
		return err
	}

	if cCtx.Bool("print") {
		printLink(l)
	}

	return nil
}

func setLinkGroup(cCtx *cli.Context) error {
	pdb := pindb.New()
	path := cCtx.String("path")
	passphrase := cCtx.String("passphrase")

	var store *pindb.Store
	var err error
	if strings.TrimSpace(passphrase) == "" {
		store, err = pdb.Read(path)
	} else {
		store, err = pdb.ReadEncrypted(path, passphrase)
	}

	if err != nil {
		return err
	}

	bid, err := uuid.Parse(cCtx.String("bucket"))
	if err != nil {
		return err
	}

	b, err := store.Bucket(bid)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(cCtx.String("uuid"))
	if err != nil {
		return err
	}

	l, err := b.Link(id)
	if err != nil {
		return err
	}

	group := pindb.NewTag(cCtx.String("group"))
	l, err = l.SetGroup(group)
	if err != nil {
		return err
	}

	if strings.TrimSpace(passphrase) == "" {
		err = store.Write(path)
	} else {
		err = store.WriteEncrypted(path, passphrase)
	}

	if err != nil {
		return err
	}

	if cCtx.Bool("print") {
		printLink(l)
	}

	return nil
}

func unsetLinkGroup(cCtx *cli.Context) error {
	pdb := pindb.New()
	path := cCtx.String("path")
	passphrase := cCtx.String("passphrase")

	var store *pindb.Store
	var err error
	if strings.TrimSpace(passphrase) == "" {
		store, err = pdb.Read(path)
	} else {
		store, err = pdb.ReadEncrypted(path, passphrase)
	}

	if err != nil {
		return err
	}

	bid, err := uuid.Parse(cCtx.String("bucket"))
	if err != nil {
		return err
	}

	b, err := store.Bucket(bid)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(cCtx.String("uuid"))
	if err != nil {
		return err
	}

	l, err := b.Link(id)
	if err != nil {
		return err
	}

	l, err = l.UnsetGroup()
	if err != nil {
		return err
	}

	if strings.TrimSpace(passphrase) == "" {
		err = store.Write(path)
	} else {
		err = store.WriteEncrypted(path, passphrase)
	}

	if err != nil {
		return err
	}

	if cCtx.Bool("print") {
		printLink(l)
	}

	return nil
}

func removeLink(cCtx *cli.Context) error {
	pdb := pindb.New()
	path := cCtx.String("path")
	passphrase := cCtx.String("passphrase")

	var store *pindb.Store
	var err error
	if strings.TrimSpace(passphrase) == "" {
		store, err = pdb.Read(path)
	} else {
		store, err = pdb.ReadEncrypted(path, passphrase)
	}

	if err != nil {
		return err
	}

	bid, err := uuid.Parse(cCtx.String("bucket"))
	if err != nil {
		return err
	}

	b, err := store.Bucket(bid)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(cCtx.String("uuid"))
	if err != nil {
		return err
	}

	l, err := b.Link(id)
	if err != nil {
		return err
	}

	err = l.Remove()
	if err != nil {
		return err
	}

	if strings.TrimSpace(passphrase) == "" {
		err = store.Write(path)
	} else {
		err = store.WriteEncrypted(path, passphrase)
	}

	if err != nil {
		return err
	}

	return nil
}

func fixLink(cCtx *cli.Context) error {
	pdb := pindb.New()
	path := cCtx.String("path")
	passphrase := cCtx.String("passphrase")

	var store *pindb.Store
	var err error
	if strings.TrimSpace(passphrase) == "" {
		store, err = pdb.Read(path)
	} else {
		store, err = pdb.ReadEncrypted(path, passphrase)
	}

	if err != nil {
		return err
	}

	bid, err := uuid.Parse(cCtx.String("bucket"))
	if err != nil {
		return err
	}

	b, err := store.Bucket(bid)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(cCtx.String("uuid"))
	if err != nil {
		return err
	}

	l, err := b.Link(id)
	if err != nil {
		return err
	}

	var warning pindb.WarningCategory
	var tag *pindb.Tag
	switch cCtx.String("warning") {
	case "mismatch_record":
		warning = pindb.MismatchRecordWarning
	case "no_uuid":
		warning = pindb.NoUUIDWarning
	case "mismatch_uuid":
		warning = pindb.MismatchUUIDWarning
	case "multiple_pindb_group_tag":
		warning = pindb.MultiplePinDBGroupTagWarning
	case "unrelated_pindb_group_tag":
		warning = pindb.UnrelatedPinDBGroupTagWarning
	case "unrelated_pindb_store_tag":
		warning = pindb.UnrelatedPinDBStoreTagWarning
	case "unrelated_pindb_bucket_tag":
		warning = pindb.UnrelatedPinDBBucketTagWarning
	default:
		return fmt.Errorf("unknown warning: %s", cCtx.String("warning"))
	}

	l, err = l.Fix(pindb.NewWarning(warning, tag))
	if err != nil {
		return err
	}

	if strings.TrimSpace(passphrase) == "" {
		err = store.Write(path)
	} else {
		err = store.WriteEncrypted(path, passphrase)
	}

	if err != nil {
		return err
	}

	if cCtx.Bool("print") {
		printLink(l)
	}

	return nil
}

func listJsonLinks(cCtx *cli.Context) error {
	pdb := pindb.New()
	path := cCtx.String("path")
	passphrase := cCtx.String("passphrase")

	var store *pindb.Store
	var err error
	if strings.TrimSpace(passphrase) == "" {
		store, err = pdb.Read(path)
	} else {
		store, err = pdb.ReadEncrypted(path, passphrase)
	}

	if err != nil {
		return err
	}

	bid, err := uuid.Parse(cCtx.String("bucket"))
	if err != nil {
		return err
	}

	b, err := store.Bucket(bid)
	if err != nil {
		return err
	}

	j := b.Links()
	d, err := json.MarshalIndent(j, "", " ")
	if err != nil {
		return nil
	}

	fmt.Println(string(d))

	return nil
}

func jsonLink(cCtx *cli.Context) error {
	pdb := pindb.New()
	path := cCtx.String("path")
	passphrase := cCtx.String("passphrase")

	var store *pindb.Store
	var err error
	if strings.TrimSpace(passphrase) == "" {
		store, err = pdb.Read(path)
	} else {
		store, err = pdb.ReadEncrypted(path, passphrase)
	}

	if err != nil {
		return err
	}

	bid, err := uuid.Parse(cCtx.String("bucket"))
	if err != nil {
		return err
	}

	b, err := store.Bucket(bid)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(cCtx.String("uuid"))
	if err != nil {
		return err
	}

	l, err := b.Link(id)
	if err != nil {
		return err
	}

	j := l.JSON()
	d, err := json.MarshalIndent(j, "", " ")
	if err != nil {
		return nil
	}

	fmt.Println(string(d))

	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/tmstn/pindb"
	"github.com/urfave/cli/v2"
)

func listBuckets(cCtx *cli.Context) error {
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

	printBuckets(store.Buckets(), cCtx.Bool("include-links"))

	return nil
}

func readBucket(cCtx *cli.Context) error {
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

	id, err := uuid.Parse(cCtx.String("uuid"))
	if err != nil {
		return err
	}

	b, err := store.Bucket(id)
	if err != nil {
		return err
	}

	printBucket(b, cCtx.Bool("include-links"))

	return nil
}

func addBucket(cCtx *cli.Context) error {
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

	name := cCtx.String("name")
	b, err := store.Add(name)
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
		printBucket(b, false)
	}

	return nil
}

func renameBucket(cCtx *cli.Context) error {
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

	id, err := uuid.Parse(cCtx.String("uuid"))
	if err != nil {
		return err
	}

	b, err := store.Bucket(id)
	if err != nil {
		return err
	}

	name := cCtx.String("name")
	b = b.Rename(name)

	if strings.TrimSpace(passphrase) == "" {
		err = store.Write(path)
	} else {
		err = store.WriteEncrypted(path, passphrase)
	}

	if err != nil {
		return err
	}

	if cCtx.Bool("print") {
		printBucket(b, false)
	}

	return nil
}

func removeBucket(cCtx *cli.Context) error {
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

	id, err := uuid.Parse(cCtx.String("uuid"))
	if err != nil {
		return err
	}

	b, err := store.Bucket(id)
	if err != nil {
		return err
	}

	remLinks := cCtx.Bool("remove-links")
	remTags := cCtx.Bool("remove-tags")
	_, err = b.Remove(remLinks, remTags)
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

func refreshBucket(cCtx *cli.Context) error {
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

	id, err := uuid.Parse(cCtx.String("uuid"))
	if err != nil {
		return err
	}

	b, err := store.Bucket(id)
	if err != nil {
		return err
	}

	force := cCtx.Bool("force")
	_, err = b.Refresh(force)
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

func listJsonBuckets(cCtx *cli.Context) error {
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

	j := store.JSON().Buckets
	if !cCtx.Bool("include-links") {
		for i := range j {
			j[i].Links = pindb.LinksJSON{}
		}
	}

	d, err := json.MarshalIndent(j, "", " ")
	if err != nil {
		return nil
	}

	fmt.Println(string(d))

	return nil
}

func jsonBucket(cCtx *cli.Context) error {
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

	id, err := uuid.Parse(cCtx.String("uuid"))
	if err != nil {
		return err
	}

	b, err := store.Bucket(id)
	if err != nil {
		return err
	}

	j := b.JSON()
	if !cCtx.Bool("include-links") {
		j.Links = pindb.LinksJSON{}
	}

	d, err := json.MarshalIndent(j, "", " ")
	if err != nil {
		return nil
	}

	fmt.Println(string(d))

	return nil
}

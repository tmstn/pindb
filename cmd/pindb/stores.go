package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/tmstn/pindb"
	"github.com/urfave/cli/v2"
)

func readStore(cCtx *cli.Context) error {
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

	printStore(store, cCtx.Bool("include-buckets"), cCtx.Bool("include-links"))

	return nil
}

func addStore(cCtx *cli.Context) error {
	pdb := pindb.New()
	path := cCtx.String("path")
	passphrase := cCtx.String("passphrase")
	token := cCtx.String("token")
	name := cCtx.String("name")

	store, err := pdb.Add(token, name)
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
		printStore(store, false, false)
	}

	return nil
}

func renameStore(cCtx *cli.Context) error {
	pdb := pindb.New()
	path := cCtx.String("path")
	passphrase := cCtx.String("passphrase")
	name := cCtx.String("name")

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

	store = store.Rename(name)
	if strings.TrimSpace(passphrase) == "" {
		err = store.Write(path)
	} else {
		err = store.WriteEncrypted(path, passphrase)
	}

	if err != nil {
		return err
	}

	if cCtx.Bool("print") {
		printStore(store, false, false)
	}

	return nil
}

func removeStore(cCtx *cli.Context) error {
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

	remLinks := cCtx.Bool("remove-links")
	remTags := cCtx.Bool("remove-tags")
	_, err = store.Remove(remLinks, remTags)
	if err != nil {
		return err
	}

	err = os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

func refreshStore(cCtx *cli.Context) error {
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

	force := cCtx.Bool("force")
	store, err = store.Refresh(force)
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

func jsonStore(cCtx *cli.Context) error {
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

	j := store.JSON()
	if !cCtx.Bool("include-buckets") {
		j.Buckets = pindb.BucketsJSON{}
	} else if !cCtx.Bool("include-links") {
		for i := range j.Buckets {
			j.Buckets[i].Links = pindb.LinksJSON{}
		}
	}

	d, err := json.MarshalIndent(j, "", " ")
	if err != nil {
		return nil
	}

	fmt.Println(string(d))

	return nil
}

package main

import "github.com/urfave/cli/v2"

func createApp() *cli.App {
	app := &cli.App{
		Name:  "pindb",
		Usage: "pinboard link database cli",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Usage:   "a path to a store db file",
				Aliases: []string{"p", "pt"},
			},
			&cli.StringFlag{
				Name:    "passphrase",
				Usage:   "a passphrase to encrypt/decrypt with",
				Aliases: []string{"pp", "pass"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "stores",
				Aliases: []string{"s", "st"},
				Usage:   "manage stores",
				Subcommands: []*cli.Command{
					{
						Name:   "read",
						Usage:  "read a store",
						Action: readStore,
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "include-buckets",
								Usage:   "include buckets in the output",
								Aliases: []string{"ib", "incb"},
							},
							&cli.BoolFlag{
								Name:    "include-links",
								Usage:   "include links in the output",
								Aliases: []string{"il", "incl"},
							},
						},
					},
					{
						Name:   "add",
						Usage:  "add a store",
						Action: addStore,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "token",
								Usage:    "a pinboard api token",
								Aliases:  []string{"t", "tkn"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "name",
								Usage:    "the name of the store",
								Aliases:  []string{"n", "nm"},
								Required: true,
							},
							&cli.BoolFlag{
								Name:    "print",
								Usage:   "print result of the operation",
								Aliases: []string{"p", "pr"},
							},
						},
					},
					{
						Name:   "rename",
						Usage:  "rename a store",
						Action: renameStore,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Usage:    "the new name for the store",
								Aliases:  []string{"n", "nm"},
								Required: true,
							},
							&cli.BoolFlag{
								Name:    "print",
								Usage:   "print result of the operation",
								Aliases: []string{"p", "pr"},
							},
						},
					},
					{
						Name:   "remove",
						Usage:  "remove a store",
						Action: removeStore,
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "remove-links",
								Usage:   "remove the links from pinboard",
								Aliases: []string{"rl", "reml"},
							},
							&cli.BoolFlag{
								Name:    "remove-tags",
								Usage:   "remove the tags from pinboard",
								Aliases: []string{"rt", "remt"},
							},
						},
					},
					{
						Name:   "refresh",
						Usage:  "refresh store from pinboard",
						Action: refreshStore,
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "force",
								Usage:   "force the refresh",
								Aliases: []string{"f", "frc"},
							},
						},
					},
					{
						Name:   "json",
						Usage:  "show a store as json",
						Action: jsonStore,
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "include-buckets",
								Usage:   "include buckets in the output",
								Aliases: []string{"ib", "incb"},
							},
							&cli.BoolFlag{
								Name:    "include-links",
								Usage:   "include links in the output",
								Aliases: []string{"il", "incl"},
							},
						},
					},
				},
			},
			{
				Name:    "buckets",
				Aliases: []string{"b", "bck"},
				Usage:   "manage buckets",
				Subcommands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "list all buckets",
						Action: listBuckets,
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "include-links",
								Usage:   "include links in the output",
								Aliases: []string{"il", "incl"},
							},
						},
					},
					{
						Name:   "read",
						Usage:  "show a single bucket",
						Action: readBucket,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "uuid",
								Usage:    "the uuid of the bucket",
								Aliases:  []string{"u", "uid"},
								Required: true,
							},
							&cli.BoolFlag{
								Name:    "include-links",
								Usage:   "include links in the output",
								Aliases: []string{"il", "incl"},
							},
						},
					},
					{
						Name:   "add",
						Usage:  "add a bucket",
						Action: addBucket,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Usage:    "the name of the bucket",
								Aliases:  []string{"n", "nm"},
								Required: true,
							},
							&cli.BoolFlag{
								Name:    "print",
								Usage:   "print result of the operation",
								Aliases: []string{"p", "pr"},
							},
						},
					},
					{
						Name:   "rename",
						Usage:  "rename a bucket",
						Action: renameBucket,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "uuid",
								Usage:    "the uuid of the bucket",
								Aliases:  []string{"u", "uid"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "name",
								Usage:    "the name of the bucket",
								Aliases:  []string{"n", "nm"},
								Required: true,
							},
							&cli.BoolFlag{
								Name:    "print",
								Usage:   "print result of the operation",
								Aliases: []string{"p", "pr"},
							},
						},
					},
					{
						Name:   "remove",
						Usage:  "remove a bucket",
						Action: removeBucket,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "uuid",
								Usage:    "the uuid of the bucket",
								Aliases:  []string{"u", "uid"},
								Required: true,
							},
							&cli.BoolFlag{
								Name:    "remove-links",
								Usage:   "remove the links from pinboard",
								Aliases: []string{"rl", "reml"},
							},
							&cli.BoolFlag{
								Name:    "remove-tags",
								Usage:   "remove the tags from pinboard",
								Aliases: []string{"rt", "remt"},
							},
						},
					},
					{
						Name:   "refresh",
						Usage:  "refresh a bucket from pinboard",
						Action: refreshBucket,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "uuid",
								Usage:    "the uuid of the bucket",
								Aliases:  []string{"u", "uid"},
								Required: true,
							},
							&cli.BoolFlag{
								Name:    "force",
								Usage:   "force the refresh",
								Aliases: []string{"f", "frc"},
							},
						},
					},
					{
						Name:   "listjson",
						Usage:  "list all buckets as json",
						Action: listJsonBuckets,
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "include-links",
								Usage:   "include links in the output",
								Aliases: []string{"il", "incl"},
							},
						},
					},
					{
						Name:   "json",
						Usage:  "show a bucket as json",
						Action: jsonBucket,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "uuid",
								Usage:    "the uuid of the bucket",
								Aliases:  []string{"u", "uid"},
								Required: true,
							},
							&cli.BoolFlag{
								Name:    "include-links",
								Usage:   "include links in the output",
								Aliases: []string{"il", "incl"},
							},
						},
					},
				},
			},
			{
				Name:    "links",
				Aliases: []string{"l", "link"},
				Usage:   "manage links",
				Subcommands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "list all links",
						Action: listLinks,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "bucket",
								Usage:    "the uuid of the bucket",
								Aliases:  []string{"b", "bck"},
								Required: true,
							},
						},
					},
					{
						Name:   "read",
						Usage:  "show a single link",
						Action: readLink,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "bucket",
								Usage:    "the uuid of the bucket",
								Aliases:  []string{"b", "bck"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "uuid",
								Usage:    "the uuid of the link",
								Aliases:  []string{"u", "uid"},
								Required: true,
							},
						},
					},
					{
						Name:   "add",
						Usage:  "add a link",
						Action: addLink,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "bucket",
								Usage:    "the uuid of the bucket",
								Aliases:  []string{"b", "bck"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "title",
								Usage:    "the title of the link",
								Aliases:  []string{"t", "ttl"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "url",
								Usage:    "the url of the link",
								Aliases:  []string{"u", "ur"},
								Required: true,
							},
							&cli.StringFlag{
								Name:    "group",
								Usage:   "the group of the link",
								Aliases: []string{"g", "grp"},
							},
							// &cli.StringSliceFlag{
							// 	Name:    "tags",
							// 	Usage:   "the tags of the link",
							// 	Aliases: []string{"tg", "tgs"},
							// },
							&cli.BoolFlag{
								Name:    "print",
								Usage:   "print result of the operation",
								Aliases: []string{"p", "pr"},
							},
						},
					},
					{
						Name:   "setgroup",
						Usage:  "set the group for a link",
						Action: setLinkGroup,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "bucket",
								Usage:    "the uuid of the bucket",
								Aliases:  []string{"b", "bck"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "uuid",
								Usage:    "the uuid of the link",
								Aliases:  []string{"u", "uid"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "group",
								Usage:    "the group of the link",
								Aliases:  []string{"g", "grp"},
								Required: true,
							},
							&cli.BoolFlag{
								Name:    "print",
								Usage:   "print result of the operation",
								Aliases: []string{"p", "pr"},
							},
						},
					},
					{
						Name:   "unsetgroup",
						Usage:  "unset the group for a link",
						Action: unsetLinkGroup,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "bucket",
								Usage:    "the uuid of the bucket",
								Aliases:  []string{"b", "bck"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "uuid",
								Usage:    "the uuid of the link",
								Aliases:  []string{"u", "uid"},
								Required: true,
							},
							&cli.BoolFlag{
								Name:    "print",
								Usage:   "print result of the operation",
								Aliases: []string{"p", "pr"},
							},
						},
					},
					{
						Name:   "remove",
						Usage:  "remove a link",
						Action: removeLink,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "bucket",
								Usage:    "the uuid of the bucket",
								Aliases:  []string{"b", "bck"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "uuid",
								Usage:    "the uuid of the link",
								Aliases:  []string{"u", "uid"},
								Required: true,
							},
						},
					},
					{
						Name:   "fix",
						Usage:  "fix a link warning",
						Action: fixLink,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "bucket",
								Usage:    "the uuid of the bucket",
								Aliases:  []string{"b", "bck"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "uuid",
								Usage:    "the uuid of the link",
								Aliases:  []string{"u", "uid"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "warning",
								Usage:    "the warning to fix",
								Aliases:  []string{"w", "war"},
								Required: true,
							},
							&cli.StringFlag{
								Name:    "tag",
								Usage:   "the tag the warning applies to",
								Aliases: []string{"t", "tg"},
							},
							&cli.BoolFlag{
								Name:    "print",
								Usage:   "print result of the operation",
								Aliases: []string{"p", "pr"},
							},
						},
					},
					{
						Name:   "listjson",
						Usage:  "list all links as json",
						Action: listJsonLinks,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "bucket",
								Usage:    "the uuid of the bucket",
								Aliases:  []string{"b", "bck"},
								Required: true,
							},
						},
					},
					{
						Name:   "json",
						Usage:  "show a link as json",
						Action: jsonLink,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "bucket",
								Usage:    "the uuid of the bucket",
								Aliases:  []string{"b", "bck"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "uuid",
								Usage:    "the uuid of the link",
								Aliases:  []string{"u", "uid"},
								Required: true,
							},
						},
					},
				},
			},
		},
	}

	return app
}

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/tmstn/pindb"
)

func printStore(s *pindb.Store, incBuckets, incLinks bool) {
	fmt.Println("========STORE:========")
	fmt.Printf("Name: %s\n", s.Name())
	fmt.Printf("UUID: %s\n", s.UUID())
	t := "Not Refreshed"
	if s.RefreshedAt() != nil {
		t = s.RefreshedAt().Format(time.RFC3339)
	}
	fmt.Printf("Refreshed At: %s\n", t)
	fmt.Printf("Tag: %s\n", s.Tag())

	if incBuckets {
		printBuckets(s.Buckets(), incLinks)
	}
}

func printBuckets(b []*pindb.Bucket, incLinks bool) {
	for _, i := range b {
		printBucket(i, incLinks)
	}
}

func printBucket(b *pindb.Bucket, incLinks bool) {
	fmt.Println("--------BUCKET:-------")
	fmt.Printf("Name: %s\n", b.Name())
	fmt.Printf("UUID: %s\n", b.UUID())
	t := "Not Refreshed"
	if b.RefreshedAt() != nil {
		t = b.RefreshedAt().Format(time.RFC3339)
	}
	fmt.Printf("Refreshed At: %s\n", t)
	fmt.Printf("Tag: %s\n", b.Tag())

	if incLinks {
		printLinks(b.Links())
	}
}

func printLinks(l []*pindb.Link) {
	for _, i := range l {
		printLink(i)
	}
}

func printLink(l *pindb.Link) {
	fmt.Println("++++++++LINK:++++++++")
	fmt.Printf("UUID: %s\n", l.UUID())
	fmt.Printf("Title: %s\n", l.Title())
	fmt.Printf("URL: %s\n", l.URL(false).String())
	fmt.Printf("Group: %s\n", l.Group())
	fmt.Printf("Tags: %s\n", strings.Join(l.Tags(false).Strings(), ", "))
	printWarnings(l.Warnings())
}

func printWarnings(w pindb.Warnings) {
	if len(w) > 0 {
		fmt.Println("......WARNINGS......")
		for i, m := range w {
			fmt.Printf("%d: %s\n", i+1, m.String())
		}
	}
}

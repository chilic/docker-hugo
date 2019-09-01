package main

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"sort"
	"time"

	"github.com/hashicorp/go-version"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

const tpl = `FROM alpine:latest

MAINTAINER IChe <me@iche.eu>

RUN apk add --no-cache \
	curl \
	git \
	openssh-client \
	rsync

RUN mkdir -p /usr/local/src \
	&& cd /usr/local/src \
	&& curl -L https://github.com/gohugoio/hugo/releases/download/v{{ .Version }}/hugo_{{ .Version }}_linux-64bit.tar.gz | tar -xz \
	&& mv hugo /usr/local/bin/hugo \
	&& addgroup -Sg 1000 hugo \
	&& adduser -SG hugo -u 1000 -h /src hugo

WORKDIR /src

EXPOSE 1313

ENTRYPOINT ["/usr/local/bin/hugo"]
CMD [ "--help" ]
`

const gitURL = "https://github.com/gohugoio/hugo.git"

func fetchTags(repo string) ([]string, error) {
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{repo},
	})

	// We can then use every Remote functions to retrieve wanted informations
	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Filters the references list and only keeps tags
	var tags []string
	for _, ref := range refs {
		if ref.Name().IsTag() {
			tags = append(tags, ref.Name().Short())
		}
	}

	if len(tags) == 0 {
		return nil, errors.New("0 tags found")
	}

	return tags, nil
}

func getLastVersion(tags []string) string {
	versions := make([]*version.Version, len(tags))
	for i, raw := range tags {
		v, _ := version.NewVersion(raw)
		versions[i] = v
	}
	// After this, the versions are properly sorted
	sort.Sort(version.Collection(versions))

	return versions[len(versions)-1].String()
}

func commitLocal(version string) {
	r, _ := git.PlainOpen("./")
	w, _ := r.Worktree()
	status, _ := w.Status()
	if status.File("Dockerfile").Worktree == git.Modified {
		_, _ = w.Add("Dockerfile")

		_, _ = w.Commit(version, &git.CommitOptions{
			Author: &object.Signature{
				Name:  "I Che",
				Email: "me@iche.eu",
				When:  time.Now(),
			},
		})

		_ = r.Push(&git.PushOptions{})
	}
}

func main() {
	log.Println("Fetching tags...")
	tags, err := fetchTags(gitURL)
	if err != nil {
		log.Fatal(err)
	}

	ver := getLastVersion(tags)
	log.Printf("Detected version: %v", ver)

	tmpl, err := template.New("Dockerfile").Parse(tpl)
	if err != nil {
		log.Fatal(err)
	}

	type H struct {
		Version string
	}
	h := H{Version: ver}

	var t bytes.Buffer
	tmpl.Execute(&t, h)

	// Write new docker file.
	ioutil.WriteFile("Dockerfile", t.Bytes(), 0644)
	commitLocal(ver)
}

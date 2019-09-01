# Hugo Docker Image

[Hugo](https://gohugo.io/) is a fast and flexible static site generator, written in Go.
Hugo flexibly works with many formats and is ideal for blogs, docs, portfolios and much more.
Hugo’s speed fosters creativity and makes building a website fun again.

This Lightweight Docker Image is based on Alpine, and comes with rsync for Continuous Deployment.

## Get Started

Print Hugo Help:

```bash
docker run --rm -ti chilic/docker-hugo help
```

Run any Hugo command:

```bash
docker run --rm -ti chilic/docker-hugo [command]
```

Create a new Hugo managed website:

```bash
docker run --rm -it -v $PWD:/src -u hugo chilic/docker-hugo new site mysite
cd mysite

# Now, you probably want to add a theme (see https://themes.gohugo.io/):
git init
git submodule add https://github.com/budparr/gohugo-theme-ananke.git themes/ananke;
echo 'theme = "ananke"' >> config.toml
```

Add some content:

```bash
docker run --rm -it -v $PWD:/src -u hugo chilic/docker-hugo new posts/my-first-post.md

# Now, you can edit this post, add your content and remove "draft" flag:
vim content/posts/my-first-post.md
```

Build your site:

```bash
docker run --rm -it -v $PWD:/src -u hugo chilic/docker-hugo -v
```

Serve your site locally:

```bash
docker run --rm -it -v $PWD:/src -p 1313:1313 -u hugo chilic/docker-hugo server -w --bind=0.0.0.0
```

Then open [`http://localhost:1313/`](http://localhost:1313/) in your browser.

To go further, read the [Hugo documentation](https://gohugo.io/documentation/).

## Continuous Deployment

I use this Docker image for Continuous Deployment. You can find some CI config examples in the [`ci-deploy`](https://github.com/chilic/docker-hugo/tree/master/ci-deploy) directory.

This Docker image also comes with:

- rsync
- git
- openssh-client

**Inspired by:** https://github.com/jguyomard/docker-hugo
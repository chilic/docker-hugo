FROM alpine:latest

MAINTAINER IChe <me@iche.eu>

RUN apk add --no-cache \
	curl \
	git \
	openssh-client \
	rsync

RUN mkdir -p /usr/local/src \
	&& cd /usr/local/src \
	&& curl -L https://github.com/gohugoio/hugo/releases/download/v0.63.2/hugo_0.63.2_linux-64bit.tar.gz | tar -xz \
	&& mv hugo /usr/local/bin/hugo \
	&& addgroup -Sg 1000 hugo \
	&& adduser -SG hugo -u 1000 -h /src hugo

WORKDIR /src

EXPOSE 1313

ENTRYPOINT ["/usr/local/bin/hugo"]
CMD [ "--help" ]

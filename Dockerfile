# Build Stage
FROM lacion/alpine-golang-buildimage:1.12.4 AS build-stage

LABEL app="build-voyager"
LABEL REPO="https://github.com/obsurvive/voyager"

ENV PROJPATH=/go/src/github.com/obsurvive/voyager

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/obsurvive/voyager
WORKDIR /go/src/github.com/obsurvive/voyager

RUN make build-alpine

# Final Stage
FROM lacion/alpine-base-image:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/obsurvive/voyager"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/voyager/bin

WORKDIR /opt/voyager/bin

COPY --from=build-stage /go/src/github.com/obsurvive/voyager/bin/voyager /opt/voyager/bin/
RUN chmod +x /opt/voyager/bin/voyager

# Create appuser
RUN adduser -D -g '' voyager
USER voyager

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/voyager/bin/voyager"]

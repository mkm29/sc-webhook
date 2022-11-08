FROM golang:1.18 AS build
LABEL maintainer="Mitch Murphy <mitch.murphy@gmail.com>"

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

WORKDIR /work
COPY . /work

# Build admission-webhook
RUN --mount=type=cache,target=/root/.cache/go-build,sharing=private \
  go build -o bin/admission-webhook .

# ---
FROM scratch AS run

ARG REGISTRY=""
ARG EXCLUDE_NAMESPACES="kube-system,kube-public,kube-node-lease"

COPY --from=build /work/bin/admission-webhook /usr/local/bin/

ENV REGISTRY="$REGISTRY" \
  EXCLUDE_NAMESPACES="$EXCLUDE_NAMESPACES"

CMD ["admission-webhook"]

FROM golang:1.18 AS build
LABEL maintainer="Mitch Murphy <mitch.murphy@gmail.com>"

ARG REGISTRY_BASE_URL = "docker.io"

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

COPY --from=build /work/bin/admission-webhook /usr/local/bin/

ENV REGISTRY_BASE_URL=${REGISTRY_BASE_URL}

CMD ["admission-webhook"]

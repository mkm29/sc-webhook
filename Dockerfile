FROM golang:1.16 AS build
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

COPY --from=build /work/bin/admission-webhook /usr/local/bin/

CMD ["admission-webhook"]

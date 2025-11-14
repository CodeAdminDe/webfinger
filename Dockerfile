LABEL org.opencontainers.image.source="https://github.com/codeadminde/webfinger"

FROM golang:1.25.4-alpine AS builder
WORKDIR /src
COPY ./src .

RUN dtsr="$(date '+%Y%m%d_%H%M%S')" && bvr="${GITHUB_SHA:-localbuild}" && buildid="$bvr-$dtsr" && sed -i "s|<<BUILD>>|$buildid|g" main.go

RUN CGO_ENABLED=0 GOOS=linux go build -o /main

FROM alpine:3.22.2
WORKDIR /
COPY --from=builder /main /main

RUN addgroup -g 1001 appusr && adduser -u 1001 -G appusr -D -H appusr && chown appusr:appusr /main

USER appusr

# ENV WEBFINGER_ISSUER_URL="true"
# ENV WEBFINGER_RESOURCE=acct:user@example.com
# ENV WEBFINGER_ISSUER_URL=https://example.com/application/issuer

EXPOSE 8080
ENTRYPOINT ["/main"]


ARG GO_VERSION=1.16.8
ARG ALPINE_VERSION=3.14
ARG ALPINECI_DOCKERTAG="awscli-latest-287984ed"
ARG APP_NAME="parzival"
ARG APP_PATH="/go/src/github.com/unfor19/parzival"

FROM unfor19/alpine-ci:${ALPINECI_DOCKERTAG} as awscli

# Dev
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS dev
RUN apk add --update git groff
COPY --from=awscli /usr/local/aws-cli/ /usr/local/aws-cli/
COPY --from=awscli /usr/local/bin/ /usr/local/bin/
ARG APP_NAME
ARG APP_PATH
ENV APP_NAME="${APP_NAME}" \
    APP_PATH="${APP_PATH}" \
    GOOS="linux"
WORKDIR "${APP_PATH}"
COPY . "${APP_PATH}"
ENTRYPOINT ["sh"]

# Build
FROM dev as build
ARG APP_NAME
ARG APP_PATH
RUN go mod download
RUN mkdir -p "/app/" && go build -o "/app/parzival"
ENTRYPOINT [ "sh" ]

# App
FROM alpine:${ALPINE_VERSION} AS app
WORKDIR "/app/"
COPY --from=build "/app/parzival" ./
RUN addgroup -S "appgroup" && adduser -S "appuser" -G "appgroup" && \
    chown "appuser":"appgroup" "parzival" && \
    chmod +x "parzival"
USER "appuser"
ENTRYPOINT ["/app/parzival"]

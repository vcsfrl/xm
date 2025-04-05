# Base Image
FROM golang:1.24-bookworm AS base
ARG username
ARG exec_user_id
RUN groupadd -g $exec_user_id -o $username
RUN useradd -r -u $exec_user_id -g $username $username -m
RUN mkdir -p /srv/xm
RUN chown $username:$username /srv/xm -R
USER $username:$username
WORKDIR /srv/xm

# Build Image
FROM base AS build
USER root:root
COPY . /srv/xm
RUN chown $username:$username /srv/xm -R
USER $username:$username
RUN go build -o /srv/xm/bin/app main.go

# Prod image
FROM golang:1.24-bookworm AS prod
COPY --from=build /srv/xm/bin/app /srv/xm/bin/app
HEALTHCHECK CMD curl --fail http://localhost:$XM_APP_PORT/api/v1/health
RUN go install github.com/divan/expvarmon@latest
CMD ["/srv/xm/bin/app","api"]

# Dev image
FROM base AS dev
CMD ["sleep", "365d"]


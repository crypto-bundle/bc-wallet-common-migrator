FROM golang:1.19-alpine AS gobuild

ENV GO111MODULE on
ENV GOSUMDB off
# add go-base repo to exceptions as a private repository.
ENV GOPRIVATE $GOPRIVATE,gitlab.heronodes.io/bc-platform

RUN apk add --no-cache bash git openssh build-base
RUN mkdir -p -m 0700 ~/.ssh && ssh-keyscan gitlab.heronodes.io >> ~/.ssh/known_hosts
RUN git config --global url."git@gitlab.heronodes.io:".insteadOf "https://gitlab.heronodes.io/"

WORKDIR /src

# Download and precompile all third party libraries, ignoring errors (some have broken tests or whatever).
COPY go.* ./

COPY . .

# Compile! Should only compile our sources since everything else is precompiled.
ARG RACE=-race
ARG CGO=1
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=ssh \
    mkdir -p /src/bin && \
    GOOS=linux CGO_ENABLED=${CGO} go build ${RACE} -v -installsuffix cgo -o ./bin/migrator -ldflags "-linkmode external -extldflags -static -s -w" ./cmd

FROM scratch

# Import the user and group files from the build stage.
#COPY --from=gobuild /etc/group /etc/passwd /etc/

ENV APP_ROOT /opt/appworker
ENV PATH /opt/appworker

COPY --from=gobuild /src/bin $APP_ROOT

USER appworker
CMD ["/opt/appworker/migrator"]
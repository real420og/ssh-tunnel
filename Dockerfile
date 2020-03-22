FROM golang:1.13-alpine as build

WORKDIR /src/
COPY main.go go.* /src/
RUN CGO_ENABLED=0 go build -o /bin/ssh-tunnel

FROM scratch
COPY --from=build /bin/ssh-tunnel /bin/ssh-tunnel
ENTRYPOINT ["/bin/ssh-tunnel"]
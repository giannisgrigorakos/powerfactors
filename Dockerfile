FROM golang:1.20.1 as builder

WORKDIR /go/powerfactors

COPY . ./
ARG version=dev
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X main.version=$version" -o powerfactors ./cmd/powerfactors/
RUN CGO_ENABLED=0 go build -o powerfactors ./cmd/powerfactors/


FROM scratch
COPY --from=builder /go/powerfactors/powerfactors .
CMD ["/powerfactors"]

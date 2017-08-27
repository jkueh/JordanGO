FROM golang:1.9 as compiler
WORKDIR /go/src/app
COPY . .
RUN go-wrapper download \
  && go-wrapper install \
  && go build -o /usr/local/sbin/jordango

FROM golang:1.9
COPY --from=compiler /usr/local/sbin/jordango ./jordango
ENTRYPOINT [ "./jordango" ]

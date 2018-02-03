# image with development tools
FROM xena/go:1.9.3
ENV GOPATH /root/go
RUN apk --no-cache add git protobuf retool
COPY . /root/go/src/github.com/Xe/lokahi
WORKDIR /root/go/src/github.com/Xe/lokahi
RUN retool build \
 && retool do mage -v generate build

# runner image
FROM xena/alpine
COPY --from=0 /root/go/src/github.com/Xe/lokahi/bin/lokahid /usr/local/bin/lokahid
COPY --from=0 /root/go/src/github.com/Xe/lokahi/bin/lokahictl /usr/local/bin/lokahictl
CMD /usr/local/bin/lokahid

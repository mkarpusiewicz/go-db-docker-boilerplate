ARG PROJECT
ARG ALPINE_VERSION=3.10.1

FROM ${PROJECT}-src AS builder
ARG VERSION
RUN make VERSION=${VERSION} -f build/build.Makefile


FROM alpine:$ALPINE_VERSION as sys-builder
RUN apk update && \
    apk add --no-cache ca-certificates tzdata && \
    update-ca-certificates
RUN adduser -D -g '' appuser

FROM alpine:$ALPINE_VERSION as entrypoint-builder
RUN apk update && \
    apk add --no-cache tini

FROM scratch as sys
#COPY --from=sys-builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=sys-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=sys-builder /etc/passwd /etc/passwd

FROM scratch as entrypoint
COPY --from=entrypoint-builder /sbin/tini /sbin/tini

FROM scratch
COPY --from=sys / /
COPY --from=entrypoint / /
COPY --from=builder /app/server /app/healthcheck /

USER appuser
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/server"]
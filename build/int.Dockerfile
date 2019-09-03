ARG PROJECT

FROM ${PROJECT}-src

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

RUN go test -c -tags "integration debug" ./cmd/server/integration

CMD ['sh', '-c', 'dockerize', '-timeout', '30s', '-wait', '$SERVER_URL/healthcheck', './integration.test', '-test.v', '-test.count=3']
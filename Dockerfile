FROM korchasa/go-build:0.1.0 as build
WORKDIR /app
ADD . .
ENV GOFLAGS="-mod=vendor"
RUN echo "#### Test" && go test -cover -v
RUN echo "#### Lint" && golangci-lint version && golangci-lint run -v
RUN echo "#### Build" && \
    gitget all && \
    go build -o app -ldflags "-s -w -X 'main.version=$(gitget version)' -X 'main.commit=$(gitget commit)' -X 'main.date=$(gitget date)'" .

FROM korchasa/go-app:0.0.1
WORKDIR /app
COPY --from=build /app/app ./app
USER nobody:nobody
ENTRYPOINT ["/app/app"]
CMD []
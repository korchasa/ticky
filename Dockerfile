FROM korchasa/go-build:0.0.1 as build
WORKDIR /app
ADD . .
ENV GOFLAGS="-mod=vendor"
RUN \
    echo "  ## Test" && \
    go test -cover -v -count=1 && \
    echo "  ## Lint" && \
    golangci-lint run ./... && \
    echo "  ## Build" && \
    gitget all && \
    go build -o app -ldflags "-s -w -X 'main.version=$(gitget version)' -X 'main.commit=$(gitget commit)' -X 'main.date=$(gitget date)'" . && \
    echo "  ## Done"

FROM korchasa/go-app:0.0.1
WORKDIR /app
COPY --from=build /app/app ./app
USER nobody:nobody
CMD ./app
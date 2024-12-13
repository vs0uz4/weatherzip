FROM golang:1.23.3 AS builder

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get dist-upgrade -y && \
    apt-get autoremove -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    echo "America/Sao_Paulo" > /etc/timezone && \
    ln -sf /usr/share/zoneinfo/America/Sao_Paulo /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata && \
    addgroup --system weatherzip && \ 
    adduser --system --ingroup weatherzip weatherzip

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./cmd      /app/cmd
COPY ./configs  /app/configs
COPY ./internal /app/internal

RUN mv /app/cmd/api/.env /app/.env && \
	CGO_ENABLED=0 GOOS=linux go build -o weatherzipapp ./cmd/api

USER weatherzip

FROM scratch

COPY --from=builder /app/.env .
COPY --from=builder /app/weatherzipapp .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/America/Sao_Paulo /usr/share/zoneinfo/America/Sao_Paulo

ENV TZ=America/Sao_Paulo

ENTRYPOINT [ "./weatherzipapp" ]
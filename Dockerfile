FROM golang:1.23.3

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
COPY ./pkg      /app/pkg

RUN mv /app/cmd/api/.env /app/.env && \
	CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

USER weatherzip

EXPOSE 8080

CMD ["./main"]

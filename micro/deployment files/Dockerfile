# FROM micro_temp

# WORKDIR /app

# COPY . .
# RUN go mod download
# RUN go build -o myapp main.go

# CMD ["./myapp"]

FROM ubuntu

WORKDIR /app

RUN apt update
RUN apt install -y golang
RUN apt install -y ca-certificates

COPY . .
RUN go mod download
RUN go run github.com/steebchen/prisma-client-go db push
RUN go build -o main main.go

CMD ["./main"]
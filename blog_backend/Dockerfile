FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Shanghai

WORKDIR /app

COPY server .
COPY configs ./configs

RUN mkdir -p uploads logs && chmod +x server

EXPOSE 8080

CMD ["./server"]

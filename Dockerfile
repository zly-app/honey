FROM golang:1.17 as builder
COPY . /src
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct CGO_ENABLED=0
RUN cd /src && \
 go build -o honey .

FROM scratch
COPY --from=builder /src/honey /src/honey
COPY --from=builder /src/configs /src/configs
WORKDIR /src
CMD ["./honey", "-c", "./configs/default.toml"]

EXPOSE 8080

# docker build -t zlyuan/honey .
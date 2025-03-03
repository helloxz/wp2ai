from debian:12-slim
RUN mkdir -p /opt/wp2ai/data
WORKDIR /opt/wp2ai
COPY wp2ai config.toml /opt/wp2ai/
COPY assets /opt/wp2ai/assets
# 暴露文件夹和端口
EXPOSE 2080
# 启动程序
CMD ["./wp2ai", "start"]
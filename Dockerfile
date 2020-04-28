# 指定基础镜像
FROM alpine:3.11


# 维护者信息
MAINTAINER yangrui


# DESCRIPTION "更改时区，替换apk源为USTC"

LABEL alpine_version=3.11.3
LABEL zoneinfo="Asia/Shanghai"
LABEL apk_repositoris="mirrors.ustc.edu.cn"

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

RUN rm -rf /var/cache/apk/* /tmp/* /var/tmp/* $HOME/.cache

# 参数
ARG SVC_NAME
COPY conf/config.yaml /
COPY bin/${SVC_NAME} /
RUN ln -s /${SVC_NAME} /app

CMD ["/app"]
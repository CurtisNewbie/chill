FROM alpine:3.17

LABEL author="Yongjie Zhuang"
LABEL descrption="Chill"

RUN apk --no-cache add tzdata

WORKDIR /usr/src/

# binary is pre-compiled
COPY app_chill ./app_chill

ENV TZ=Asia/Shanghai

CMD ["./app_chill", "configFile=/usr/src/config/conf.yml"]

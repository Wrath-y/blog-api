FROM golang:latest

MAINTAINER Ysama "ysama.cn"

ENV ACCESS_KEY_ID id
ENV ACCESS_KEY_SECRET secret
ENV HTTP_URL http://127.0.0.1:8081
ENV DB_URL http://127.0.0.1:3306
ENV MYSQL_ROOT_PASSWORD 123456
ENV JWT_SECRET jwt_secret
ENV GIT_USERNAME username
ENV GIT_PASSWORD 123456

WORKDIR /var/www

RUN cd /var/www; \
    touch ~/.git-credentials; \
    echo https://$GIT_USERNAME:$GIT_PASSWORD@gitee.com >> ~/.git-credentials; \
    git config --global credential.helper store; \
    git clone https://gitee.com/ysama/go-blog.git; \
    cd go-blog; \
    cp config/config.example.yaml config/config.yaml; \
    sed -in-place -e "s/url_example/$HTTP_URL/g" config/config.yaml; \
    sed -in-place -e "s/db_url_example/$DB_URL/g" config/config.yaml; \
    sed -in-place -e "s/db_password_example/$MYSQL_ROOT_PASSWORD/g" config/config.yaml; \
    sed -in-place -e "s/jwt_secret_example/$JWT_SECRET/g" config/config.yaml; \
    sed -in-place -e "s/accessKeyId_example/$ACCESS_KEY_ID/g" config/config.yaml; \
    sed -in-place -e "s/accessKeySecret_example/$ACCESS_KEY_SECRET/g" config/config.yaml; \
    go build ./

EXPOSE 8081

ENTRYPOINT ["./go-blog"]
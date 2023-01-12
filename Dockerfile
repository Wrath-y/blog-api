FROM golang:latest

ENV RUN_MODE release
ENV ACCESS_KEY_ID id
ENV ACCESS_KEY_SECRET secret
ENV HTTP_URL http://127.0.0.1:8081
ENV DB_URL http://127.0.0.1:3306
ENV MYSQL_ROOT_PASSWORD 123456
ENV JWT_SECRET jwt_secret

WORKDIR /var/www

RUN cd /var/www; \
    git clone https://github.com/Wrath-y/blog-api.git; \
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
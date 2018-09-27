FROM scratch
EXPOSE 8080
ENTRYPOINT ["/contrelease"]
COPY ./bin/ /

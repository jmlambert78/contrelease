FROM centos
MAINTAINER Jean Marc Lambert <jmlambert78@gmail.com>

RUN useradd -ms /bin/bash myuser
USER myuser
WORKDIR /home/myuser

ADD main main
EXPOSE 1323
ENTRYPOINT ["./main"]

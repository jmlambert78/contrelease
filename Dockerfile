FROM centos
MAINTAINER Jean Marc Lambert <jmlambert78@gmail.com>

RUN useradd -ms /bin/bash myuser
USER myuser
WORKDIR /home/myuser

ADD main main
ADD public public
EXPOSE 8080
ENTRYPOINT ["./main"]

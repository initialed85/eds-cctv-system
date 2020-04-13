FROM golang:1.13

RUN apt-get update && apt-get install -y cpulimit ffmpeg imagemagick python3 python3-pip

RUN pip3 install pytest python-dateutil

WORKDIR /srv/

COPY . /srv/

CMD ./native_build.sh && ./native_test.sh

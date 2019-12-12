FROM ubuntu:18.04 AS motion_build

RUN apt-get update && apt-get install -y \
    autoconf automake pkgconf libtool libjpeg8-dev build-essential libzip-dev gettext libmicrohttpd-dev \
    libavformat-dev libavcodec-dev libavutil-dev libswscale-dev libavdevice-dev git autopoint

WORKDIR /srv/

RUN git clone https://github.com/Motion-Project/motion.git

WORKDIR /srv/motion

RUN autoreconf -fiv

RUN ./configure --prefix=/

RUN make

RUN make install

FROM golang:1.13 AS go_build

WORKDIR /srv/

COPY cmd /srv/cmd
COPY internal /srv/internal
COPY pkg /srv/pkg
COPY go.mod /srv/go.mod
COPY go.sum /srv/go.sum

RUN go build -o motion_log_event_streamer cmd/motion_log_event_streamer/main.go

FROM ubuntu:18.04 AS run

ENV TZ Australia/Perth

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get install -y \
    tzdata supervisor libjpeg8-dev libmicrohttpd-dev libavformat-dev libavcodec-dev libavutil-dev libswscale-dev libavdevice-dev

RUN dpkg-reconfigure -f noninteractive tzdata

COPY --from=motion_build /bin/motion /bin/motion
COPY --from=motion_build /etc/motion/motion-dist.conf /etc/motion/examples/motion.conf
COPY --from=motion_build /etc/motion/camera1-dist.conf /etc/motion/examples/camera1.conf
COPY --from=motion_build /etc/motion/camera2-dist.conf /etc/motion/examples/camera2.conf
COPY --from=motion_build /etc/motion/camera3-dist.conf /etc/motion/examples/camera3.conf

RUN mkdir -p /etc/motion/conf.d

COPY --from=go_build /srv/motion_log_event_streamer /usr/local/bin/motion_log_event_streamer

COPY docker/eds-cctv-system-supervisor.conf /etc/supervisor/conf.d/

COPY docker/entrypoint.sh /srv/entrypoint.sh

# motion webcontrol
EXPOSE 8080

# motion stream
EXPOSE 8081

# motion_log_event_streamer
EXPOSE 8082

ENTRYPOINT ["/srv/entrypoint.sh"]

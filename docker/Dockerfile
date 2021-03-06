FROM linuxserver/ffmpeg AS base

RUN apt-get update && apt-get upgrade -y libfontconfig1 && apt-get install -y --reinstall \
    libfontconfig1 libfontconfig1-dev fontconfig-config

FROM base AS motion_build

RUN apt-get update && apt-get install -y \
    git autoconf automake pkgconf libtool build-essential libzip-dev gettext autopoint \
    libmicrohttpd-dev libjpeg8-dev libavutil-dev libavformat-dev libavcodec-dev libswscale-dev libavdevice-dev

WORKDIR /srv/
RUN git clone https://github.com/Motion-Project/motion.git

WORKDIR /srv/motion
RUN autoreconf -fiv
RUN ./configure --prefix=/
RUN make
RUN make install

FROM base AS run_base

ENV TZ Australia/Perth
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get install -y \
    tzdata supervisor imagemagick nginx logrotate python3-pip cpulimit x264 \
    libmicrohttpd-dev libjpeg8-dev libavutil-dev libavformat-dev libavcodec-dev libswscale-dev libavdevice-dev

RUN dpkg-reconfigure -f noninteractive tzdata

RUN pip3 install python-dateutil

FROM golang:1.13 AS go_build

WORKDIR /srv/
COPY cmd /srv/cmd
COPY internal /srv/internal
COPY pkg /srv/pkg
COPY go.mod /srv/go.mod
COPY go.sum /srv/go.sum

RUN go build -v -o motion_config_segment_recorder cmd/motion_config_segment_recorder/main.go
RUN go build -v -o motion_log_event_handler cmd/motion_log_event_handler/main.go
RUN go build -v -o segment_folder_event_handler cmd/segment_folder_event_handler/main.go
RUN go build -v -o event_store_updater_page_renderer cmd/event_store_updater_page_renderer/main.go
RUN go build -v -o static_file_server cmd/static_file_server/main.go
RUN go build -v -o event_store_deduplicator cmd/event_store_deduplicator/main.go

FROM run_base AS run

# router
EXPOSE 80

# motion web ui
EXPOSE 8080

# motion streams
EXPOSE 8081

# event api
EXPOSE 8082

# segment api
EXPOSE 8083

# static file server
EXPOSE 8084

COPY --from=motion_build /bin/motion /bin/motion
COPY --from=motion_build /etc/motion/motion-dist.conf /etc/motion/examples/motion.conf
COPY --from=motion_build /etc/motion/camera1-dist.conf /etc/motion/examples/camera1.conf
COPY --from=motion_build /etc/motion/camera2-dist.conf /etc/motion/examples/camera2.conf
COPY --from=motion_build /etc/motion/camera3-dist.conf /etc/motion/examples/camera3.conf

RUN mkdir -p /etc/motion/conf.d

COPY --from=go_build /srv/motion_config_segment_recorder /usr/local/bin/motion_config_segment_recorder
COPY --from=go_build /srv/motion_log_event_handler /usr/local/bin/motion_log_event_handler
COPY --from=go_build /srv/segment_folder_event_handler /usr/local/bin/segment_folder_event_handler
COPY --from=go_build /srv/event_store_updater_page_renderer /usr/local/bin/event_store_updater_page_renderer
COPY --from=go_build /srv/static_file_server /usr/local/bin/static_file_server
COPY --from=go_build /srv/event_store_deduplicator /usr/local/bin/event_store_deduplicator

COPY docker/nginx.conf /etc/nginx/nginx.conf
COPY docker/supervisor-apps.conf /etc/supervisor/conf.d/supervisor-apps.conf
COPY docker/entrypoint.sh /srv/entrypoint.sh
COPY utils/__init__.py /srv/utils/__init__.py
COPY utils/common.py /srv/utils/common.py
COPY utils/event_store_rebuilder_for_events.py /srv/utils/event_store_rebuilder_for_events.py
COPY utils/event_store_rebuilder_for_segments.py /srv/utils/event_store_rebuilder_for_segments.py

WORKDIR /srv/

ENTRYPOINT ["/srv/entrypoint.sh"]

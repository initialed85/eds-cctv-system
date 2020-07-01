# eds-cctv-system

This repo contains some Go code that integrates the following services into a rudimentary CCTV system:

NOTE: As of 1 Jul 2020 `willprice/nvidia-ffmpeg` is used as the base Docker image and so I can't promise this will work well on a machine withour an Nvidia GPU.

## Overview

- [Supervisor](https://github.com/Supervisor/supervisor)
    - to manage all processes
- [Motion](https://github.com/Motion-Project/motion)
    - to generate videos for "events" (detected motion events)
- [FFMpeg](https://github.com/FFmpeg/FFmpeg)
    - to generate videos for "segments" (5-minute segments)
    - to convert high-res videos to lo-res (for previews)
    - to pull thumbnails from videos
- [Cpulimit](https://github.com/opsengine/cpulimit)
    - limit CPU use for ffmpeg conversion
- [ImageMagick](https://github.com/ImageMagick/ImageMagick/)
    - to convert high-res thumbnails to low-res (for previews) 
- Go code
    - monitor the log from Motion and populate a datastore
    - monitor the folder for segments and populate a datastore
    - watch the event datastore and segment datastore and generate templated HTML
        - sorry about it- I never progressed past the "eating crayons" stage of front-end development
    - generate websocket events when something is added to the event datastore or segment datastore
        - not used by anything as yet (feel free to use externally!)
    
## How does it work?

I'll explain this by going through the processes in a running Docker container:

- `supervisord`
    - responsible for lifecycle of all processes (so, it's the root process)
- `static_file_server` (Go)
    - serve up the contents of /srv/target_dir (events and segments) to port 8084
    - TODO: I was lazy, this could probably be handled by Nginx
- `nginx`
    - router/proxy to internal services
         - `/ = file:///srv/root`
         - `/motion/ = http://127.0.0.1:8080/`
         - `/motion-stream/ = http://127.0.0.1:8081/`
         - `/events/ = http://127.0.0.1:8084/events/`
         - `/event_api/ = http://127.0.0.1:8082/` 
         - `/segments/ = http://127.0.0.1:8084/segments/`
         - `/segment_api/ = http://127.0.0.1:8083/` 
         - `/browse/ = file:///srv/target_dir/` 
- `logrotate_loop`
    - logrotate on an infinite loop
- `motion`
    - generate videos for detected motion events based on configuration at `/etc/motion`  
- `motion_log_event_handler` (Go)
    - watch logs from `motion` and identify motion events
    - generate low-res videos and thumbnails from high-res counterparts
    - write them to a datastore
    - expose them via web API
    - write them a websocket
- `event_store_updater_page_renderer` for events (Go)
    - watch events from a datastore
    - generate templated HTML with the events
- `motion_config_segment_recorder` (Go)
    - read the configs from `/etc/motion` and spawn `ffmpeg` instances to generate 5-minute video segments
- `segment_folder_event_handler` (Go)
    - watch the segments folder and identify segment events
    - generate low-res videos and thumbnails from high-res counterparts
    - write them to a datastore
    - expose them via web API
    - write them a websocket
- `event_store_updater_page_renderer` for segments (Go)
    - watch events from a datastore
    - generate templated HTML with the events

## How do I build it?

- in Docker
    - `./build.sh`
- natively
    - `./native_build.sh`

## How do I test it?

- in Docker
    - `./test.sh`
- natively
    - `./native_test.sh` 

## How do I run it?

- in Docker
    - `./deploy.sh` ensuring you've set the following environment variables
        - `CCTV_MOTION_CONFIGS` path to folder containing `motion.conf` and camera configs
            - see `motion-configs` for my configs or `motion-configs-examples` for the defaults
        - `CCTV_EVENTS_PATH` the path to store event videos and templated HTML
        - `CCTV_EVENTS_QUOTA` events path quota in GB
        - `CCTV_SEGMENTS_PATH` the path to store segments videos and templated HTML
        - `CCTV_SEGMENTS_QUOTA` segments path quota in GB
- natively
    - not recommended (if you desperately want to though, look through the Dockerfile to see what's done)

The quotas are managed by another Go tool I've written called [quotanizer](https://github.com/initialed85/quotanizer).

## How do I use it?

Once you've deployed the service (assuming localhost in this example), you can access the following URLs:

- http://localhost:81/events/events.html
    - index for templated events HTML
- http://localhost:81/event_api/events
    - all events as JSON
- http://localhost:81/event_api/events_by_date
    - all events as JSON, grouped by date
- ws://localhost:8082/stream
    - a WebSocket stream for any new events
    - you should be able to hit this from http://localhost:81/event_api/stream but you can't for some reason (connects, no messages)
- http://localhost:81/segments/events.html
    - index for templatd segments HTML
- http://localhost:81/segment_api/events
    - all segments as JSON
- http://localhost:81/segment_api/events_by_date
    - all segments as JSON, grouped by date
- ws://localhost:8083/stream
    - a WebSocket stream for any new segments
    - you should be able to hit this from http://localhost:81/segment_api/stream but you can't for some reason (connects, no messages)

## TODO

- Use Nginx for static_file_server
- Replace datastore with actual datastore (rather than hacky JSON Lines approach)
- Figure out why the WebSocket piece doesn't work through the Nginx proxy

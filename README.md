# eds-cctv-system

This repo contains some Go code that integrates the following services into a rudimentary CCTV system:

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
    - `./docker_build.sh`
- natively
    - `./build.sh`

## How do I test it?

- in Docker
    - `./docker_test.sh`
- natively
    - `./docker_build.sh` 

## How do I run it?

- in Docker
    - `./docker_run.sh` or `./docker_run_in_background.sh` (note: the latter is configured for my own setup)
- natively
    - not recommended (if you desperately want to though, look through the Dockerfile to see what's done)

## TODO

- Use Nginx for static_file_server
- Replace datastore with actual datastore (rather than hacky JSON Lines approach)

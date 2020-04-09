import datetime
import json
import os
from pathlib import Path
from types import SimpleNamespace
from typing import List
from typing import NamedTuple, Union, Optional, Callable
from uuid import uuid3, NAMESPACE_DNS

from dateutil.parser import parse

_VIDEO_SUFFIXES = [".mkv", ".mp4"]
_IMAGE_SUFFIXES = [".jpg"]
_PERMITTED_EXTENSIONS = _VIDEO_SUFFIXES + _IMAGE_SUFFIXES


class PathDetails(NamedTuple):
    path: Path
    event_id: Optional[int]
    camera_id: Optional[int]
    timestamp: datetime.datetime
    camera_name: str
    is_image: bool
    is_lowres: bool


class Event(SimpleNamespace):
    event_id: str
    timestamp: Union[datetime.datetime, str]
    camera_name: str
    high_res_image_path: str
    low_res_image_path: str
    high_res_video_path: str
    low_res_video_path: str


def get_sorted_paths(path: Path) -> List[Path]:
    return sorted(Path(path).iterdir(), key=os.path.getmtime)


def format_timestamp_for_go(timestamp: Union[datetime.datetime, str]) -> str:
    if isinstance(timestamp, str):
        timestamp = parse(timestamp)

    us = timestamp.strftime("%f")

    tz_raw = timestamp.strftime("%z")
    tz = "{}:{}".format(tz_raw[0:3], tz_raw[3:])

    return timestamp.strftime(f"%Y-%m-%dT%H:%M:%S.{us}00{tz}")


def parse_paths(paths: List[Path], tzinfo: datetime.tzinfo, parse_method: Callable) -> List[PathDetails]:
    return [
        y
        for y in [parse_method(path=x, tzinfo=tzinfo) for x in paths if x is not None]
        if y is not None
    ]


def build_event_for_some_path_details(some_path_details: List[PathDetails], path: Path):
    if len(some_path_details) != 4:
        raise ValueError(
            f"expected some_path_details to be 4 long (and related); instead it was {len(some_path_details)} long"
        )

    event_ids = list(set([x.event_id for x in some_path_details]))
    if len(event_ids) != 1:
        raise ValueError(
            f"expected all PathDetails to have a common event_id; instead they were {event_ids}"
        )

    camera_ids = list(set([x.camera_id for x in some_path_details]))
    if len(camera_ids) != 1:
        raise ValueError(
            f"expected all PathDetails to have a common camera_id; instead they were {camera_ids}"
        )

    camera_names = list(set([x.camera_name for x in some_path_details]))
    if len(camera_names) != 1:
        raise ValueError(
            f"expected all PathDetails to have a common camera_name; instead they were {camera_names}"
        )

    high_res_image_paths = list(
        set([x.path for x in some_path_details if x.is_image and not x.is_lowres])
    )
    if len(high_res_image_paths) != 1:
        raise ValueError(
            f"expected to find 1 high_res_image_path from PathDetails; instead found {high_res_image_paths}"
        )

    low_res_image_paths = list(
        set([x.path for x in some_path_details if x.is_image and x.is_lowres])
    )
    if len(low_res_image_paths) != 1:
        raise ValueError(
            f"expected to find 1 low_res_image_path from PathDetails; instead found {low_res_image_paths}"
        )

    high_res_video_paths = list(
        set([x.path for x in some_path_details if not x.is_image and not x.is_lowres])
    )
    if len(high_res_video_paths) != 1:
        raise ValueError(
            f"expected to find 1 high_res_video_path from PathDetails; instead found {high_res_video_paths}"
        )

    low_res_video_paths = list(
        set([x.path for x in some_path_details if not x.is_image and x.is_lowres])
    )
    if len(low_res_video_paths) != 1:
        raise ValueError(
            f"expected to find 1 low_res_video_path from PathDetails; instead found {low_res_video_paths}"
        )

    timestamp = sorted([x.timestamp for x in some_path_details])[0]
    high_res_image_path = high_res_image_paths[0]
    low_res_image_path = low_res_image_paths[0]
    high_res_video_path = high_res_video_paths[0]
    low_res_video_path = low_res_video_paths[0]

    # in Go:
    # eventId := uuid.NewSHA1(
    #     uuid.NameSpaceDNS,
    #     []byte(fmt.Sprintf("%v, %v, %v, %v, %v", timestamp, highResImagePath, lowResImagePath, highResVideoPath, lowResVideoPath)),
    # )

    event_id = uuid3(
        NAMESPACE_DNS,
        f"{format_timestamp_for_go(timestamp)}, {high_res_image_path}, {low_res_image_path}, {high_res_video_path}, {low_res_video_path}",
    )

    return Event(
        event_id=str(event_id),
        timestamp=timestamp,
        camera_name=camera_names[0],
        high_res_image_path=str(path / high_res_image_path),
        low_res_image_path=str(path / low_res_image_path),
        high_res_video_path=str(path / high_res_video_path),
        low_res_video_path=str(path / low_res_video_path),
    )


def relate_path_details(
        some_path_details: List[PathDetails],
        get_key_methods: List[Callable]
) -> List[List[PathDetails]]:
    some_path_details_by_key = {}
    for path_details in some_path_details:
        keys = [x(path_details) for x in get_key_methods]

        for key in keys:
            some_path_details_by_key.setdefault(key, [])
            some_path_details_by_key[key] += [path_details]

    viable_some_path_details_by_key = {
        k: v for k, v in some_path_details_by_key.items() if len(v) == 4
    }

    deduplicated_path_details = []
    for some_path_details in viable_some_path_details_by_key.values():
        if some_path_details not in deduplicated_path_details:
            deduplicated_path_details += [some_path_details]

    return deduplicated_path_details


def build_events_for_related_path_details(
        related_path_details: List[List[PathDetails]], path: Path
) -> List[Event]:
    events: List[Event] = []
    for some_path_details in related_path_details:
        events += [
            build_event_for_some_path_details(
                some_path_details=some_path_details, path=path
            )
        ]

    sorted_events = sorted(events, key=lambda x: x.timestamp)

    for event in sorted_events:
        event.timestamp = format_timestamp_for_go(timestamp=event.timestamp)

    return sorted_events


def build_json_lines_from_events(events: List[Event]) -> str:
    return "\n".join(
        [
            json.dumps(
                {
                    "event_id": x.event_id,
                    "timestamp": x.timestamp,
                    "camera_name": x.camera_name,
                    "high_res_image_path": x.high_res_image_path,
                    "low_res_image_path": x.low_res_image_path,
                    "high_res_video_path": x.high_res_video_path,
                    "low_res_video_path": x.low_res_video_path,
                }
            )
            for x in events
        ]
    )


def write_to_file(path: Path, data: str):
    with open(str(path), "w") as f:
        f.write(data)


def rebuild_event_store(root_path: Path, tzinfo: datetime.tzinfo, json_path: Path, parse_method: Callable, get_key_methods: List[Callable]):
    print(f"getting sorted paths from {root_path}...")
    sorted_paths = get_sorted_paths(path=root_path)
    print(f"got {len(sorted_paths)} sorted paths")

    print("parsing sorted paths...")
    some_path_details = parse_paths(paths=sorted_paths, tzinfo=tzinfo, parse_method=parse_method)
    print(f"got {len(some_path_details)} parsed paths")

    print("relating parsed paths...")
    related_path_details = relate_path_details(some_path_details=some_path_details,
                                               get_key_methods=get_key_methods)
    print(f"got {len(related_path_details)} related paths")

    print("building events...")
    events = build_events_for_related_path_details(
        related_path_details=related_path_details, path=root_path
    )
    print(f"built {len(events)} events")

    print("building json lines...")
    json_lines = build_json_lines_from_events(events=events)
    print(f"built {len(json_lines)} bytes")

    print(f"writing to {json_path}")
    write_to_file(path=json_path, data=json_lines)
    print("done.")

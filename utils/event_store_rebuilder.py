#!/usr/bin/python3

import datetime
from pathlib import Path
from types import SimpleNamespace
from typing import NamedTuple, Tuple, Dict, List, Optional, Union
from uuid import uuid3, NAMESPACE_DNS

import json
import os

_VIDEO_SUFFIXES = [".mkv", ".mp4"]
_IMAGE_SUFFIXES = [".jpg"]
_PERMITTED_EXTENSIONS = _VIDEO_SUFFIXES + _IMAGE_SUFFIXES


class PathDetails(NamedTuple):
    path: Path
    event_id: int
    camera_id: int
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


def parse_path(path: Path, tzinfo: datetime.tzinfo) -> Optional[PathDetails]:
    if path.suffix.lower() not in _PERMITTED_EXTENSIONS:
        return None

    if path.name.lower().startswith("segment"):
        raise ValueError("cannot process segments; only events")

    parts = path.name.split("__")

    event_id = int(parts[0])

    camera_id = int(parts[1])

    timestamp = datetime.datetime.strptime(parts[2], "%Y-%m-%d_%H-%M-%S")
    timestamp = timestamp.replace(tzinfo=tzinfo)

    camera_name = parts[3].split(".")[0].split("-")[0]

    return PathDetails(
        path=path,
        event_id=event_id,
        camera_id=camera_id,
        timestamp=timestamp,
        camera_name=camera_name,
        is_image=path.suffix.lower() in _IMAGE_SUFFIXES,
        is_lowres="-lowres" in path.name.lower(),
    )


def parse_paths(paths: List[Path], tzinfo: datetime.tzinfo) -> List[PathDetails]:
    return [parse_path(path=x, tzinfo=tzinfo) for x in paths if x is not None]


def _get_key(path_details: PathDetails) -> Tuple[int, int, str]:
    return (path_details.event_id, path_details.camera_id, path_details.camera_name)


def relate_path_details(
    some_path_details: List[PathDetails],
) -> Dict[Tuple[int, int, str], List[PathDetails]]:
    some_path_details_by_key: Dict[Tuple[int, int, str], List[PathDetails]] = {}
    for some_path_details in some_path_details:
        key = _get_key(some_path_details)
        some_path_details_by_key.setdefault(key, [])
        some_path_details_by_key[key] += [some_path_details]

    lens = [len(x) for x in some_path_details_by_key.values()]

    # most likely matched everything- nice; return early
    if len(lens) > 0 and min(lens) == 4 and max(lens) == 4:
        return some_path_details_by_key

    raise RuntimeError(
        "min matches = {}, max matches = {}".format(min(lens), max(lens))
    )


def format_timestamp_for_go(timestamp: datetime.datetime) -> str:
    us = timestamp.strftime("%f")

    tz_raw = timestamp.strftime("%z")
    tz = "{}:{}".format(tz_raw[0:3], tz_raw[3:])

    return timestamp.strftime(f"%Y-%m-%dT%H:%M:%S.{us}00{tz}")


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


def build_events_for_path_details_by_key(
    some_path_details_by_key: Dict[Tuple[int, int, str], List[PathDetails]], path: Path
) -> List[Event]:
    events: List[Event] = []
    for some_path_details in some_path_details_by_key.values():
        events += [
            build_event_for_some_path_details(
                some_path_details=some_path_details, path=path
            )
        ]

    sorted_events: List[Tuple[datetime.datetime, Event]] = sorted(
        (x.timestamp, x) for x in events
    )

    for _, event in sorted_events:
        event.timestamp = str(event.timestamp)

    return [x[1] for x in sorted_events]


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


def rebuild_event_store(root_path: Path, tzinfo: datetime.tzinfo, json_path: Path):
    sorted_paths = get_sorted_paths(path=root_path)

    some_path_details = parse_paths(paths=sorted_paths, tzinfo=tzinfo)

    some_path_details_by_key = relate_path_details(some_path_details=some_path_details)

    events = build_events_for_path_details_by_key(
        some_path_details_by_key=some_path_details_by_key, path=root_path
    )

    json_lines = build_json_lines_from_events(events=events)

    write_to_file(path=json_path, data=json_lines)

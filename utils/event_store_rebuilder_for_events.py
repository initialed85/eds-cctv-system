#!/usr/bin/python3

import datetime
from pathlib import Path
from typing import Tuple, Optional

from common import _IMAGE_SUFFIXES, _PERMITTED_EXTENSIONS, PathDetails, rebuild_event_store


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


def _get_key_pass_1(path_details: PathDetails) -> Tuple[int, int, str]:
    return (
        path_details.event_id,
        path_details.camera_id,
        path_details.camera_name,
    )


def _get_key_pass_2(path_details: PathDetails) -> Tuple[int, int, str, str]:
    return (
        path_details.event_id,
        path_details.camera_id,
        path_details.camera_name,
        path_details.timestamp.strftime("%Y-%m-%d %H"),
    )


def _get_key_pass_3(path_details: PathDetails) -> Tuple[int, int, str, str]:
    return (
        path_details.event_id,
        path_details.camera_id,
        path_details.camera_name,
        path_details.timestamp.strftime("%Y-%m-%d %H:%M"),
    )


if __name__ == "__main__":
    import argparse
    from dateutil.tz import tzoffset

    parser = argparse.ArgumentParser()

    parser.add_argument("-r", "--root-path", type=str, required=True)
    parser.add_argument("-j", "--json-path", type=str, required=True)

    args = parser.parse_args()

    rebuild_event_store(
        root_path=args.root_path,
        tzinfo=tzoffset(name="WST-8", offset=8 * 60 * 60),
        json_path=args.json_path,
        parse_method=parse_path,
        get_key_methods=[_get_key_pass_1, _get_key_pass_2, _get_key_pass_3]
    )

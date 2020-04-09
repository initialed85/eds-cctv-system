#!/usr/bin/python3

import datetime
from pathlib import Path
from typing import Optional, Tuple

from common import _IMAGE_SUFFIXES, _PERMITTED_EXTENSIONS, PathDetails, rebuild_event_store


def parse_path(path: Path, tzinfo: datetime.tzinfo) -> Optional[PathDetails]:
    if path.suffix.lower() not in _PERMITTED_EXTENSIONS:
        return None

    if path.name.lower().startswith("event"):
        raise ValueError("cannot process events; only segments")

    parts = path.name.split('_')

    timestamp = datetime.datetime.strptime(f'{parts[1]}_{parts[2]}', "%Y-%m-%d_%H-%M-%S")
    timestamp = timestamp.replace(tzinfo=tzinfo)

    camera_name = parts[3].split('.')[0]
    if camera_name.endswith('-lowres'):
        camera_name = camera_name.split('-lowres')[0]

    return PathDetails(
        path=path,
        event_id=None,
        camera_id=None,
        timestamp=timestamp,
        camera_name=camera_name,
        is_image=path.suffix.lower() in _IMAGE_SUFFIXES,
        is_lowres="-lowres" in path.name.lower(),
    )


def _get_key(path_details: PathDetails) -> Tuple[str, str]:
    return (
        path_details.camera_name,
        path_details.timestamp.strftime("%Y-%m-%d %H:%M:%S")
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
        get_key_methods=[_get_key]
    )

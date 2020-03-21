import datetime
import unittest
from pathlib import Path, PosixPath

from dateutil.tz import tzoffset
from tempfile import mkdtemp

from .event_store_rebuilder import (
    get_sorted_paths,
    PathDetails,
    parse_path,
    parse_paths,
    relate_path_details,
    build_event_for_some_path_details,
    format_timestamp_for_go,
    build_events_for_related_path_details,
    Event,
    build_json_lines_from_events,
    rebuild_event_store,
)
from .test_data import _sorted_paths_events

_TZINFO = tzoffset(name="WST-8", offset=8 * 60 * 60)

_TEST_TIMESTAMP = datetime.datetime(
    year=1992,
    month=2,
    day=6,
    hour=11,
    minute=11,
    second=11,
    microsecond=11,
    tzinfo=_TZINFO,
)
_TEST_TIMESTAMP_IN_GO_FORMAT = "1992-02-06T11:11:11.00001100+08:00"

_SOME_PATHS = _sorted_paths_events[0:64]

"""_SOME_PATHS = [
    PosixPath('01__102__2020-03-03_00-26-28__FrontDoor.mkv'),
    PosixPath('01__102__2020-03-03_00-26-30__FrontDoor.jpg'),
    PosixPath('01__102__2020-03-03_00-26-28__FrontDoor-lowres.mkv'),
    PosixPath('01__102__2020-03-03_00-26-30__FrontDoor-lowres.jpg'),
    PosixPath('02__102__2020-03-03_00-27-15__FrontDoor.jpg'),
    PosixPath('02__102__2020-03-03_00-27-13__FrontDoor.mkv'),
    PosixPath('02__102__2020-03-03_00-27-13__FrontDoor-lowres.mkv'),
    PosixPath('02__102__2020-03-03_00-27-15__FrontDoor-lowres.jpg'),
    PosixPath('01__101__2020-03-03_00-43-24__Driveway.jpg'),
    PosixPath('01__101__2020-03-03_00-43-22__Driveway.mkv'),
    PosixPath('01__101__2020-03-03_00-43-22__Driveway-lowres.mkv'),
    PosixPath('01__101__2020-03-03_00-43-24__Driveway-lowres.jpg'),
    PosixPath('02__101__2020-03-03_00-47-16__Driveway.jpg'),
    PosixPath('02__101__2020-03-03_00-47-14__Driveway.mkv'),
    PosixPath('02__101__2020-03-03_00-47-14__Driveway-lowres.mkv'),
    PosixPath('02__101__2020-03-03_00-47-16__Driveway-lowres.jpg'),
    PosixPath('03__101__2020-03-03_00-49-52__Driveway.jpg'),
    PosixPath('03__101__2020-03-03_00-49-50__Driveway.mkv'),
    PosixPath('03__101__2020-03-03_00-49-50__Driveway-lowres.mkv'),
    PosixPath('03__101__2020-03-03_00-49-52__Driveway-lowres.jpg'),
    PosixPath('04__101__2020-03-03_01-00-39__Driveway.jpg'),
    PosixPath('04__101__2020-03-03_01-00-37__Driveway.mkv'),
    PosixPath('04__101__2020-03-03_01-00-37__Driveway-lowres.mkv'),
    PosixPath('04__101__2020-03-03_01-00-39__Driveway-lowres.jpg'),
    PosixPath('05__101__2020-03-03_01-01-24__Driveway.jpg'),
    PosixPath('05__101__2020-03-03_01-01-23__Driveway.mkv'),
    PosixPath('05__101__2020-03-03_01-01-23__Driveway-lowres.mkv'),
    PosixPath('05__101__2020-03-03_01-01-24__Driveway-lowres.jpg'),
    PosixPath('03__102__2020-03-03_01-01-51__FrontDoor.mkv'),
    PosixPath('03__102__2020-03-03_01-01-53__FrontDoor.jpg'),
    PosixPath('03__102__2020-03-03_01-01-51__FrontDoor-lowres.mkv'),
    PosixPath('03__102__2020-03-03_01-01-53__FrontDoor-lowres.jpg'),
    PosixPath('06__101__2020-03-03_01-02-03__Driveway.jpg'),
    PosixPath('06__101__2020-03-03_01-02-01__Driveway.mkv'),
    PosixPath('06__101__2020-03-03_01-02-01__Driveway-lowres.mkv'),
    PosixPath('06__101__2020-03-03_01-02-03__Driveway-lowres.jpg'),
    PosixPath('07__101__2020-03-03_01-02-30__Driveway.jpg'),
    PosixPath('07__101__2020-03-03_01-02-28__Driveway.mkv'),
    PosixPath('07__101__2020-03-03_01-02-28__Driveway-lowres.mkv'),
    PosixPath('07__101__2020-03-03_01-02-30__Driveway-lowres.jpg'),
    PosixPath('08__101__2020-03-03_01-04-35__Driveway.jpg'),
    PosixPath('08__101__2020-03-03_01-04-34__Driveway.mkv'),
    PosixPath('08__101__2020-03-03_01-04-34__Driveway-lowres.mkv'),
    PosixPath('08__101__2020-03-03_01-04-35__Driveway-lowres.jpg'),
    PosixPath('04__102__2020-03-03_01-05-05__FrontDoor.mkv'),
    PosixPath('04__102__2020-03-03_01-05-07__FrontDoor.jpg'),
    PosixPath('04__102__2020-03-03_01-05-05__FrontDoor-lowres.mkv'),
    PosixPath('04__102__2020-03-03_01-05-07__FrontDoor-lowres.jpg'),
    PosixPath('09__101__2020-03-03_01-27-27__Driveway.jpg'),
    PosixPath('09__101__2020-03-03_01-27-25__Driveway.mkv'),
    PosixPath('09__101__2020-03-03_01-27-25__Driveway-lowres.mkv'),
    PosixPath('09__101__2020-03-03_01-27-27__Driveway-lowres.jpg'),
    PosixPath('10__101__2020-03-03_01-48-27__Driveway.jpg'),
    PosixPath('10__101__2020-03-03_01-48-26__Driveway.mkv'),
    PosixPath('10__101__2020-03-03_01-48-26__Driveway-lowres.mkv'),
    PosixPath('10__101__2020-03-03_01-48-27__Driveway-lowres.jpg'),
    PosixPath('11__101__2020-03-03_01-54-19__Driveway.jpg'),
    PosixPath('11__101__2020-03-03_01-54-17__Driveway.mkv'),
    PosixPath('11__101__2020-03-03_01-54-17__Driveway-lowres.mkv'),
    PosixPath('11__101__2020-03-03_01-54-19__Driveway-lowres.jpg'),
    PosixPath('12__101__2020-03-03_01-58-52__Driveway.jpg'),
    PosixPath('12__101__2020-03-03_01-58-50__Driveway.mkv'),
    PosixPath('12__101__2020-03-03_01-58-50__Driveway-lowres.mkv'),
    PosixPath('12__101__2020-03-03_01-58-52__Driveway-lowres.jpg')
]"""

_PATH_DETAILS = PathDetails(
    path=PosixPath("01__102__2020-03-03_00-26-30__FrontDoor-lowres.jpg"),
    event_id=1,
    camera_id=102,
    timestamp=datetime.datetime(2020, 3, 3, 0, 26, 30, tzinfo=tzoffset("WST-8", 28800)),
    camera_name="FrontDoor",
    is_image=True,
    is_lowres=True,
)

_SOME_PATH_DETAILS = [
    PathDetails(
        path=PosixPath("01__102__2020-03-03_00-26-28__FrontDoor.mkv"),
        event_id=1,
        camera_id=102,
        timestamp=datetime.datetime(
            2020, 3, 3, 0, 26, 28, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="FrontDoor",
        is_image=False,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("01__102__2020-03-03_00-26-30__FrontDoor.jpg"),
        event_id=1,
        camera_id=102,
        timestamp=datetime.datetime(
            2020, 3, 3, 0, 26, 30, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="FrontDoor",
        is_image=True,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("01__102__2020-03-03_00-26-28__FrontDoor-lowres.mkv"),
        event_id=1,
        camera_id=102,
        timestamp=datetime.datetime(
            2020, 3, 3, 0, 26, 28, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="FrontDoor",
        is_image=False,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("01__102__2020-03-03_00-26-30__FrontDoor-lowres.jpg"),
        event_id=1,
        camera_id=102,
        timestamp=datetime.datetime(
            2020, 3, 3, 0, 26, 30, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="FrontDoor",
        is_image=True,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("02__102__2020-03-03_00-27-15__FrontDoor.jpg"),
        event_id=2,
        camera_id=102,
        timestamp=datetime.datetime(
            2020, 3, 3, 0, 27, 15, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="FrontDoor",
        is_image=True,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("02__102__2020-03-03_00-27-13__FrontDoor.mkv"),
        event_id=2,
        camera_id=102,
        timestamp=datetime.datetime(
            2020, 3, 3, 0, 27, 13, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="FrontDoor",
        is_image=False,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("02__102__2020-03-03_00-27-13__FrontDoor-lowres.mkv"),
        event_id=2,
        camera_id=102,
        timestamp=datetime.datetime(
            2020, 3, 3, 0, 27, 13, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="FrontDoor",
        is_image=False,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("02__102__2020-03-03_00-27-15__FrontDoor-lowres.jpg"),
        event_id=2,
        camera_id=102,
        timestamp=datetime.datetime(
            2020, 3, 3, 0, 27, 15, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="FrontDoor",
        is_image=True,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("01__101__2020-03-03_00-43-24__Driveway.jpg"),
        event_id=1,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 0, 43, 24, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("01__101__2020-03-03_00-43-22__Driveway.mkv"),
        event_id=1,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 0, 43, 22, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("01__101__2020-03-03_00-43-22__Driveway-lowres.mkv"),
        event_id=1,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 0, 43, 22, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("01__101__2020-03-03_00-43-24__Driveway-lowres.jpg"),
        event_id=1,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 0, 43, 24, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("02__101__2020-03-03_00-47-16__Driveway.jpg"),
        event_id=2,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 0, 47, 16, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("02__101__2020-03-03_00-47-14__Driveway.mkv"),
        event_id=2,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 0, 47, 14, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("02__101__2020-03-03_00-47-14__Driveway-lowres.mkv"),
        event_id=2,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 0, 47, 14, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("02__101__2020-03-03_00-47-16__Driveway-lowres.jpg"),
        event_id=2,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 0, 47, 16, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("03__101__2020-03-03_00-49-52__Driveway.jpg"),
        event_id=3,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 0, 49, 52, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("03__101__2020-03-03_00-49-50__Driveway.mkv"),
        event_id=3,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 0, 49, 50, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("03__101__2020-03-03_00-49-50__Driveway-lowres.mkv"),
        event_id=3,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 0, 49, 50, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("03__101__2020-03-03_00-49-52__Driveway-lowres.jpg"),
        event_id=3,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 0, 49, 52, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("04__101__2020-03-03_01-00-39__Driveway.jpg"),
        event_id=4,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 0, 39, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("04__101__2020-03-03_01-00-37__Driveway.mkv"),
        event_id=4,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 0, 37, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("04__101__2020-03-03_01-00-37__Driveway-lowres.mkv"),
        event_id=4,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 0, 37, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("04__101__2020-03-03_01-00-39__Driveway-lowres.jpg"),
        event_id=4,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 0, 39, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("05__101__2020-03-03_01-01-24__Driveway.jpg"),
        event_id=5,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 1, 24, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("05__101__2020-03-03_01-01-23__Driveway.mkv"),
        event_id=5,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 1, 23, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("05__101__2020-03-03_01-01-23__Driveway-lowres.mkv"),
        event_id=5,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 1, 23, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("05__101__2020-03-03_01-01-24__Driveway-lowres.jpg"),
        event_id=5,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 1, 24, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("03__102__2020-03-03_01-01-51__FrontDoor.mkv"),
        event_id=3,
        camera_id=102,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 1, 51, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="FrontDoor",
        is_image=False,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("03__102__2020-03-03_01-01-53__FrontDoor.jpg"),
        event_id=3,
        camera_id=102,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 1, 53, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="FrontDoor",
        is_image=True,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("03__102__2020-03-03_01-01-51__FrontDoor-lowres.mkv"),
        event_id=3,
        camera_id=102,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 1, 51, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="FrontDoor",
        is_image=False,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("03__102__2020-03-03_01-01-53__FrontDoor-lowres.jpg"),
        event_id=3,
        camera_id=102,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 1, 53, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="FrontDoor",
        is_image=True,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("06__101__2020-03-03_01-02-03__Driveway.jpg"),
        event_id=6,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 2, 3, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("06__101__2020-03-03_01-02-01__Driveway.mkv"),
        event_id=6,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 2, 1, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("06__101__2020-03-03_01-02-01__Driveway-lowres.mkv"),
        event_id=6,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 2, 1, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("06__101__2020-03-03_01-02-03__Driveway-lowres.jpg"),
        event_id=6,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 2, 3, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("07__101__2020-03-03_01-02-30__Driveway.jpg"),
        event_id=7,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 2, 30, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("07__101__2020-03-03_01-02-28__Driveway.mkv"),
        event_id=7,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 2, 28, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("07__101__2020-03-03_01-02-28__Driveway-lowres.mkv"),
        event_id=7,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 2, 28, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("07__101__2020-03-03_01-02-30__Driveway-lowres.jpg"),
        event_id=7,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 2, 30, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("08__101__2020-03-03_01-04-35__Driveway.jpg"),
        event_id=8,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 4, 35, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("08__101__2020-03-03_01-04-34__Driveway.mkv"),
        event_id=8,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 4, 34, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("08__101__2020-03-03_01-04-34__Driveway-lowres.mkv"),
        event_id=8,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 4, 34, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("08__101__2020-03-03_01-04-35__Driveway-lowres.jpg"),
        event_id=8,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 4, 35, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("04__102__2020-03-03_01-05-05__FrontDoor.mkv"),
        event_id=4,
        camera_id=102,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 5, 5, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="FrontDoor",
        is_image=False,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("04__102__2020-03-03_01-05-07__FrontDoor.jpg"),
        event_id=4,
        camera_id=102,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 5, 7, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="FrontDoor",
        is_image=True,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("04__102__2020-03-03_01-05-05__FrontDoor-lowres.mkv"),
        event_id=4,
        camera_id=102,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 5, 5, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="FrontDoor",
        is_image=False,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("04__102__2020-03-03_01-05-07__FrontDoor-lowres.jpg"),
        event_id=4,
        camera_id=102,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 5, 7, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="FrontDoor",
        is_image=True,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("09__101__2020-03-03_01-27-27__Driveway.jpg"),
        event_id=9,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 27, 27, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("09__101__2020-03-03_01-27-25__Driveway.mkv"),
        event_id=9,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 27, 25, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("09__101__2020-03-03_01-27-25__Driveway-lowres.mkv"),
        event_id=9,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 27, 25, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("09__101__2020-03-03_01-27-27__Driveway-lowres.jpg"),
        event_id=9,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 27, 27, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("10__101__2020-03-03_01-48-27__Driveway.jpg"),
        event_id=10,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 48, 27, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("10__101__2020-03-03_01-48-26__Driveway.mkv"),
        event_id=10,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 48, 26, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("10__101__2020-03-03_01-48-26__Driveway-lowres.mkv"),
        event_id=10,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 48, 26, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("10__101__2020-03-03_01-48-27__Driveway-lowres.jpg"),
        event_id=10,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 48, 27, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("11__101__2020-03-03_01-54-19__Driveway.jpg"),
        event_id=11,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 54, 19, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("11__101__2020-03-03_01-54-17__Driveway.mkv"),
        event_id=11,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 54, 17, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("11__101__2020-03-03_01-54-17__Driveway-lowres.mkv"),
        event_id=11,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 54, 17, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("11__101__2020-03-03_01-54-19__Driveway-lowres.jpg"),
        event_id=11,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 54, 19, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("12__101__2020-03-03_01-58-52__Driveway.jpg"),
        event_id=12,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 58, 52, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("12__101__2020-03-03_01-58-50__Driveway.mkv"),
        event_id=12,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 58, 50, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=False,
    ),
    PathDetails(
        path=PosixPath("12__101__2020-03-03_01-58-50__Driveway-lowres.mkv"),
        event_id=12,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 58, 50, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=False,
        is_lowres=True,
    ),
    PathDetails(
        path=PosixPath("12__101__2020-03-03_01-58-52__Driveway-lowres.jpg"),
        event_id=12,
        camera_id=101,
        timestamp=datetime.datetime(
            2020, 3, 3, 1, 58, 52, tzinfo=tzoffset("WST-8", 28800)
        ),
        camera_name="Driveway",
        is_image=True,
        is_lowres=True,
    ),
]

_PATH_DETAILS_BY_KEY_IMMEDIATE_MATCH = [
    [
        PathDetails(
            path=PosixPath("01__102__2020-03-03_00-26-28__FrontDoor.mkv"),
            event_id=1,
            camera_id=102,
            timestamp=datetime.datetime(
                2020, 3, 3, 0, 26, 28, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="FrontDoor",
            is_image=False,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("01__102__2020-03-03_00-26-30__FrontDoor.jpg"),
            event_id=1,
            camera_id=102,
            timestamp=datetime.datetime(
                2020, 3, 3, 0, 26, 30, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="FrontDoor",
            is_image=True,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("01__102__2020-03-03_00-26-28__FrontDoor-lowres.mkv"),
            event_id=1,
            camera_id=102,
            timestamp=datetime.datetime(
                2020, 3, 3, 0, 26, 28, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="FrontDoor",
            is_image=False,
            is_lowres=True,
        ),
        PathDetails(
            path=PosixPath("01__102__2020-03-03_00-26-30__FrontDoor-lowres.jpg"),
            event_id=1,
            camera_id=102,
            timestamp=datetime.datetime(
                2020, 3, 3, 0, 26, 30, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="FrontDoor",
            is_image=True,
            is_lowres=True,
        ),
    ],
    [
        PathDetails(
            path=PosixPath("02__102__2020-03-03_00-27-15__FrontDoor.jpg"),
            event_id=2,
            camera_id=102,
            timestamp=datetime.datetime(
                2020, 3, 3, 0, 27, 15, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="FrontDoor",
            is_image=True,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("02__102__2020-03-03_00-27-13__FrontDoor.mkv"),
            event_id=2,
            camera_id=102,
            timestamp=datetime.datetime(
                2020, 3, 3, 0, 27, 13, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="FrontDoor",
            is_image=False,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("02__102__2020-03-03_00-27-13__FrontDoor-lowres.mkv"),
            event_id=2,
            camera_id=102,
            timestamp=datetime.datetime(
                2020, 3, 3, 0, 27, 13, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="FrontDoor",
            is_image=False,
            is_lowres=True,
        ),
        PathDetails(
            path=PosixPath("02__102__2020-03-03_00-27-15__FrontDoor-lowres.jpg"),
            event_id=2,
            camera_id=102,
            timestamp=datetime.datetime(
                2020, 3, 3, 0, 27, 15, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="FrontDoor",
            is_image=True,
            is_lowres=True,
        ),
    ],
    [
        PathDetails(
            path=PosixPath("01__101__2020-03-03_00-43-24__Driveway.jpg"),
            event_id=1,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 0, 43, 24, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("01__101__2020-03-03_00-43-22__Driveway.mkv"),
            event_id=1,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 0, 43, 22, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("01__101__2020-03-03_00-43-22__Driveway-lowres.mkv"),
            event_id=1,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 0, 43, 22, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=True,
        ),
        PathDetails(
            path=PosixPath("01__101__2020-03-03_00-43-24__Driveway-lowres.jpg"),
            event_id=1,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 0, 43, 24, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=True,
        ),
    ],
    [
        PathDetails(
            path=PosixPath("02__101__2020-03-03_00-47-16__Driveway.jpg"),
            event_id=2,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 0, 47, 16, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("02__101__2020-03-03_00-47-14__Driveway.mkv"),
            event_id=2,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 0, 47, 14, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("02__101__2020-03-03_00-47-14__Driveway-lowres.mkv"),
            event_id=2,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 0, 47, 14, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=True,
        ),
        PathDetails(
            path=PosixPath("02__101__2020-03-03_00-47-16__Driveway-lowres.jpg"),
            event_id=2,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 0, 47, 16, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=True,
        ),
    ],
    [
        PathDetails(
            path=PosixPath("03__101__2020-03-03_00-49-52__Driveway.jpg"),
            event_id=3,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 0, 49, 52, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("03__101__2020-03-03_00-49-50__Driveway.mkv"),
            event_id=3,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 0, 49, 50, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("03__101__2020-03-03_00-49-50__Driveway-lowres.mkv"),
            event_id=3,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 0, 49, 50, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=True,
        ),
        PathDetails(
            path=PosixPath("03__101__2020-03-03_00-49-52__Driveway-lowres.jpg"),
            event_id=3,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 0, 49, 52, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=True,
        ),
    ],
    [
        PathDetails(
            path=PosixPath("04__101__2020-03-03_01-00-39__Driveway.jpg"),
            event_id=4,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 0, 39, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("04__101__2020-03-03_01-00-37__Driveway.mkv"),
            event_id=4,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 0, 37, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("04__101__2020-03-03_01-00-37__Driveway-lowres.mkv"),
            event_id=4,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 0, 37, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=True,
        ),
        PathDetails(
            path=PosixPath("04__101__2020-03-03_01-00-39__Driveway-lowres.jpg"),
            event_id=4,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 0, 39, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=True,
        ),
    ],
    [
        PathDetails(
            path=PosixPath("05__101__2020-03-03_01-01-24__Driveway.jpg"),
            event_id=5,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 1, 24, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("05__101__2020-03-03_01-01-23__Driveway.mkv"),
            event_id=5,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 1, 23, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("05__101__2020-03-03_01-01-23__Driveway-lowres.mkv"),
            event_id=5,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 1, 23, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=True,
        ),
        PathDetails(
            path=PosixPath("05__101__2020-03-03_01-01-24__Driveway-lowres.jpg"),
            event_id=5,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 1, 24, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=True,
        ),
    ],
    [
        PathDetails(
            path=PosixPath("03__102__2020-03-03_01-01-51__FrontDoor.mkv"),
            event_id=3,
            camera_id=102,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 1, 51, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="FrontDoor",
            is_image=False,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("03__102__2020-03-03_01-01-53__FrontDoor.jpg"),
            event_id=3,
            camera_id=102,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 1, 53, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="FrontDoor",
            is_image=True,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("03__102__2020-03-03_01-01-51__FrontDoor-lowres.mkv"),
            event_id=3,
            camera_id=102,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 1, 51, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="FrontDoor",
            is_image=False,
            is_lowres=True,
        ),
        PathDetails(
            path=PosixPath("03__102__2020-03-03_01-01-53__FrontDoor-lowres.jpg"),
            event_id=3,
            camera_id=102,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 1, 53, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="FrontDoor",
            is_image=True,
            is_lowres=True,
        ),
    ],
    [
        PathDetails(
            path=PosixPath("06__101__2020-03-03_01-02-03__Driveway.jpg"),
            event_id=6,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 2, 3, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("06__101__2020-03-03_01-02-01__Driveway.mkv"),
            event_id=6,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 2, 1, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("06__101__2020-03-03_01-02-01__Driveway-lowres.mkv"),
            event_id=6,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 2, 1, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=True,
        ),
        PathDetails(
            path=PosixPath("06__101__2020-03-03_01-02-03__Driveway-lowres.jpg"),
            event_id=6,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 2, 3, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=True,
        ),
    ],
    [
        PathDetails(
            path=PosixPath("07__101__2020-03-03_01-02-30__Driveway.jpg"),
            event_id=7,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 2, 30, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("07__101__2020-03-03_01-02-28__Driveway.mkv"),
            event_id=7,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 2, 28, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("07__101__2020-03-03_01-02-28__Driveway-lowres.mkv"),
            event_id=7,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 2, 28, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=True,
        ),
        PathDetails(
            path=PosixPath("07__101__2020-03-03_01-02-30__Driveway-lowres.jpg"),
            event_id=7,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 2, 30, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=True,
        ),
    ],
    [
        PathDetails(
            path=PosixPath("08__101__2020-03-03_01-04-35__Driveway.jpg"),
            event_id=8,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 4, 35, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("08__101__2020-03-03_01-04-34__Driveway.mkv"),
            event_id=8,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 4, 34, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("08__101__2020-03-03_01-04-34__Driveway-lowres.mkv"),
            event_id=8,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 4, 34, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=True,
        ),
        PathDetails(
            path=PosixPath("08__101__2020-03-03_01-04-35__Driveway-lowres.jpg"),
            event_id=8,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 4, 35, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=True,
        ),
    ],
    [
        PathDetails(
            path=PosixPath("04__102__2020-03-03_01-05-05__FrontDoor.mkv"),
            event_id=4,
            camera_id=102,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 5, 5, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="FrontDoor",
            is_image=False,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("04__102__2020-03-03_01-05-07__FrontDoor.jpg"),
            event_id=4,
            camera_id=102,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 5, 7, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="FrontDoor",
            is_image=True,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("04__102__2020-03-03_01-05-05__FrontDoor-lowres.mkv"),
            event_id=4,
            camera_id=102,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 5, 5, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="FrontDoor",
            is_image=False,
            is_lowres=True,
        ),
        PathDetails(
            path=PosixPath("04__102__2020-03-03_01-05-07__FrontDoor-lowres.jpg"),
            event_id=4,
            camera_id=102,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 5, 7, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="FrontDoor",
            is_image=True,
            is_lowres=True,
        ),
    ],
    [
        PathDetails(
            path=PosixPath("09__101__2020-03-03_01-27-27__Driveway.jpg"),
            event_id=9,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 27, 27, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("09__101__2020-03-03_01-27-25__Driveway.mkv"),
            event_id=9,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 27, 25, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("09__101__2020-03-03_01-27-25__Driveway-lowres.mkv"),
            event_id=9,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 27, 25, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=True,
        ),
        PathDetails(
            path=PosixPath("09__101__2020-03-03_01-27-27__Driveway-lowres.jpg"),
            event_id=9,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 27, 27, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=True,
        ),
    ],
    [
        PathDetails(
            path=PosixPath("10__101__2020-03-03_01-48-27__Driveway.jpg"),
            event_id=10,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 48, 27, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("10__101__2020-03-03_01-48-26__Driveway.mkv"),
            event_id=10,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 48, 26, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("10__101__2020-03-03_01-48-26__Driveway-lowres.mkv"),
            event_id=10,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 48, 26, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=True,
        ),
        PathDetails(
            path=PosixPath("10__101__2020-03-03_01-48-27__Driveway-lowres.jpg"),
            event_id=10,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 48, 27, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=True,
        ),
    ],
    [
        PathDetails(
            path=PosixPath("11__101__2020-03-03_01-54-19__Driveway.jpg"),
            event_id=11,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 54, 19, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("11__101__2020-03-03_01-54-17__Driveway.mkv"),
            event_id=11,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 54, 17, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("11__101__2020-03-03_01-54-17__Driveway-lowres.mkv"),
            event_id=11,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 54, 17, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=True,
        ),
        PathDetails(
            path=PosixPath("11__101__2020-03-03_01-54-19__Driveway-lowres.jpg"),
            event_id=11,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 54, 19, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=True,
        ),
    ],
    [
        PathDetails(
            path=PosixPath("12__101__2020-03-03_01-58-52__Driveway.jpg"),
            event_id=12,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 58, 52, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("12__101__2020-03-03_01-58-50__Driveway.mkv"),
            event_id=12,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 58, 50, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=False,
        ),
        PathDetails(
            path=PosixPath("12__101__2020-03-03_01-58-50__Driveway-lowres.mkv"),
            event_id=12,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 58, 50, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=False,
            is_lowres=True,
        ),
        PathDetails(
            path=PosixPath("12__101__2020-03-03_01-58-52__Driveway-lowres.jpg"),
            event_id=12,
            camera_id=101,
            timestamp=datetime.datetime(
                2020, 3, 3, 1, 58, 52, tzinfo=tzoffset("WST-8", 28800)
            ),
            camera_name="Driveway",
            is_image=True,
            is_lowres=True,
        ),
    ],
]

_EVENT_FROM_PATH_DETAILS_BY_KEY_IMMEDIATE_MATCH_FIRST_RECORD = Event(
    event_id="028a7569-549a-396c-9a95-ed4a2d8118f1",
    timestamp=datetime.datetime(2020, 3, 3, 0, 26, 28, tzinfo=tzoffset("WST-8", 28800)),
    camera_name="FrontDoor",
    high_res_image_path="/srv/target_dir/events/01__102__2020-03-03_00-26-30__FrontDoor.jpg",
    low_res_image_path="/srv/target_dir/events/01__102__2020-03-03_00-26-30__FrontDoor-lowres.jpg",
    high_res_video_path="/srv/target_dir/events/01__102__2020-03-03_00-26-28__FrontDoor.mkv",
    low_res_video_path="/srv/target_dir/events/01__102__2020-03-03_00-26-28__FrontDoor-lowres.mkv",
)

_SORTED_EVENTS = [
    Event(
        camera_name="FrontDoor",
        event_id="028a7569-549a-396c-9a95-ed4a2d8118f1",
        high_res_image_path="/srv/target_dir/events/01__102__2020-03-03_00-26-30__FrontDoor.jpg",
        high_res_video_path="/srv/target_dir/events/01__102__2020-03-03_00-26-28__FrontDoor.mkv",
        low_res_image_path="/srv/target_dir/events/01__102__2020-03-03_00-26-30__FrontDoor-lowres.jpg",
        low_res_video_path="/srv/target_dir/events/01__102__2020-03-03_00-26-28__FrontDoor-lowres.mkv",
        timestamp="2020-03-03T00:26:28.00000000+08:00",
    ),
    Event(
        camera_name="FrontDoor",
        event_id="b3517c2d-a987-3296-94bc-9ab25c223129",
        high_res_image_path="/srv/target_dir/events/02__102__2020-03-03_00-27-15__FrontDoor.jpg",
        high_res_video_path="/srv/target_dir/events/02__102__2020-03-03_00-27-13__FrontDoor.mkv",
        low_res_image_path="/srv/target_dir/events/02__102__2020-03-03_00-27-15__FrontDoor-lowres.jpg",
        low_res_video_path="/srv/target_dir/events/02__102__2020-03-03_00-27-13__FrontDoor-lowres.mkv",
        timestamp="2020-03-03T00:27:13.00000000+08:00",
    ),
    Event(
        camera_name="Driveway",
        event_id="bfbcf1a9-00a1-317f-8901-ff140725fb2c",
        high_res_image_path="/srv/target_dir/events/01__101__2020-03-03_00-43-24__Driveway.jpg",
        high_res_video_path="/srv/target_dir/events/01__101__2020-03-03_00-43-22__Driveway.mkv",
        low_res_image_path="/srv/target_dir/events/01__101__2020-03-03_00-43-24__Driveway-lowres.jpg",
        low_res_video_path="/srv/target_dir/events/01__101__2020-03-03_00-43-22__Driveway-lowres.mkv",
        timestamp="2020-03-03T00:43:22.00000000+08:00",
    ),
    Event(
        camera_name="Driveway",
        event_id="cced7bca-a745-3f74-95b1-8654eaf95878",
        high_res_image_path="/srv/target_dir/events/02__101__2020-03-03_00-47-16__Driveway.jpg",
        high_res_video_path="/srv/target_dir/events/02__101__2020-03-03_00-47-14__Driveway.mkv",
        low_res_image_path="/srv/target_dir/events/02__101__2020-03-03_00-47-16__Driveway-lowres.jpg",
        low_res_video_path="/srv/target_dir/events/02__101__2020-03-03_00-47-14__Driveway-lowres.mkv",
        timestamp="2020-03-03T00:47:14.00000000+08:00",
    ),
    Event(
        camera_name="Driveway",
        event_id="ee4eb6f2-710d-3e02-aee0-0ba8fdb423b4",
        high_res_image_path="/srv/target_dir/events/03__101__2020-03-03_00-49-52__Driveway.jpg",
        high_res_video_path="/srv/target_dir/events/03__101__2020-03-03_00-49-50__Driveway.mkv",
        low_res_image_path="/srv/target_dir/events/03__101__2020-03-03_00-49-52__Driveway-lowres.jpg",
        low_res_video_path="/srv/target_dir/events/03__101__2020-03-03_00-49-50__Driveway-lowres.mkv",
        timestamp="2020-03-03T00:49:50.00000000+08:00",
    ),
    Event(
        camera_name="Driveway",
        event_id="8d799eeb-542b-372a-beac-b6ea9727410e",
        high_res_image_path="/srv/target_dir/events/04__101__2020-03-03_01-00-39__Driveway.jpg",
        high_res_video_path="/srv/target_dir/events/04__101__2020-03-03_01-00-37__Driveway.mkv",
        low_res_image_path="/srv/target_dir/events/04__101__2020-03-03_01-00-39__Driveway-lowres.jpg",
        low_res_video_path="/srv/target_dir/events/04__101__2020-03-03_01-00-37__Driveway-lowres.mkv",
        timestamp="2020-03-03T01:00:37.00000000+08:00",
    ),
    Event(
        camera_name="Driveway",
        event_id="40198285-b26c-31e0-b51c-647ad52f6a83",
        high_res_image_path="/srv/target_dir/events/05__101__2020-03-03_01-01-24__Driveway.jpg",
        high_res_video_path="/srv/target_dir/events/05__101__2020-03-03_01-01-23__Driveway.mkv",
        low_res_image_path="/srv/target_dir/events/05__101__2020-03-03_01-01-24__Driveway-lowres.jpg",
        low_res_video_path="/srv/target_dir/events/05__101__2020-03-03_01-01-23__Driveway-lowres.mkv",
        timestamp="2020-03-03T01:01:23.00000000+08:00",
    ),
    Event(
        camera_name="FrontDoor",
        event_id="be1e1c6b-be6a-3db4-81df-2774f61bc9f2",
        high_res_image_path="/srv/target_dir/events/03__102__2020-03-03_01-01-53__FrontDoor.jpg",
        high_res_video_path="/srv/target_dir/events/03__102__2020-03-03_01-01-51__FrontDoor.mkv",
        low_res_image_path="/srv/target_dir/events/03__102__2020-03-03_01-01-53__FrontDoor-lowres.jpg",
        low_res_video_path="/srv/target_dir/events/03__102__2020-03-03_01-01-51__FrontDoor-lowres.mkv",
        timestamp="2020-03-03T01:01:51.00000000+08:00",
    ),
    Event(
        camera_name="Driveway",
        event_id="437aa95d-ac50-3626-be8f-e5498e17d7b7",
        high_res_image_path="/srv/target_dir/events/06__101__2020-03-03_01-02-03__Driveway.jpg",
        high_res_video_path="/srv/target_dir/events/06__101__2020-03-03_01-02-01__Driveway.mkv",
        low_res_image_path="/srv/target_dir/events/06__101__2020-03-03_01-02-03__Driveway-lowres.jpg",
        low_res_video_path="/srv/target_dir/events/06__101__2020-03-03_01-02-01__Driveway-lowres.mkv",
        timestamp="2020-03-03T01:02:01.00000000+08:00",
    ),
    Event(
        camera_name="Driveway",
        event_id="c1ad867b-0a43-332b-8eb7-adf6c025fed5",
        high_res_image_path="/srv/target_dir/events/07__101__2020-03-03_01-02-30__Driveway.jpg",
        high_res_video_path="/srv/target_dir/events/07__101__2020-03-03_01-02-28__Driveway.mkv",
        low_res_image_path="/srv/target_dir/events/07__101__2020-03-03_01-02-30__Driveway-lowres.jpg",
        low_res_video_path="/srv/target_dir/events/07__101__2020-03-03_01-02-28__Driveway-lowres.mkv",
        timestamp="2020-03-03T01:02:28.00000000+08:00",
    ),
    Event(
        camera_name="Driveway",
        event_id="ca77ae39-92e6-3d1a-8ff3-9434b889aa15",
        high_res_image_path="/srv/target_dir/events/08__101__2020-03-03_01-04-35__Driveway.jpg",
        high_res_video_path="/srv/target_dir/events/08__101__2020-03-03_01-04-34__Driveway.mkv",
        low_res_image_path="/srv/target_dir/events/08__101__2020-03-03_01-04-35__Driveway-lowres.jpg",
        low_res_video_path="/srv/target_dir/events/08__101__2020-03-03_01-04-34__Driveway-lowres.mkv",
        timestamp="2020-03-03T01:04:34.00000000+08:00",
    ),
    Event(
        camera_name="FrontDoor",
        event_id="1ab2183e-d859-3821-9081-0927469d55c2",
        high_res_image_path="/srv/target_dir/events/04__102__2020-03-03_01-05-07__FrontDoor.jpg",
        high_res_video_path="/srv/target_dir/events/04__102__2020-03-03_01-05-05__FrontDoor.mkv",
        low_res_image_path="/srv/target_dir/events/04__102__2020-03-03_01-05-07__FrontDoor-lowres.jpg",
        low_res_video_path="/srv/target_dir/events/04__102__2020-03-03_01-05-05__FrontDoor-lowres.mkv",
        timestamp="2020-03-03T01:05:05.00000000+08:00",
    ),
    Event(
        camera_name="Driveway",
        event_id="eefba6ae-7be3-3158-a94f-59a5c191c385",
        high_res_image_path="/srv/target_dir/events/09__101__2020-03-03_01-27-27__Driveway.jpg",
        high_res_video_path="/srv/target_dir/events/09__101__2020-03-03_01-27-25__Driveway.mkv",
        low_res_image_path="/srv/target_dir/events/09__101__2020-03-03_01-27-27__Driveway-lowres.jpg",
        low_res_video_path="/srv/target_dir/events/09__101__2020-03-03_01-27-25__Driveway-lowres.mkv",
        timestamp="2020-03-03T01:27:25.00000000+08:00",
    ),
    Event(
        camera_name="Driveway",
        event_id="4966768d-cec7-3c88-a183-d47c2093d5eb",
        high_res_image_path="/srv/target_dir/events/10__101__2020-03-03_01-48-27__Driveway.jpg",
        high_res_video_path="/srv/target_dir/events/10__101__2020-03-03_01-48-26__Driveway.mkv",
        low_res_image_path="/srv/target_dir/events/10__101__2020-03-03_01-48-27__Driveway-lowres.jpg",
        low_res_video_path="/srv/target_dir/events/10__101__2020-03-03_01-48-26__Driveway-lowres.mkv",
        timestamp="2020-03-03T01:48:26.00000000+08:00",
    ),
    Event(
        camera_name="Driveway",
        event_id="ff1a771c-720f-3d51-af87-f557996934a0",
        high_res_image_path="/srv/target_dir/events/11__101__2020-03-03_01-54-19__Driveway.jpg",
        high_res_video_path="/srv/target_dir/events/11__101__2020-03-03_01-54-17__Driveway.mkv",
        low_res_image_path="/srv/target_dir/events/11__101__2020-03-03_01-54-19__Driveway-lowres.jpg",
        low_res_video_path="/srv/target_dir/events/11__101__2020-03-03_01-54-17__Driveway-lowres.mkv",
        timestamp="2020-03-03T01:54:17.00000000+08:00",
    ),
    Event(
        camera_name="Driveway",
        event_id="8fd397f2-8dd4-36d0-8d54-bcbb09daf7c4",
        high_res_image_path="/srv/target_dir/events/12__101__2020-03-03_01-58-52__Driveway.jpg",
        high_res_video_path="/srv/target_dir/events/12__101__2020-03-03_01-58-50__Driveway.mkv",
        low_res_image_path="/srv/target_dir/events/12__101__2020-03-03_01-58-52__Driveway-lowres.jpg",
        low_res_video_path="/srv/target_dir/events/12__101__2020-03-03_01-58-50__Driveway-lowres.mkv",
        timestamp="2020-03-03T01:58:50.00000000+08:00",
    ),
]

_JSON_LINES = """{"event_id": "028a7569-549a-396c-9a95-ed4a2d8118f1", "timestamp": "2020-03-03T00:26:28.00000000+08:00", "camera_name": "FrontDoor", "high_res_image_path": "/srv/target_dir/events/01__102__2020-03-03_00-26-30__FrontDoor.jpg", "low_res_image_path": "/srv/target_dir/events/01__102__2020-03-03_00-26-30__FrontDoor-lowres.jpg", "high_res_video_path": "/srv/target_dir/events/01__102__2020-03-03_00-26-28__FrontDoor.mkv", "low_res_video_path": "/srv/target_dir/events/01__102__2020-03-03_00-26-28__FrontDoor-lowres.mkv"}
{"event_id": "b3517c2d-a987-3296-94bc-9ab25c223129", "timestamp": "2020-03-03T00:27:13.00000000+08:00", "camera_name": "FrontDoor", "high_res_image_path": "/srv/target_dir/events/02__102__2020-03-03_00-27-15__FrontDoor.jpg", "low_res_image_path": "/srv/target_dir/events/02__102__2020-03-03_00-27-15__FrontDoor-lowres.jpg", "high_res_video_path": "/srv/target_dir/events/02__102__2020-03-03_00-27-13__FrontDoor.mkv", "low_res_video_path": "/srv/target_dir/events/02__102__2020-03-03_00-27-13__FrontDoor-lowres.mkv"}
{"event_id": "bfbcf1a9-00a1-317f-8901-ff140725fb2c", "timestamp": "2020-03-03T00:43:22.00000000+08:00", "camera_name": "Driveway", "high_res_image_path": "/srv/target_dir/events/01__101__2020-03-03_00-43-24__Driveway.jpg", "low_res_image_path": "/srv/target_dir/events/01__101__2020-03-03_00-43-24__Driveway-lowres.jpg", "high_res_video_path": "/srv/target_dir/events/01__101__2020-03-03_00-43-22__Driveway.mkv", "low_res_video_path": "/srv/target_dir/events/01__101__2020-03-03_00-43-22__Driveway-lowres.mkv"}
{"event_id": "cced7bca-a745-3f74-95b1-8654eaf95878", "timestamp": "2020-03-03T00:47:14.00000000+08:00", "camera_name": "Driveway", "high_res_image_path": "/srv/target_dir/events/02__101__2020-03-03_00-47-16__Driveway.jpg", "low_res_image_path": "/srv/target_dir/events/02__101__2020-03-03_00-47-16__Driveway-lowres.jpg", "high_res_video_path": "/srv/target_dir/events/02__101__2020-03-03_00-47-14__Driveway.mkv", "low_res_video_path": "/srv/target_dir/events/02__101__2020-03-03_00-47-14__Driveway-lowres.mkv"}
{"event_id": "ee4eb6f2-710d-3e02-aee0-0ba8fdb423b4", "timestamp": "2020-03-03T00:49:50.00000000+08:00", "camera_name": "Driveway", "high_res_image_path": "/srv/target_dir/events/03__101__2020-03-03_00-49-52__Driveway.jpg", "low_res_image_path": "/srv/target_dir/events/03__101__2020-03-03_00-49-52__Driveway-lowres.jpg", "high_res_video_path": "/srv/target_dir/events/03__101__2020-03-03_00-49-50__Driveway.mkv", "low_res_video_path": "/srv/target_dir/events/03__101__2020-03-03_00-49-50__Driveway-lowres.mkv"}
{"event_id": "8d799eeb-542b-372a-beac-b6ea9727410e", "timestamp": "2020-03-03T01:00:37.00000000+08:00", "camera_name": "Driveway", "high_res_image_path": "/srv/target_dir/events/04__101__2020-03-03_01-00-39__Driveway.jpg", "low_res_image_path": "/srv/target_dir/events/04__101__2020-03-03_01-00-39__Driveway-lowres.jpg", "high_res_video_path": "/srv/target_dir/events/04__101__2020-03-03_01-00-37__Driveway.mkv", "low_res_video_path": "/srv/target_dir/events/04__101__2020-03-03_01-00-37__Driveway-lowres.mkv"}
{"event_id": "40198285-b26c-31e0-b51c-647ad52f6a83", "timestamp": "2020-03-03T01:01:23.00000000+08:00", "camera_name": "Driveway", "high_res_image_path": "/srv/target_dir/events/05__101__2020-03-03_01-01-24__Driveway.jpg", "low_res_image_path": "/srv/target_dir/events/05__101__2020-03-03_01-01-24__Driveway-lowres.jpg", "high_res_video_path": "/srv/target_dir/events/05__101__2020-03-03_01-01-23__Driveway.mkv", "low_res_video_path": "/srv/target_dir/events/05__101__2020-03-03_01-01-23__Driveway-lowres.mkv"}
{"event_id": "be1e1c6b-be6a-3db4-81df-2774f61bc9f2", "timestamp": "2020-03-03T01:01:51.00000000+08:00", "camera_name": "FrontDoor", "high_res_image_path": "/srv/target_dir/events/03__102__2020-03-03_01-01-53__FrontDoor.jpg", "low_res_image_path": "/srv/target_dir/events/03__102__2020-03-03_01-01-53__FrontDoor-lowres.jpg", "high_res_video_path": "/srv/target_dir/events/03__102__2020-03-03_01-01-51__FrontDoor.mkv", "low_res_video_path": "/srv/target_dir/events/03__102__2020-03-03_01-01-51__FrontDoor-lowres.mkv"}
{"event_id": "437aa95d-ac50-3626-be8f-e5498e17d7b7", "timestamp": "2020-03-03T01:02:01.00000000+08:00", "camera_name": "Driveway", "high_res_image_path": "/srv/target_dir/events/06__101__2020-03-03_01-02-03__Driveway.jpg", "low_res_image_path": "/srv/target_dir/events/06__101__2020-03-03_01-02-03__Driveway-lowres.jpg", "high_res_video_path": "/srv/target_dir/events/06__101__2020-03-03_01-02-01__Driveway.mkv", "low_res_video_path": "/srv/target_dir/events/06__101__2020-03-03_01-02-01__Driveway-lowres.mkv"}
{"event_id": "c1ad867b-0a43-332b-8eb7-adf6c025fed5", "timestamp": "2020-03-03T01:02:28.00000000+08:00", "camera_name": "Driveway", "high_res_image_path": "/srv/target_dir/events/07__101__2020-03-03_01-02-30__Driveway.jpg", "low_res_image_path": "/srv/target_dir/events/07__101__2020-03-03_01-02-30__Driveway-lowres.jpg", "high_res_video_path": "/srv/target_dir/events/07__101__2020-03-03_01-02-28__Driveway.mkv", "low_res_video_path": "/srv/target_dir/events/07__101__2020-03-03_01-02-28__Driveway-lowres.mkv"}
{"event_id": "ca77ae39-92e6-3d1a-8ff3-9434b889aa15", "timestamp": "2020-03-03T01:04:34.00000000+08:00", "camera_name": "Driveway", "high_res_image_path": "/srv/target_dir/events/08__101__2020-03-03_01-04-35__Driveway.jpg", "low_res_image_path": "/srv/target_dir/events/08__101__2020-03-03_01-04-35__Driveway-lowres.jpg", "high_res_video_path": "/srv/target_dir/events/08__101__2020-03-03_01-04-34__Driveway.mkv", "low_res_video_path": "/srv/target_dir/events/08__101__2020-03-03_01-04-34__Driveway-lowres.mkv"}
{"event_id": "1ab2183e-d859-3821-9081-0927469d55c2", "timestamp": "2020-03-03T01:05:05.00000000+08:00", "camera_name": "FrontDoor", "high_res_image_path": "/srv/target_dir/events/04__102__2020-03-03_01-05-07__FrontDoor.jpg", "low_res_image_path": "/srv/target_dir/events/04__102__2020-03-03_01-05-07__FrontDoor-lowres.jpg", "high_res_video_path": "/srv/target_dir/events/04__102__2020-03-03_01-05-05__FrontDoor.mkv", "low_res_video_path": "/srv/target_dir/events/04__102__2020-03-03_01-05-05__FrontDoor-lowres.mkv"}
{"event_id": "eefba6ae-7be3-3158-a94f-59a5c191c385", "timestamp": "2020-03-03T01:27:25.00000000+08:00", "camera_name": "Driveway", "high_res_image_path": "/srv/target_dir/events/09__101__2020-03-03_01-27-27__Driveway.jpg", "low_res_image_path": "/srv/target_dir/events/09__101__2020-03-03_01-27-27__Driveway-lowres.jpg", "high_res_video_path": "/srv/target_dir/events/09__101__2020-03-03_01-27-25__Driveway.mkv", "low_res_video_path": "/srv/target_dir/events/09__101__2020-03-03_01-27-25__Driveway-lowres.mkv"}
{"event_id": "4966768d-cec7-3c88-a183-d47c2093d5eb", "timestamp": "2020-03-03T01:48:26.00000000+08:00", "camera_name": "Driveway", "high_res_image_path": "/srv/target_dir/events/10__101__2020-03-03_01-48-27__Driveway.jpg", "low_res_image_path": "/srv/target_dir/events/10__101__2020-03-03_01-48-27__Driveway-lowres.jpg", "high_res_video_path": "/srv/target_dir/events/10__101__2020-03-03_01-48-26__Driveway.mkv", "low_res_video_path": "/srv/target_dir/events/10__101__2020-03-03_01-48-26__Driveway-lowres.mkv"}
{"event_id": "ff1a771c-720f-3d51-af87-f557996934a0", "timestamp": "2020-03-03T01:54:17.00000000+08:00", "camera_name": "Driveway", "high_res_image_path": "/srv/target_dir/events/11__101__2020-03-03_01-54-19__Driveway.jpg", "low_res_image_path": "/srv/target_dir/events/11__101__2020-03-03_01-54-19__Driveway-lowres.jpg", "high_res_video_path": "/srv/target_dir/events/11__101__2020-03-03_01-54-17__Driveway.mkv", "low_res_video_path": "/srv/target_dir/events/11__101__2020-03-03_01-54-17__Driveway-lowres.mkv"}
{"event_id": "8fd397f2-8dd4-36d0-8d54-bcbb09daf7c4", "timestamp": "2020-03-03T01:58:50.00000000+08:00", "camera_name": "Driveway", "high_res_image_path": "/srv/target_dir/events/12__101__2020-03-03_01-58-52__Driveway.jpg", "low_res_image_path": "/srv/target_dir/events/12__101__2020-03-03_01-58-52__Driveway-lowres.jpg", "high_res_video_path": "/srv/target_dir/events/12__101__2020-03-03_01-58-50__Driveway.mkv", "low_res_video_path": "/srv/target_dir/events/12__101__2020-03-03_01-58-50__Driveway-lowres.mkv"}"""


class EventStoreRebuilderTest(unittest.TestCase):
    maxDiff = 65536 * 8

    def test_get_sorted_paths(self):
        tmpdir = Path(mkdtemp())

        paths = []
        for i in range(0, 10):
            path = tmpdir / Path(f"file_{i}.txt")
            path.touch()
            paths += [path]

        self.assertEqual(paths, get_sorted_paths(tmpdir))

    def test_parse_path(self):
        self.assertEqual(_PATH_DETAILS, parse_path(path=_SOME_PATHS[3], tzinfo=_TZINFO))

    def test_parse_paths(self):
        # print(parse_paths(paths=_SOME_PATHS, tzinfo=_TZINFO))

        self.assertEqual(
            _SOME_PATH_DETAILS, parse_paths(paths=_SOME_PATHS, tzinfo=_TZINFO)
        )

    def test_relate_path_details(self):
        # print(relate_path_details(_SOME_PATH_DETAILS))

        self.assertEqual(
            _PATH_DETAILS_BY_KEY_IMMEDIATE_MATCH,
            relate_path_details(_SOME_PATH_DETAILS),
        )

    def test_format_timestamp_for_go(self):
        self.assertEqual(
            _TEST_TIMESTAMP_IN_GO_FORMAT, format_timestamp_for_go(_TEST_TIMESTAMP)
        )

    def test_build_event_for_some_path_details_from_path_details_by_key_immediate_match_first_record(
        self,
    ):
        self.assertEqual(
            _EVENT_FROM_PATH_DETAILS_BY_KEY_IMMEDIATE_MATCH_FIRST_RECORD,
            build_event_for_some_path_details(
                some_path_details=_PATH_DETAILS_BY_KEY_IMMEDIATE_MATCH[0],
                path=PosixPath("/srv/target_dir/events/"),
            ),
        )

    def test_build_events_for_path_details_by_key(self):
        # print(
        #     build_events_for_path_details_by_key(
        #         some_path_details_by_key=_PATH_DETAILS_BY_KEY_IMMEDIATE_MATCH,
        #         path=PosixPath("/srv/target_dir/events"),
        #     )
        # )

        self.assertEqual(
            _SORTED_EVENTS,
            build_events_for_related_path_details(
                related_path_details=_PATH_DETAILS_BY_KEY_IMMEDIATE_MATCH,
                path=PosixPath("/srv/target_dir/events"),
            ),
        )

    def test_build_json_lines_from_events(self):
        # print(build_json_lines_from_events(_SORTED_EVENTS))

        self.assertEqual(_JSON_LINES, build_json_lines_from_events(_SORTED_EVENTS))

    def test_rebuild_event_store(self):
        tmpdir = Path(mkdtemp())

        for path in _SOME_PATHS:
            full_path = tmpdir / path
            full_path.touch()

        json_path = tmpdir / Path("events.jsonl")

        rebuild_event_store(root_path=tmpdir, tzinfo=_TZINFO, json_path=json_path)

        with open(str(json_path), "r") as f:
            data = f.read().strip()

        self.assertEqual(_JSON_LINES.count("\n"), data.count("\n"))

        for path in _SORTED_EVENTS:
            self.assertIn(path.high_res_image_path.split("/")[-1], data)
            self.assertIn(path.low_res_image_path.split("/")[-1], data)
            self.assertIn(path.high_res_video_path.split("/")[-1], data)
            self.assertIn(path.low_res_video_path.split("/")[-1], data)

    def test_many_paths_with_poor_data_sanity(self):
        sorted_paths = _sorted_paths_events
        tzinfo = _TZINFO

        some_path_details = parse_paths(paths=sorted_paths, tzinfo=tzinfo)
        some_path_details_by_key = relate_path_details(
            some_path_details=some_path_details
        )

        self.assertEqual(5275, len(some_path_details_by_key))

import datetime
import unittest
from pathlib import Path, PosixPath
from tempfile import mkdtemp

from dateutil.tz import tzoffset

from .common import (
    parse_paths,
    relate_path_details,
)
from .common_test import _TZINFO, _PATH_DETAILS_BY_KEY_IMMEDIATE_MATCH, _SORTED_EVENTS, _JSON_LINES
from .event_store_rebuilder_for_events import (
    PathDetails,
    parse_path,
    rebuild_event_store,
    _get_key_pass_1,
    _get_key_pass_2,
    _get_key_pass_3
)
from .test_data_for_events import _sorted_paths_events

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


class EventStoreRebuilderForEventsTest(unittest.TestCase):
    maxDiff = 65536 * 8

    def test_parse_path(self):
        self.assertEqual(_PATH_DETAILS, parse_path(path=_SOME_PATHS[3], tzinfo=_TZINFO))

    def test_parse_paths(self):
        # print(parse_paths(paths=_SOME_PATHS, tzinfo=_TZINFO))

        self.assertEqual(
            _SOME_PATH_DETAILS, parse_paths(paths=_SOME_PATHS, tzinfo=_TZINFO, parse_method=parse_path)
        )

    def test_relate_path_details(self):
        # print(relate_path_details(_SOME_PATH_DETAILS))

        self.assertEqual(
            _PATH_DETAILS_BY_KEY_IMMEDIATE_MATCH,
            relate_path_details(_SOME_PATH_DETAILS, get_key_methods=[
                _get_key_pass_1, _get_key_pass_2, _get_key_pass_3
            ]),
        )

    def test_rebuild_event_store(self):
        tmpdir = Path(mkdtemp())

        for path in _SOME_PATHS:
            full_path = tmpdir / path
            full_path.touch()

        json_path = tmpdir / Path("events.jsonl")

        rebuild_event_store(root_path=tmpdir, tzinfo=_TZINFO, json_path=json_path, parse_method=parse_path,
                            get_key_methods=[_get_key_pass_1, _get_key_pass_2, _get_key_pass_3])

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

        some_path_details = parse_paths(paths=sorted_paths, tzinfo=tzinfo, parse_method=parse_path)
        some_path_details_by_key = relate_path_details(
            some_path_details=some_path_details, get_key_methods=[
                _get_key_pass_1, _get_key_pass_2, _get_key_pass_3
            ]
        )

        self.assertEqual(5275, len(some_path_details_by_key))

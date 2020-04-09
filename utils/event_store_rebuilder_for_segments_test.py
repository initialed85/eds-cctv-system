import datetime
import unittest
from pathlib import PosixPath, Path
from tempfile import mkdtemp

from dateutil.tz import tzoffset

from .common import parse_paths, relate_path_details, rebuild_event_store
from .common_test import _TZINFO, PathDetails
from .event_store_rebuilder_for_segments import parse_path, _get_key
from .test_data_for_segments import _sorted_paths_segments

_SOME_PATHS = _sorted_paths_segments[-64:]

_PATH_DETAILS = PathDetails(path=PosixPath('Segment_2020-04-10_03-20-01_SideGate.mp4'), event_id=None,
                            camera_id=None, timestamp=datetime.datetime(2020, 4, 10, 3, 20, 1, tzinfo=tzoffset('WST-8', 28800)),
                            camera_name='SideGate', is_image=False, is_lowres=False)

_SOME_PATH_DETAILS = [
    PathDetails(path=PosixPath('Segment_2020-04-10_03-15-01_Driveway-lowres.jpg'), event_id=None,
                camera_id=None, timestamp=datetime.datetime(2020, 4, 10, 3, 15, 1, tzinfo=tzoffset('WST-8', 28800)), camera_name='Driveway',
                is_image=True, is_lowres=True),
    PathDetails(path=PosixPath('Segment_2020-04-10_03-15-01_FrontDoor-lowres.mp4'), event_id=None,
                camera_id=None, timestamp=datetime.datetime(2020, 4, 10, 3, 15, 1, tzinfo=tzoffset('WST-8', 28800)),
                camera_name='FrontDoor', is_image=False, is_lowres=True),
    PathDetails(path=PosixPath('Segment_2020-04-10_03-15-01_FrontDoor-lowres.jpg'), event_id=None,
                camera_id=None, timestamp=datetime.datetime(2020, 4, 10, 3, 15, 1, tzinfo=tzoffset('WST-8', 28800)),
                camera_name='FrontDoor', is_image=True, is_lowres=True),
    PathDetails(path=PosixPath('Segment_2020-04-10_03-20-01_SideGate.mp4'), event_id=None, camera_id=None,
                timestamp=datetime.datetime(2020, 4, 10, 3, 20, 1, tzinfo=tzoffset('WST-8', 28800)), camera_name='SideGate', is_image=False,
                is_lowres=False),
    PathDetails(path=PosixPath('Segment_2020-04-10_03-20-01_FrontDoor.mp4'), event_id=None, camera_id=None,
                timestamp=datetime.datetime(2020, 4, 10, 3, 20, 1, tzinfo=tzoffset('WST-8', 28800)), camera_name='FrontDoor',
                is_image=False, is_lowres=False),
    PathDetails(path=PosixPath('Segment_2020-04-10_03-20-01_Driveway.mp4'), event_id=None, camera_id=None,
                timestamp=datetime.datetime(2020, 4, 10, 3, 20, 1, tzinfo=tzoffset('WST-8', 28800)), camera_name='Driveway', is_image=False,
                is_lowres=False),
    PathDetails(path=PosixPath('Segment_2020-04-10_03-24-32_SideGate.mp4'), event_id=None, camera_id=None,
                timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)), camera_name='SideGate',
                is_image=False, is_lowres=False),
    PathDetails(path=PosixPath('Segment_2020-04-10_03-24-32_Driveway.mp4'), event_id=None, camera_id=None,
                timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)), camera_name='Driveway',
                is_image=False, is_lowres=False),
    PathDetails(path=PosixPath('Segment_2020-04-10_03-24-32_FrontDoor.mp4'), event_id=None, camera_id=None,
                timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)), camera_name='FrontDoor',
                is_image=False, is_lowres=False),
    PathDetails(path=PosixPath('Segment_2020-04-10_03-24-32_SideGate.jpg'), event_id=None, camera_id=None,
                timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)), camera_name='SideGate', is_image=True,
                is_lowres=False),
    PathDetails(path=PosixPath('Segment_2020-04-10_03-24-32_FrontDoor.jpg'), event_id=None, camera_id=None,
                timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)), camera_name='FrontDoor',
                is_image=True, is_lowres=False),
    PathDetails(path=PosixPath('Segment_2020-04-10_03-24-32_Driveway.jpg'), event_id=None, camera_id=None,
                timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)), camera_name='Driveway', is_image=True,
                is_lowres=False),
    PathDetails(path=PosixPath('Segment_2020-04-10_03-24-32_SideGate-lowres.mp4'), event_id=None,
                camera_id=None, timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)),
                camera_name='SideGate', is_image=False, is_lowres=True),
    PathDetails(path=PosixPath('Segment_2020-04-10_03-24-32_SideGate-lowres.jpg'), event_id=None,
                camera_id=None, timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)),
                camera_name='SideGate', is_image=True, is_lowres=True),
    PathDetails(path=PosixPath('Segment_2020-04-10_03-24-32_FrontDoor-lowres.mp4'), event_id=None,
                camera_id=None, timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)),
                camera_name='FrontDoor', is_image=False, is_lowres=True),
    PathDetails(path=PosixPath('Segment_2020-04-10_03-24-32_FrontDoor-lowres.jpg'), event_id=None,
                camera_id=None, timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)),
                camera_name='FrontDoor', is_image=True, is_lowres=True),
    PathDetails(path=PosixPath('Segment_2020-04-10_03-24-32_Driveway-lowres.mp4'), event_id=None,
                camera_id=None, timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)),
                camera_name='Driveway', is_image=False, is_lowres=True),
    PathDetails(path=PosixPath('Segment_2020-04-10_03-24-32_Driveway-lowres.jpg'), event_id=None,
                camera_id=None, timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)),
                camera_name='Driveway', is_image=True, is_lowres=True),
    PathDetails(path=PosixPath('Segment_2020-04-10_03-25-00_SideGate.mp4'), event_id=None, camera_id=None,
                timestamp=datetime.datetime(2020, 4, 10, 3, 25, tzinfo=tzoffset('WST-8', 28800)), camera_name='SideGate', is_image=False,
                is_lowres=False),
    PathDetails(path=PosixPath('Segment_2020-04-10_03-25-00_Driveway.mp4'), event_id=None, camera_id=None,
                timestamp=datetime.datetime(2020, 4, 10, 3, 25, tzinfo=tzoffset('WST-8', 28800)), camera_name='Driveway', is_image=False,
                is_lowres=False),
    PathDetails(path=PosixPath('Segment_2020-04-10_03-25-00_FrontDoor.mp4'), event_id=None, camera_id=None,
                timestamp=datetime.datetime(2020, 4, 10, 3, 25, tzinfo=tzoffset('WST-8', 28800)), camera_name='FrontDoor', is_image=False,
                is_lowres=False)]

_PATH_DETAILS_BY_KEY_IMMEDIATE_MATCH = [[PathDetails(
    path=PosixPath('Segment_2020-04-10_03-24-32_SideGate.mp4'), event_id=None, camera_id=None,
    timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)), camera_name='SideGate', is_image=False,
    is_lowres=False), PathDetails(path=PosixPath('Segment_2020-04-10_03-24-32_SideGate.jpg'), event_id=None,
                                  camera_id=None, timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)),
                                  camera_name='SideGate', is_image=True, is_lowres=False), PathDetails(
    path=PosixPath('Segment_2020-04-10_03-24-32_SideGate-lowres.mp4'), event_id=None, camera_id=None,
    timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)), camera_name='SideGate', is_image=False,
    is_lowres=True), PathDetails(path=PosixPath('Segment_2020-04-10_03-24-32_SideGate-lowres.jpg'),
                                 event_id=None, camera_id=None,
                                 timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)),
                                 camera_name='SideGate', is_image=True, is_lowres=True)], [PathDetails(
    path=PosixPath('Segment_2020-04-10_03-24-32_Driveway.mp4'), event_id=None, camera_id=None,
    timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)), camera_name='Driveway', is_image=False,
    is_lowres=False), PathDetails(path=PosixPath('Segment_2020-04-10_03-24-32_Driveway.jpg'), event_id=None,
                                  camera_id=None, timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)),
                                  camera_name='Driveway', is_image=True, is_lowres=False), PathDetails(
    path=PosixPath('Segment_2020-04-10_03-24-32_Driveway-lowres.mp4'), event_id=None, camera_id=None,
    timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)), camera_name='Driveway', is_image=False,
    is_lowres=True), PathDetails(path=PosixPath('Segment_2020-04-10_03-24-32_Driveway-lowres.jpg'),
                                 event_id=None, camera_id=None,
                                 timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)),
                                 camera_name='Driveway', is_image=True, is_lowres=True)], [PathDetails(
    path=PosixPath('Segment_2020-04-10_03-24-32_FrontDoor.mp4'), event_id=None, camera_id=None,
    timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)), camera_name='FrontDoor', is_image=False,
    is_lowres=False), PathDetails(path=PosixPath('Segment_2020-04-10_03-24-32_FrontDoor.jpg'),
                                  event_id=None, camera_id=None,
                                  timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)),
                                  camera_name='FrontDoor', is_image=True, is_lowres=False), PathDetails(
    path=PosixPath('Segment_2020-04-10_03-24-32_FrontDoor-lowres.mp4'), event_id=None, camera_id=None,
    timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)), camera_name='FrontDoor', is_image=False,
    is_lowres=True), PathDetails(path=PosixPath('Segment_2020-04-10_03-24-32_FrontDoor-lowres.jpg'),
                                 event_id=None, camera_id=None,
                                 timestamp=datetime.datetime(2020, 4, 10, 3, 24, 32, tzinfo=tzoffset('WST-8', 28800)),
                                 camera_name='FrontDoor', is_image=True, is_lowres=True)]]

_SORTED_EVENTS = []

_JSON_LINES = """"""


class EventStoreRebuilderForEventsTest(unittest.TestCase):
    maxDiff = 65536 * 8

    def test_parse_path(self):
        self.assertEqual(_PATH_DETAILS, parse_path(path=_SOME_PATHS[3], tzinfo=_TZINFO))

    def test_parse_paths(self):
        print(parse_paths(paths=_SOME_PATHS, tzinfo=_TZINFO, parse_method=parse_path))

        self.assertEqual(
            _SOME_PATH_DETAILS, parse_paths(paths=_SOME_PATHS, tzinfo=_TZINFO, parse_method=parse_path)
        )

    def test_relate_path_details(self):
        # print(relate_path_details(_SOME_PATH_DETAILS, get_key_methods=[
        #     _get_key
        # ]))

        self.assertEqual(
            _PATH_DETAILS_BY_KEY_IMMEDIATE_MATCH,
            relate_path_details(_SOME_PATH_DETAILS, get_key_methods=[
                _get_key
            ]),
        )

    def test_rebuild_event_store(self):
        tmpdir = Path(mkdtemp())

        for path in _SOME_PATHS:
            full_path = tmpdir / path
            full_path.touch()

        json_path = tmpdir / Path("events.jsonl")

        rebuild_event_store(root_path=tmpdir, tzinfo=_TZINFO, json_path=json_path, parse_method=parse_path,
                            get_key_methods=[_get_key])

        with open(str(json_path), "r") as f:
            data = f.read().strip()

        self.assertEqual(_JSON_LINES.count("\n"), data.count("\n"))

        for path in _SORTED_EVENTS:
            self.assertIn(path.high_res_image_path.split("/")[-1], data)
            self.assertIn(path.low_res_image_path.split("/")[-1], data)
            self.assertIn(path.high_res_video_path.split("/")[-1], data)
            self.assertIn(path.low_res_video_path.split("/")[-1], data)

    def test_many_paths_with_poor_data_sanity(self):
        sorted_paths = _sorted_paths_segments
        tzinfo = _TZINFO

        some_path_details = parse_paths(paths=sorted_paths, tzinfo=tzinfo, parse_method=parse_path)
        some_path_details_by_key = relate_path_details(
            some_path_details=some_path_details, get_key_methods=[
                _get_key
            ]
        )

        self.assertEqual(18265, len(some_path_details_by_key))

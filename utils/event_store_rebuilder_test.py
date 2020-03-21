import datetime
import unittest

from .event_store_rebuilder import format_timestamp_for_go

_TEST_TIMESTAMP = datetime.datetime(year=1992, month=2, day=6, hour=11, minute=11, second=11, microsecond=11)


class EventStoreRebuilderTest(unittest.TestCase):
    def test_format_timestamp_for_go(self):
        self.assertEqual(
            '',
            format_timestamp_for_go()
        )

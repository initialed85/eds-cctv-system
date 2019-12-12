import datetime
import os
import unittest
from collections import OrderedDict

from event_parser import parse_events

_TIMESTAMP = datetime.datetime(year=1992, month=2, day=6, hour=13, minute=33, second=37)

_HTMLS = OrderedDict([('events_2019-06-07.html',
                       '</html>\n<head>\n<title>Events for 2019-06-07 as at 1992-02-06 13:33:37</title>\n<style type="text/css">\n\nBODY {\n    font-family: Tahoma;\n    font-size: 8pt;\n    font-weight: none;\n    text-align: center;\n}\n\nTH {\n    font-family: Tahoma;\n    font-size: 8pt;\n    font-weight: bold;\n    text-align: center;\n}\n\nTD {\n    font-family: Tahoma;\n    font-size: 8pt;\n    font-weight: none;\n    text-align: center;\n    border: 1px solid gray; \n}\n\n</style>\n</head>\n\n<body>\n<h1>Events for 2019-06-07 as at 1992-02-06 13:33:37</h1>\n\n<center>\n<table width="90%">\n\n<tr>\n<th>Event ID</th>\n<th>Camera ID</th>\n<th>Timestamp</th>\n<th>Size</th>\n<th>Camera</th>\n<th>Screenshot</th>\n<th>Download</th>\n</tr>\n\n<tr>\n<td>73</td>\n<td>102</td>\n<td>2019-06-07 00:09:29</td>\n<td>0.0 MB</td>\n<td>FrontDoor</td>\n<td style="width: 320px";><a target="_blank" href="/browse58__102__2019-06-06_22-49-40__FrontDoor.jpg"><img src="/browse58__102__2019-06-06_22-49-40__FrontDoor.jpg" alt="58__102__2019-06-06_22-49-40__FrontDoor.jpg" width="320" height="180" /></a></td>\n<td><a href="/browse73__102__2019-06-07_00-09-29__FrontDoor.mkv">Download</a></td>\n</tr>\n\n<tr>\n<td>72</td>\n<td>102</td>\n<td>2019-06-07 00:08:58</td>\n<td>0.0 MB</td>\n<td>FrontDoor</td>\n<td style="width: 320px";><a target="_blank" href="/browse72__102__2019-06-07_00-08-59__FrontDoor.jpg"><img src="/browse72__102__2019-06-07_00-08-59__FrontDoor.jpg" alt="72__102__2019-06-07_00-08-59__FrontDoor.jpg" width="320" height="180" /></a></td>\n<td><a href="/browse72__102__2019-06-07_00-08-58__FrontDoor.mkv">Download</a></td>\n</tr>\n\n</table>\n<center>\n\n</body>\n</html>'),
                      ('events_2019-06-06.html',
                       '</html>\n<head>\n<title>Events for 2019-06-06 as at 1992-02-06 13:33:37</title>\n<style type="text/css">\n\nBODY {\n    font-family: Tahoma;\n    font-size: 8pt;\n    font-weight: none;\n    text-align: center;\n}\n\nTH {\n    font-family: Tahoma;\n    font-size: 8pt;\n    font-weight: bold;\n    text-align: center;\n}\n\nTD {\n    font-family: Tahoma;\n    font-size: 8pt;\n    font-weight: none;\n    text-align: center;\n    border: 1px solid gray; \n}\n\n</style>\n</head>\n\n<body>\n<h1>Events for 2019-06-06 as at 1992-02-06 13:33:37</h1>\n\n<center>\n<table width="90%">\n\n<tr>\n<th>Event ID</th>\n<th>Camera ID</th>\n<th>Timestamp</th>\n<th>Size</th>\n<th>Camera</th>\n<th>Screenshot</th>\n<th>Download</th>\n</tr>\n\n<tr>\n<td>61</td>\n<td>101</td>\n<td>2019-06-06 23:44:20</td>\n<td>0.0 MB</td>\n<td>Driveway</td>\n<td style="width: 320px";><a target="_blank" href="/browse61__101__2019-06-06_23-44-22__Driveway.jpg"><img src="/browse61__101__2019-06-06_23-44-22__Driveway.jpg" alt="61__101__2019-06-06_23-44-22__Driveway.jpg" width="320" height="180" /></a></td>\n<td><a href="/browse61__101__2019-06-06_23-44-20__Driveway.mkv">Download</a></td>\n</tr>\n\n<tr>\n<td>60</td>\n<td>101</td>\n<td>2019-06-06 23:43:18</td>\n<td>0.0 MB</td>\n<td>Driveway</td>\n<td style="width: 320px";><a target="_blank" href="/browse60__101__2019-06-06_23-43-19__Driveway.jpg"><img src="/browse60__101__2019-06-06_23-43-19__Driveway.jpg" alt="60__101__2019-06-06_23-43-19__Driveway.jpg" width="320" height="180" /></a></td>\n<td><a href="/browse60__101__2019-06-06_23-43-18__Driveway.mkv">Download</a></td>\n</tr>\n\n<tr>\n<td>59</td>\n<td>101</td>\n<td>2019-06-06 23:38:03</td>\n<td>0.0 MB</td>\n<td>Driveway</td>\n<td style="width: 320px";><a target="_blank" href="/browse59__101__2019-06-06_23-38-05__Driveway.jpg"><img src="/browse59__101__2019-06-06_23-38-05__Driveway.jpg" alt="59__101__2019-06-06_23-38-05__Driveway.jpg" width="320" height="180" /></a></td>\n<td><a href="/browse59__101__2019-06-06_23-38-03__Driveway.mkv">Download</a></td>\n</tr>\n\n<tr>\n<td>71</td>\n<td>102</td>\n<td>2019-06-06 23:32:34</td>\n<td>0.0 MB</td>\n<td>FrontDoor</td>\n<td style="width: 320px";><a target="_blank" href="/browse71__102__2019-06-06_23-32-36__FrontDoor.jpg"><img src="/browse71__102__2019-06-06_23-32-36__FrontDoor.jpg" alt="71__102__2019-06-06_23-32-36__FrontDoor.jpg" width="320" height="180" /></a></td>\n<td><a href="/browse71__102__2019-06-06_23-32-34__FrontDoor.mkv">Download</a></td>\n</tr>\n\n<tr>\n<td>58</td>\n<td>101</td>\n<td>2019-06-06 23:32:34</td>\n<td>0.0 MB</td>\n<td>Driveway</td>\n<td style="width: 320px";><a target="_blank" href="/browse59__102__2019-06-06_22-50-18__FrontDoor.jpg"><img src="/browse59__102__2019-06-06_22-50-18__FrontDoor.jpg" alt="59__102__2019-06-06_22-50-18__FrontDoor.jpg" width="320" height="180" /></a></td>\n<td><a href="/browse58__101__2019-06-06_23-32-34__Driveway.mkv">Download</a></td>\n</tr>\n\n<tr>\n<td>70</td>\n<td>102</td>\n<td>2019-06-06 23:32:13</td>\n<td>0.0 MB</td>\n<td>FrontDoor</td>\n<td style="width: 320px";><a target="_blank" href="/browse70__102__2019-06-06_23-32-15__FrontDoor.jpg"><img src="/browse70__102__2019-06-06_23-32-15__FrontDoor.jpg" alt="70__102__2019-06-06_23-32-15__FrontDoor.jpg" width="320" height="180" /></a></td>\n<td><a href="/browse70__102__2019-06-06_23-32-13__FrontDoor.mkv">Download</a></td>\n</tr>\n\n<tr>\n<td>69</td>\n<td>102</td>\n<td>2019-06-06 23:31:29</td>\n<td>0.0 MB</td>\n<td>FrontDoor</td>\n<td style="width: 320px";><a target="_blank" href="/browse69__102__2019-06-06_23-31-31__FrontDoor.jpg"><img src="/browse69__102__2019-06-06_23-31-31__FrontDoor.jpg" alt="69__102__2019-06-06_23-31-31__FrontDoor.jpg" width="320" height="180" /></a></td>\n<td><a href="/browse69__102__2019-06-06_23-31-29__FrontDoor.mkv">Download</a></td>\n</tr>\n\n<tr>\n<td>68</td>\n<td>102</td>\n<td>2019-06-06 23:17:13</td>\n<td>0.0 MB</td>\n<td>FrontDoor</td>\n<td style="width: 320px";><a target="_blank" href="/browse68__102__2019-06-06_23-17-15__FrontDoor.jpg"><img src="/browse68__102__2019-06-06_23-17-15__FrontDoor.jpg" alt="68__102__2019-06-06_23-17-15__FrontDoor.jpg" width="320" height="180" /></a></td>\n<td><a href="/browse68__102__2019-06-06_23-17-13__FrontDoor.mkv">Download</a></td>\n</tr>\n\n<tr>\n<td>67</td>\n<td>102</td>\n<td>2019-06-06 23:10:52</td>\n<td>0.0 MB</td>\n<td>FrontDoor</td>\n<td style="width: 320px";><a target="_blank" href="/browse67__102__2019-06-06_23-10-54__FrontDoor.jpg"><img src="/browse67__102__2019-06-06_23-10-54__FrontDoor.jpg" alt="67__102__2019-06-06_23-10-54__FrontDoor.jpg" width="320" height="180" /></a></td>\n<td><a href="/browse67__102__2019-06-06_23-10-52__FrontDoor.mkv">Download</a></td>\n</tr>\n\n<tr>\n<td>66</td>\n<td>102</td>\n<td>2019-06-06 23:09:00</td>\n<td>0.0 MB</td>\n<td>FrontDoor</td>\n<td style="width: 320px";><a target="_blank" href="/browse66__102__2019-06-06_23-09-02__FrontDoor.jpg"><img src="/browse66__102__2019-06-06_23-09-02__FrontDoor.jpg" alt="66__102__2019-06-06_23-09-02__FrontDoor.jpg" width="320" height="180" /></a></td>\n<td><a href="/browse66__102__2019-06-06_23-09-00__FrontDoor.mkv">Download</a></td>\n</tr>\n\n<tr>\n<td>65</td>\n<td>102</td>\n<td>2019-06-06 23:03:16</td>\n<td>0.0 MB</td>\n<td>FrontDoor</td>\n<td style="width: 320px";><a target="_blank" href="/browse65__102__2019-06-06_23-03-18__FrontDoor.jpg"><img src="/browse65__102__2019-06-06_23-03-18__FrontDoor.jpg" alt="65__102__2019-06-06_23-03-18__FrontDoor.jpg" width="320" height="180" /></a></td>\n<td><a href="/browse65__102__2019-06-06_23-03-16__FrontDoor.mkv">Download</a></td>\n</tr>\n\n<tr>\n<td>64</td>\n<td>102</td>\n<td>2019-06-06 22:59:54</td>\n<td>0.0 MB</td>\n<td>FrontDoor</td>\n<td style="width: 320px";><a target="_blank" href="/browse64__102__2019-06-06_22-59-55__FrontDoor.jpg"><img src="/browse64__102__2019-06-06_22-59-55__FrontDoor.jpg" alt="64__102__2019-06-06_22-59-55__FrontDoor.jpg" width="320" height="180" /></a></td>\n<td><a href="/browse64__102__2019-06-06_22-59-54__FrontDoor.mkv">Download</a></td>\n</tr>\n\n<tr>\n<td>63</td>\n<td>102</td>\n<td>2019-06-06 22:55:52</td>\n<td>0.0 MB</td>\n<td>FrontDoor</td>\n<td style="width: 320px";><a target="_blank" href="/browse63__102__2019-06-06_22-55-54__FrontDoor.jpg"><img src="/browse63__102__2019-06-06_22-55-54__FrontDoor.jpg" alt="63__102__2019-06-06_22-55-54__FrontDoor.jpg" width="320" height="180" /></a></td>\n<td><a href="/browse63__102__2019-06-06_22-55-52__FrontDoor.mkv">Download</a></td>\n</tr>\n\n<tr>\n<td>62</td>\n<td>102</td>\n<td>2019-06-06 22:54:38</td>\n<td>0.0 MB</td>\n<td>FrontDoor</td>\n<td style="width: 320px";><a target="_blank" href="/browse62__102__2019-06-06_22-54-40__FrontDoor.jpg"><img src="/browse62__102__2019-06-06_22-54-40__FrontDoor.jpg" alt="62__102__2019-06-06_22-54-40__FrontDoor.jpg" width="320" height="180" /></a></td>\n<td><a href="/browse62__102__2019-06-06_22-54-38__FrontDoor.mkv">Download</a></td>\n</tr>\n\n<tr>\n<td>61</td>\n<td>102</td>\n<td>2019-06-06 22:52:25</td>\n<td>0.0 MB</td>\n<td>FrontDoor</td>\n<td style="width: 320px";><a target="_blank" href="/browse61__102__2019-06-06_22-52-27__FrontDoor.jpg"><img src="/browse61__102__2019-06-06_22-52-27__FrontDoor.jpg" alt="61__102__2019-06-06_22-52-27__FrontDoor.jpg" width="320" height="180" /></a></td>\n<td><a href="/browse61__102__2019-06-06_22-52-25__FrontDoor.mkv">Download</a></td>\n</tr>\n\n<tr>\n<td>60</td>\n<td>102</td>\n<td>2019-06-06 22:52:04</td>\n<td>0.0 MB</td>\n<td>FrontDoor</td>\n<td style="width: 320px";><a target="_blank" href="/browse60__102__2019-06-06_22-52-05__FrontDoor.jpg"><img src="/browse60__102__2019-06-06_22-52-05__FrontDoor.jpg" alt="60__102__2019-06-06_22-52-05__FrontDoor.jpg" width="320" height="180" /></a></td>\n<td><a href="/browse60__102__2019-06-06_22-52-04__FrontDoor.mkv">Download</a></td>\n</tr>\n\n<tr>\n<td>59</td>\n<td>102</td>\n<td>2019-06-06 22:50:17</td>\n<td>0.0 MB</td>\n<td>FrontDoor</td>\n<td style="width: 320px";><a target="_blank" href="/browseevents/missing.png"><img src="/browseevents/missing.png" alt="events/missing.png" width="320" height="180" /></a></td>\n<td><a href="/browse59__102__2019-06-06_22-50-17__FrontDoor.mkv">Download</a></td>\n</tr>\n\n<tr>\n<td>58</td>\n<td>102</td>\n<td>2019-06-06 22:49:39</td>\n<td>0.0 MB</td>\n<td>FrontDoor</td>\n<td style="width: 320px";><a target="_blank" href="/browseevents/missing.png"><img src="/browseevents/missing.png" alt="events/missing.png" width="320" height="180" /></a></td>\n<td><a href="/browse58__102__2019-06-06_22-49-39__FrontDoor.mkv">Download</a></td>\n</tr>\n\n</table>\n<center>\n\n</body>\n</html>'),
                      ('events.html',
                       '</html>\n<head>\n<title>All events as at 1992-02-06 13:33:37</title>\n<style type="text/css">\n\nBODY {\n    font-family: Tahoma;\n    font-size: 8pt;\n    font-weight: none;\n    text-align: center;\n}\n\nTH {\n    font-family: Tahoma;\n    font-size: 8pt;\n    font-weight: bold;\n    text-align: center;\n}\n\nTD {\n    font-family: Tahoma;\n    font-size: 8pt;\n    font-weight: none;\n    text-align: center;\n    border: 1px solid gray; \n}\n\n</style>\n</head>\n\n<body>\n<h2>Events as at 1992-02-06 13:33:37</h2>\n\n<center>\n<table width="90%">\n\n<tr>\n<th>Date</th>\n<th>Events</th>\n</tr>\n\n\n<tr>\n<td><a target="event" href="events_2019-06-07.html">2019-06-07</a></td>\n<td>2</td>\n</tr>\n\n\n\n<tr>\n<td><a target="event" href="events_2019-06-06.html">2019-06-06</a></td>\n<td>18</td>\n</tr>\n\n\n</table>\n<center>\n\n</body>\n</html>\n')])


class EventParserTest(unittest.TestCase):
    maxDiff = 65536 * 8

    def setUp(self):
        os.chdir(os.path.dirname(os.path.realpath(__file__)))

    def test_parse_events(self):
        print(parse_events(
            target_dir='../test_files',
            browse_url_prefix='/browse',
            run_timestamp=_TIMESTAMP
        ))

        self.assertEqual(
            parse_events(
                target_dir='../test_files',
                browse_url_prefix='/browse',
                run_timestamp=_TIMESTAMP
            ),
            _HTMLS
        )

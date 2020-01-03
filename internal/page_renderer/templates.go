package page_renderer

// Common

var StyleSheet = `BODY {
    font-family: Tahoma;
    font-size: 8pt;
    font-weight: none;
    text-align: center;
}

TH {
    font-family: Tahoma;
    font-size: 8pt;
    font-weight: bold;
    text-align: center;
}

TD {
    font-family: Tahoma;
    font-size: 8pt;
    font-weight: none;
    text-align: center;
    border: 1px solid gray; 
}`

// For the summary view

var SummaryTableRowHTML = `
	<tr>
		<td><a target="event" href="{{.EventsURL}}">{{.EventsDate}}</a></td>
		<td>{{.EventCount}}</td>
	</tr>
`

var SummaryHTML = `</html>
<head>
<title>{{.Title}} as at {{.Now}}</title>
<style type="text/css">
{{.StyleSheet}}
</style>
</head>

<body>
<h2>{{.Title}} as at {{.Now}}</h2>

<center>
<table width="90%">

	<tr>
		<th>Date</th>
		<th>Events</th>
	</tr>
{{.TableRows}}
</table>
<center>

</body>
</html>`

// For the drill-down view

var PageTableRowHTML = `
	<tr>
		<td>{{.EventID}}</td>
		<td>{{.Timestamp}}</td>
		<td>{{.CameraName}}</td>
		<td style="width: 320px";><a target="_blank" href="{{.HighResImageURL}}"><img src="{{.LowResImageURL}}" width="320" height="180" /></a></td>
		<td>Download <a href="{{.HighResVideoURL}}">high-res</a> or <a href="{{.LowResVideoURL}}">low-res</a></td>
	</tr>
`

var PageHTML = `</html>
<head>
<title>{{.Title}} for {{.EventsDate}} as at {{.Now}}</title>
<style type="text/css">
{{.StyleSheet}}
</style>
</head>

<body>
<h1>{{.Title}} for {{.EventsDate}} as at {{.Now}}</h1>

<center>
<table width="90%">

	<tr>
		<th>Event ID</th>
		<th>Timestamp</th>
		<th>Camera</th>
		<th>Screenshot</th>
		<th>Download</th>
	</tr>
{{.TableRows}}
</table>
<center>

</body>
</html>`

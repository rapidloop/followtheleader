/*

followtheleader - follow the presidential candidates with Go and Twitter

Copyright (c) 2015 RapidLoop

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

var tmpl *template.Template

func startWeb() {
	tmpl = template.Must(template.New(".").Parse(html))
	http.HandleFunc("/", webServer)
	log.Fatal(http.ListenAndServe(LISTEN_ADDRESS, nil))
}

type templateData struct {
	At      time.Time
	Twitter []TwitterInfo
}

func webServer(w http.ResponseWriter, r *http.Request) {
	a, t := stats.Get()
	if err := tmpl.Execute(w, templateData{a, t}); err != nil {
		log.Print(err)
	}
}

const html = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>The 45th POTUS</title>
    <link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css" rel="stylesheet">
    <script type="text/javascript" src="https://www.google.com/jsapi"></script>
    <link href='https://fonts.googleapis.com/css?family=PT+Sans' rel='stylesheet' type='text/css'>
    <style>
        body { font-family: 'PT Sans', 'Avenir', sans-serif; }
    </style>
    <script type="text/javascript">
        google.load('visualization', '1', {packages: ['corechart', 'bar']});
        google.setOnLoadCallback(drawChart);
        function drawChart() {
          var data = google.visualization.arrayToDataTable([
            ['Candidate', 'Followers', {role: 'style'}],
            {{range .Twitter}}
            ['{{.Name}}', {{.Followers}}, {{if .Democrat}}'#00AEF3'{{else}}'#DF1A22'{{end}}, ],
            {{end}}
          ]);

          var options = {
            chartArea: {width: '61%', height: '90%'},
            hAxis: {
              minValue: 0
            },
            dataOpacity: 0.8,
            legend: {
                position: 'none',
            },
            fontName: 'PT Sans',
          };

          var chart = new google.visualization.BarChart(document.getElementById('chart'));

          chart.draw(data, options);
        }
    </script>
  </head>
  <body>
    <div class="container">
      <div class="row">
        <h1 style="text-align: center">Who Has More Followers on Twitter?</h1>
        <div style="width: 100%; text-align: center; color: #888">Data as on {{.At.Format "1/2/2006 3:04 PM MST"}}</div>
      </div>
      <div class="row" style="text-align: center">
        <div class="col-sm-12">
          <div id="chart" style="height: 500px;"></div>
        </div>
      </div>
      <div class="row" style="text-align: center; margin-top: 5em">
        <div class="col-sm-12" style="color: #aaa">
          Brought to you by <a
          href="https://www.rapidloop.com/">RapidLoop</a> as a fun project, not intended for serious use!<br>
          If you love keeping an eye on data, checkout <a href="https://www.opsdash.com/">OpsDash</a> - Server and Service Monitoring at $1/server/month.
        </div>
      </div>
    </div>
  </body>
</html>
`

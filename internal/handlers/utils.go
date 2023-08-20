package handlers

import (
	"bytes"
	"certalert/internal/certificates"
	"fmt"
	"log"
	"strings"
	"text/template"
	"time"
)

// remainingDuration returns the remaining duration from the given epoch time
var remainingDuration = func(epoch int64) time.Duration {
	return time.Until(time.Unix(epoch, 0))
}

// getRowColor returns the color of the row based on the expiry date
func getRowColor(epoch int64) string {
	if epoch == 0 {
		return ""
	}

	d := remainingDuration(epoch)

	// expired
	if d <= 0 {
		return "red-row"
	}

	// expires in the next 3 days
	if d <= 3*24*time.Hour {
		return "red-row"
	}

	// expires in the next 30 days
	if d <= 30*24*time.Hour {
		return "orange-row"
	}

	// expires in the next 60 days
	if d <= 60*24*time.Hour {
		return "yellow-row"
	}

	return ""
}

// epochToHumanReadable converts the epoch time to human readable format
func epochToHumanReadable(epoch int64) string {
	if epoch == 0 {
		return "-"
	}

	d := remainingDuration(epoch)

	// expired
	if d <= 0 {
		return "now"
	}

	days := int(d / (24 * time.Hour))
	d -= time.Duration(days) * 24 * time.Hour

	hours := int(d / time.Hour)
	d -= time.Duration(hours) * time.Hour

	minutes := int(d / time.Minute)
	d -= time.Duration(minutes) * time.Minute

	seconds := int(d / time.Second)

	parts := []string{}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d days", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d hours", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d minutes", minutes))
	}
	if seconds > 0 {
		parts = append(parts, fmt.Sprintf("%d seconds", seconds))
	}

	return fmt.Sprintf(strings.Join(parts, ", "))
}

// formatTime formats the given time with the given format
func formatTime(t time.Time, format string) string {
	// check if the time is zero or time is not set
	if t.IsZero() || t.Unix() == 0 {
		return "-"
	}
	return t.Format(format)
}

func renderTemplate(baseTplStr string, tplStr string, data interface{}) string {
	funcMap := template.FuncMap{
		"formatTime":    formatTime,
		"humanReadable": epochToHumanReadable,
		"getRowColor":   getRowColor,
	}

	// Create a new template and parse the base template into it.
	t := template.New("base").Funcs(funcMap)
	t, err := t.Parse(baseTplStr)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new template that is associated with the previous one, and parse the specific template into it.
	t, err = t.New("content").Parse(tplStr)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "base", data); err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

type TemplateData struct {
	CSS       string
	Endpoints []Endpoint
	CertInfos []certificates.CertificateInfo
}

const CSS = `
.table {
	border-collapse: collapse;
	width: 60%;
	margin: 0 auto;
	border: 1px solid #ddd;
	font-size: 16px;
}
.table th,
.table td {
	text-align: left;
	padding: 12px;
}
.table tr:not(.table-header) {
	border-bottom: 1px solid #ddd;
}
.table tr:not(.table-header):hover {
	background-color: #f1f1f1;
}
.table-header {
	background-color: #BDB76B;
}

thead th {
	position: sticky;
	top: 0;
	z-index: 1;
	background: #BDB76B;
}

.error-symbol:hover {
	opacity: 0.7;
}

.row-yellow {
	background-color: #FFD700;
}

.row-orange {
	background-color: #FFA500;
}

.row-red {
	background-color: #FF4500;
}
`

const tplBase = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<style>
		{{ .CSS }}
	</style>
</head>
<body>
	<div style="text-align:center;">
		{{ template "content" . }}
	</div>
</body>
</html>
`

const tplCertificates = `
{{ define "content" }}
	<table class="table">
			<thead>
					<tr class="table-header">
							<th scope="col"></th>
							<th scope="col">Name</th>
							<th scope="col">Subject</th>
							<th scope="col">Type</th>
							<th scope="col">Expiry Date</th>
							<th scope="col">Expiration</th>
					</tr>
			</thead>
			<tbody>
					{{range .CertInfos}}
					<tr class="{{ getRowColor .Epoch }}">
							<td>
									{{if .Error}}
											<span class="error-symbol" title="{{.Error}}" style="color: red;">✖</span>
									{{else}}
											<span style="color: green;">✔</span>
									{{end}}
							</td>
							<td>{{.Name}}</td>
							<td>{{.Subject}}</td>
							<td>{{.Type}}</td>
							<td>{{ formatTime .ExpiryAsTime "2006-01-02" }}</td>
							<td>{{ humanReadable .Epoch }}</td>
					</tr>
					{{end}}
			</tbody>
	</table>
{{ end }}
`

const tplEndpoints = `
{{ define "content" }}
	<table class="table">
		<thead>
				<tr class="table-header">
						<th>Endpoint</th>
						<th>Methods</th>
						<th>Purpose</th>
				</tr>
		</thead>
		<tbody>
				{{range .Endpoints}}
				<tr>
						<td><a href="{{.Path}}">{{.Path}}</a></td>
						<td>{{range .Methods}}{{.}} {{end}}</td>
						<td>{{.Description}}</td>
				</tr>
				{{end}}
		</tbody>
	</table>
{{ end }}
`

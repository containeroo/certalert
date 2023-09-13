package handlers

// CSS is the CSS that is used in the template
const CSS string = `
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

.sortable:hover {
  cursor: pointer;
  text-decoration: underline;
}

.sort-asc:after {
  content: " ↑";
}

.sort-desc:after {
  content: " ↓";
}

code {
    background-color: #f0f0f0;
    color: #333;
    padding: 2px 4px;
    border: 1px solid #ccc;
    border-radius: 3px;
    font-family: 'Courier New', monospace;
}

`

// JS is the JS that is used in the template
const JS string = `
let currentColumn = -1;
let sortAscending = true;

function sortTable(columnIndex) {
	const table = document.querySelector('.table');
	const headers = Array.from(table.querySelectorAll('thead th'));
	const rows = Array.from(table.querySelectorAll('tbody tr'));

	// Remove previous sort direction classes
	headers.forEach(header => {
		header.classList.remove('sort-asc', 'sort-desc');
	});

	if (columnIndex === currentColumn) {
		sortAscending = !sortAscending;
	} else {
		currentColumn = columnIndex;
		sortAscending = true;
	}

	headers[columnIndex].classList.add(sortAscending ? 'sort-asc' : 'sort-desc');

	rows.sort((a, b) => {
		const cellA = a.cells[columnIndex].textContent;
		const cellB = b.cells[columnIndex].textContent;
		return sortAscending ? cellA.localeCompare(cellB) : cellB.localeCompare(cellA);
	});

	const tbody = table.querySelector('tbody');
	rows.forEach(row => tbody.appendChild(row));
}
`

// tplBase is the base template
// Add "content" to this base template
const tplBase string = `
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
<script>
	{{ .JS }}
</script>
</html>
`

// tplCertificates is the template for the /certificates route
const tplCertificates string = `
{{ define "content" }}
	<table class="table">
			<thead>
					<tr class="table-header">
							<th scope="col"></th>
							<th class="sortable" onclick="sortTable(1)">Name</th>
							<th class="sortable" onclick="sortTable(2)">Subject</th>
							<th class="sortable" onclick="sortTable(3)">Type</th>
							<th class="sortable" onclick="sortTable(4)">Expiry Date</th>
							<th class="sortable" onclick="sortTable(5)">Expiration</th>
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

// tplEndpoints is the template for the / route
const tplEndpoints string = `
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
            <td>
              {{range .Methods}}
                <code>{{.}}</code>
              {{end}}
            </td>
						<td>{{.Description}}</td>
				</tr>
				{{end}}
		</tbody>
	</table>
{{ end }}
`

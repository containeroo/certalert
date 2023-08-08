package handlers

import (
	"bytes"
	"certalert/internal/certificates"
	"log"
	"text/template"
)

const tpl = `
<!DOCTYPE html>
<html>
<head>
		<meta charset="UTF-8">
    <style>
		#myTable {
			border-collapse: collapse;
			width: 60%;
			margin: 0 auto;
			border: 1px solid #ddd;
			font-size: 16px;
		}
		#myTable th,
		#myTable td {
			text-align: left;
			padding: 12px;
		}
		#myTable tr:not(.header) {
			border-bottom: 1px solid #ddd;
		}
		#myTable tr:not(.header):hover {
			background-color: #f1f1f1;
		}
		#myTable tr.header {
			background-color: #BDB76B;
		}

		/* Add hover effect for error symbol */
		.error-symbol:hover {
				opacity: 0.7;
		}
	</style>
</head>
<body>
	<table id="myTable">
        <thead>
            <tr class="header">
								<th></th>
                <th>Name</th>
                <th>Subject</th>
                <th>Type</th>
                <th>Expiry Date</th>
            </tr>
        </thead>
        <tbody>
            {{range .}}
            <tr>
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
                <td>{{.ExpiryAsTime.Format "2006-01-02"}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
</body>
</html>

`

// renderCertificateInfo renders the certificate information as HTML
func renderCertificateInfo(certInfo []certificates.CertificateInfo) string {
	t, err := template.New("certInfo").Parse(tpl)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, certInfo); err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

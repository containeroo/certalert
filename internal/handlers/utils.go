package handlers

import (
	"bytes"
	"certalert/internal/certificates"
	"fmt"
	"strings"
	"text/template"
	"time"
)

// TemplateData is the data that is passed to the template
type TemplateData struct {
	CSS       string
	JS        string
	Endpoints []Handler
	CertInfos []certificates.CertificateInfo
}

// renderTemplate renders the given template with the given data
func renderTemplate(baseTplStr string, tplStr string, data interface{}) (string, error) {
	funcMap := template.FuncMap{
		"formatTime":    formatTime,
		"humanReadable": epochToHumanReadable,
		"getRowColor":   getRowColor,
	}

	// Create a new template and parse the base template into it.
	t, err := template.New("base").Funcs(funcMap).Parse(baseTplStr)
	if err != nil {
		return "", err
	}

	// Create a new template that is associated with the previous one, and parse the specific template into it.
	t, err = t.New("content").Parse(tplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "base", data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// remainingDuration returns the remaining duration from the given epoch time
// Only used so getRowColor can be tested
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

	return fmt.Sprint(strings.Join(parts, ", "))
}

// formatTime formats the given time with the given format
func formatTime(t time.Time, format string) string {
	// check if the time is zero or time is not set
	if t.IsZero() || t.Unix() == 0 {
		return "-"
	}
	return t.Format(format)
}

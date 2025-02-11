package formatter

// OutputFormat is a resulting type of file that is converted/formatted from XML to HTML (for example)
type OutputFormat string

const (
	// HTMLOutput constant defines OutputFormat is HyperText Markup Language which can be viewed using browsers
	HTMLOutput OutputFormat = "html"
	// CSVOutput constant defines OutputFormat for Comma-Separated Values CSV file which is viewed most of the time in Excel
	CSVOutput OutputFormat = "csv"
	// MarkdownOutput constant defines OutputFormat for Markdown, which is handy and easy format to read-write
	MarkdownOutput OutputFormat = "md"
	// JSONOutput constant defines OutputFormat for JavaScript Object Notation, which is more useful for machine-related operations (parsing)
	JSONOutput OutputFormat = "json"
)

// IsValid checks whether requested output format is valid
func (of OutputFormat) IsValid() bool {
	// markdown & md is essentially the same thing
	switch of {
	case "markdown", "md", "html", "csv", "json":
		return true
	}
	return false
}

// FileOutputFormat returns appropriate file format, users can provide short
func (of OutputFormat) FileOutputFormat() OutputFormat {
	switch of {
	case "markdown", "md":
		return MarkdownOutput
	case "html":
		return HTMLOutput
	case "csv":
		return CSVOutput
	case "json":
		return JSONOutput
	}
	return HTMLOutput
}

package formatter

import (
	"encoding/csv"
	"fmt"
)

// CSVFormatter is struct defined for CSV Output use-case
type CSVFormatter struct {
	config *Config
}

// Format the data to CSV and output it to appropriate io.Writer
func (f *CSVFormatter) Format(td *TemplateData) (err error) {
	return csv.NewWriter(f.config.Writer).WriteAll(f.convert(td))
}

// convert uses NMAPRun struct to convert all data to [][]string type
func (f *CSVFormatter) convert(td *TemplateData) (data [][]string) {
	data = append(data, []string{"IP", "Port", "Protocol", "State", "Service", "Reason", "Product", "Version", "Extra info"})
	for i := range td.NMAPRun.Host {
		var host *Host = &td.NMAPRun.Host[i]
		// Skipping hosts that are down
		if td.OutputOptions.SkipDownHosts && host.Status.State != "up" {
			continue
		}
		address := fmt.Sprintf("%s (%s)", host.HostAddress.Address, host.Status.State)
		data = append(data, []string{address, "", "", "", "", "", "", "", ""})
		for j := range host.Ports.Port {
			var port *Port = &host.Ports.Port[j]
			data = append(
				data,
				[]string{
					"",
					port.PortID,
					port.Protocol,
					port.State.State,
					port.Service.Name,
					port.State.Reason,
					port.Service.Product,
					port.Service.Version,
					port.Service.ExtraInfo,
				},
			)
		}
	}
	return
}

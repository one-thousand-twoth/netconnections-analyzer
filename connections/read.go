// Package connections reads declared connections and exposes ordered host and
// declaration data independently of the CLI and DOT writer.
package connections

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// Host is a named endpoint. Metadata describes the host, not a declaration.
type Host struct {
	Name     string
	Metadata string
}

// DeclaredConnection is one asserted source-to-destination relationship.
// Repeated assertions remain separate entries. Ports, protocol, and metadata
// are plain text rather than verified network configuration.
type DeclaredConnection struct {
	Source      string
	Destination string
	Ports       string
	Protocol    string
	Metadata    string
}

// Data holds one host per case-sensitive name in first-seen order, considering
// the source before the destination of each valid input record. Connections
// retain input-record order, including duplicates and self-loops. Its slices
// and entries can also be constructed by callers without CSV.
type Data struct {
	Hosts       []Host
	Connections []DeclaredConnection
}

// Mode controls whether invalid nonblank input records can be skipped.
type Mode uint8

const (
	DefaultMode Mode = iota
	StrictMode
)

// Warning identifies a skipped input record by physical starting line and field.
type Warning struct {
	Line    int
	Field   string
	Message string
}

// Read converts a minimal src;dest CSV from r into ordered data. A fatal error
// returns no usable data or earlier warnings, even if earlier records were valid.
func Read(r io.Reader, mode Mode) (Data, []Warning, error) {
	csvReader := csv.NewReader(r)
	csvReader.Comma = ';'
	header, err := csvReader.Read()
	if err != nil {
		return Data{}, nil, fmt.Errorf("read CSV header: %w", err)
	}
	if len(header) != 2 || header[0] != "src" || header[1] != "dest" {
		return Data{}, nil, fmt.Errorf("CSV header must be src;dest")
	}

	var data Data
	var warnings []Warning
	seen := make(map[string]bool)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return Data{}, nil, fmt.Errorf("read CSV record: %w", err)
		}
		line, _ := csvReader.FieldPos(0)
		src := strings.TrimSpace(record[0])
		dest := strings.TrimSpace(record[1])
		if src == "" && dest == "" {
			continue
		}
		if src == "" || dest == "" {
			field := "src"
			if src != "" {
				field = "dest"
			}
			message := fmt.Sprintf("line %d: missing %s", line, field)
			if mode == StrictMode {
				return Data{}, nil, fmt.Errorf("%s", message)
			}
			warnings = append(warnings, Warning{Line: line, Field: field, Message: message})
			continue
		}
		for _, name := range []string{src, dest} {
			if !seen[name] {
				data.Hosts = append(data.Hosts, Host{Name: name})
				seen[name] = true
			}
		}
		data.Connections = append(data.Connections, DeclaredConnection{Source: src, Destination: dest})
	}
	return data, warnings, nil
}

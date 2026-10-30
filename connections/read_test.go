package connections_test

import (
	"errors"
	"io"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/one-thousand-twoth/netconnections-analyzer/connections"
)

func TestReadMinimalDeclarationsInInputOrder(t *testing.T) {
	input := "src;dest\nWeb;db\nWeb;db\ndb;db\nweb;Web\n"
	data, warnings, err := connections.Read(strings.NewReader(input), connections.DefaultMode)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(warnings) != 0 {
		t.Fatalf("warnings: %+v", warnings)
	}
	wantHosts := []connections.Host{{Name: "Web"}, {Name: "db"}, {Name: "web"}}
	wantConnections := []connections.DeclaredConnection{
		{Source: "Web", Destination: "db"},
		{Source: "Web", Destination: "db"},
		{Source: "db", Destination: "db"},
		{Source: "web", Destination: "Web"},
	}
	if !slices.Equal(data.Hosts, wantHosts) || !slices.Equal(data.Connections, wantConnections) {
		t.Fatalf("got hosts %+v and declarations %+v; want %+v and %+v", data.Hosts, data.Connections, wantHosts, wantConnections)
	}
}

func TestReadSkipsBlankRecordsAndWarnsWithPhysicalStartingLines(t *testing.T) {
	input := "\nsrc;dest\n\"front;\nend\"; back \n ; \n;db\nweb;\n"
	data, warnings, err := connections.Read(strings.NewReader(input), connections.DefaultMode)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	want := []connections.DeclaredConnection{{Source: "front;\nend", Destination: "back"}}
	if !slices.Equal(data.Connections, want) {
		t.Fatalf("declarations: %+v, want %+v", data.Connections, want)
	}
	if !slices.Equal(data.Hosts, []connections.Host{{Name: "front;\nend"}, {Name: "back"}}) {
		t.Fatalf("hosts: %+v", data.Hosts)
	}
	if len(warnings) != 2 || warnings[0].Line != 6 || warnings[0].Field != "src" || warnings[1].Line != 7 || warnings[1].Field != "dest" {
		t.Fatalf("warnings: %+v", warnings)
	}
}

func TestReadStrictRejectsMissingEndpointWithoutPartialData(t *testing.T) {
	input := "src;dest\nfirst;second\n\"multi\nline\";\n"
	data, warnings, err := connections.Read(strings.NewReader(input), connections.StrictMode)
	if err == nil || !strings.Contains(err.Error(), "line 3") || !strings.Contains(err.Error(), "dest") {
		t.Fatalf("want line 3 missing dest error, got %v", err)
	}
	if len(data.Hosts) != 0 || len(data.Connections) != 0 || len(warnings) != 0 {
		t.Fatalf("fatal error returned data %+v or warnings %+v", data, warnings)
	}
}

func TestReadMinimalCSVFixture(t *testing.T) {
	file, err := os.Open("testdata/minimal.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	data, warnings, err := connections.Read(file, connections.DefaultMode)
	if err != nil || len(warnings) != 0 {
		t.Fatalf("read: %v, warnings: %+v", err, warnings)
	}
	wantHosts := []connections.Host{{Name: "Web"}, {Name: "db"}, {Name: "web;01"}, {Name: "line\ntwo"}, {Name: "web"}}
	wantConnections := []connections.DeclaredConnection{
		{Source: "Web", Destination: "db"},
		{Source: "web;01", Destination: "line\ntwo"},
		{Source: "Web", Destination: "db"},
		{Source: "db", Destination: "db"},
		{Source: "web", Destination: "Web"},
	}
	if !slices.Equal(data.Hosts, wantHosts) || !slices.Equal(data.Connections, wantConnections) {
		t.Fatalf("got hosts %+v and declarations %+v; want %+v and %+v", data.Hosts, data.Connections, wantHosts, wantConnections)
	}
}

func TestReadFatalCSVStructureReturnsNoEarlierWarningsOrData(t *testing.T) {
	for _, tc := range []struct {
		name, input, location string
	}{
		{"wrong width, even if empty", "src;dest\nfirst;second\n;db\n;;\n", "line 4"},
		{"malformed quote", "src;dest\nfirst;second\n;db\n\"not closed;dest\n", "line 4"},
		{"space before quote is invalid", "src;dest\nfirst;second\n;db\n \"web\";db\n", "line 4"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data, warnings, err := connections.Read(strings.NewReader(tc.input), connections.DefaultMode)
			if err == nil || !strings.Contains(err.Error(), tc.location) {
				t.Fatalf("want fatal CSV error on %s, got %v", tc.location, err)
			}
			if len(data.Hosts) != 0 || len(data.Connections) != 0 || len(warnings) != 0 {
				t.Fatalf("fatal error returned data %+v or warnings %+v", data, warnings)
			}
		})
	}
}

type failedReader struct{ err error }

func (r failedReader) Read([]byte) (int, error) { return 0, r.err }

func TestReadReaderFailureReturnsNoEarlierWarningsOrData(t *testing.T) {
	wantErr := errors.New("input disconnected")
	r := io.MultiReader(strings.NewReader("src;dest\nfirst;second\n;db\n"), failedReader{wantErr})
	data, warnings, err := connections.Read(r, connections.DefaultMode)
	if !errors.Is(err, wantErr) {
		t.Fatalf("want reader failure, got %v", err)
	}
	if len(data.Hosts) != 0 || len(data.Connections) != 0 || len(warnings) != 0 {
		t.Fatalf("fatal error returned data %+v or warnings %+v", data, warnings)
	}
}

func TestReadRequiresHeaderAndAllowsEmptyGraph(t *testing.T) {
	for _, input := range []string{"", "source;destination\nweb;db\n", "src;dest\nweb;db;extra\n"} {
		data, warnings, err := connections.Read(strings.NewReader(input), connections.DefaultMode)
		if err == nil || len(data.Hosts) != 0 || len(data.Connections) != 0 || len(warnings) != 0 {
			t.Fatalf("input %q: want no data or warnings and fatal error, got %+v, %+v, %v", input, data, warnings, err)
		}
	}
	for _, mode := range []connections.Mode{connections.DefaultMode, connections.StrictMode} {
		data, warnings, err := connections.Read(strings.NewReader("\nsrc;dest\n;\n\n ; \n"), mode)
		if err != nil || len(data.Hosts) != 0 || len(data.Connections) != 0 || len(warnings) != 0 {
			t.Fatalf("mode %d: want empty graph, got %+v, %+v, %v", mode, data, warnings, err)
		}
	}
	data, warnings, err := connections.Read(strings.NewReader("src;dest\n;db\n"), connections.DefaultMode)
	if err != nil || len(data.Hosts) != 0 || len(data.Connections) != 0 || len(warnings) != 1 || warnings[0].Line != 2 {
		t.Fatalf("all records skipped: want empty graph and warning, got %+v, %+v, %v", data, warnings, err)
	}
}

package main

import (
	"strconv"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TestFillTableHeaders(t *testing.T) {
	table := tview.NewTable()
	fillTable(table, nil)

	headers := []string{"Protocol", "Address", "Port", "PID", "Process", "Status"}
	if got := table.GetRowCount(); got != 1 {
		t.Fatalf("expected 1 row, got %d", got)
	}
	if got := table.GetColumnCount(); got != len(headers) {
		t.Fatalf("expected %d columns, got %d", len(headers), got)
	}

	for col, header := range headers {
		cell := table.GetCell(0, col)
		if cell == nil {
			t.Fatalf("missing header cell at column %d", col)
		}
		if cell.Text != header {
			t.Fatalf("expected header %q at column %d, got %q", header, col, cell.Text)
		}
		if !cell.NotSelectable {
			t.Fatalf("expected header cell at column %d to be non-selectable", col)
		}
		_, _, attrs := cell.Style.Decompose()
		if attrs&tcell.AttrBold == 0 {
			t.Fatalf("expected header cell at column %d to be bold", col)
		}
	}
}

func TestFillTableRows(t *testing.T) {
	table := tview.NewTable()
	ports := []portInfo{
		{
			protocol: "tcp",
			address:  "127.0.0.1",
			port:     8080,
			pid:      1234,
			name:     "server",
			status:   "LISTEN",
		},
		{
			protocol: "udp",
			address:  "0.0.0.0",
			port:     5353,
			pid:      5678,
			name:     "mdns",
			status:   "NONE",
		},
	}

	fillTable(table, ports)

	if got := table.GetRowCount(); got != len(ports)+1 {
		t.Fatalf("expected %d rows, got %d", len(ports)+1, got)
	}

	first := ports[0]
	tests := []struct {
		row  int
		col  int
		want string
	}{
		{1, 0, "TCP"},
		{1, 1, first.address},
		{1, 2, strconv.FormatUint(uint64(first.port), 10)},
		{1, 3, strconv.FormatInt(int64(first.pid), 10)},
		{1, 4, first.name},
		{1, 5, first.status},
	}

	for _, tt := range tests {
		cell := table.GetCell(tt.row, tt.col)
		if cell == nil {
			t.Fatalf("missing cell at row %d col %d", tt.row, tt.col)
		}
		if cell.Text != tt.want {
			t.Fatalf("expected cell %d,%d to be %q, got %q", tt.row, tt.col, tt.want, cell.Text)
		}
	}

	portCell := table.GetCell(1, 2)
	if portCell == nil {
		t.Fatal("missing port cell at row 1 col 2")
	}
	fg, _, _ := portCell.Style.Decompose()
	if fg != tcell.ColorLightCyan {
		t.Fatalf("expected port cell color %v, got %v", tcell.ColorLightCyan, fg)
	}
}

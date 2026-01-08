package main

import "testing"

func TestParseSortKey(t *testing.T) {
	tests := []struct {
		input string
		want  portSortKey
	}{
		{input: "", want: sortByProcess},
		{input: "process", want: sortByProcess},
		{input: "PORT", want: sortByPort},
		{input: "pid", want: sortByPID},
		{input: "protocol", want: sortByProtocol},
		{input: "address", want: sortByAddress},
		{input: "status", want: sortByStatus},
	}

	for _, tt := range tests {
		got, err := parseSortKey(tt.input)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tt.input, err)
		}
		if got != tt.want {
			t.Fatalf("expected %q to parse as %v, got %v", tt.input, tt.want, got)
		}
	}

	if _, err := parseSortKey("wat"); err == nil {
		t.Fatal("expected error for invalid sort key")
	}
}

func TestSortPortsByProcess(t *testing.T) {
	ports := []portInfo{
		{name: "beta", port: 200, protocol: "udp", pid: 300},
		{name: "alpha", port: 300, protocol: "tcp", pid: 100},
		{name: "beta", port: 100, protocol: "tcp", pid: 200},
		{name: "", port: 50, protocol: "udp", pid: 400},
	}

	sortPorts(ports, sortByProcess)

	got := []uint32{ports[0].port, ports[1].port, ports[2].port, ports[3].port}
	want := []uint32{300, 100, 200, 50}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected process sort order %v, got %v", want, got)
		}
	}
}

func TestSortPortsByPort(t *testing.T) {
	ports := []portInfo{
		{name: "beta", port: 200, protocol: "udp", pid: 300},
		{name: "alpha", port: 300, protocol: "tcp", pid: 100},
		{name: "beta", port: 100, protocol: "tcp", pid: 200},
		{name: "", port: 50, protocol: "udp", pid: 400},
	}

	sortPorts(ports, sortByPort)

	got := []uint32{ports[0].port, ports[1].port, ports[2].port, ports[3].port}
	want := []uint32{50, 100, 200, 300}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected port sort order %v, got %v", want, got)
		}
	}
}

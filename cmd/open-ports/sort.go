package main

import (
	"fmt"
	"sort"
	"strings"
)

type portSortKey int

const (
	sortByProcess portSortKey = iota
	sortByPort
	sortByPID
	sortByProtocol
	sortByAddress
	sortByStatus
)

const (
	defaultSortKey   = sortByProcess
	defaultSortLabel = "process"
	validSortKeys    = "process, port, pid, protocol, address, status"
)

func parseSortKey(value string) (portSortKey, error) {
	trimmed := strings.ToLower(strings.TrimSpace(value))
	if trimmed == "" {
		return defaultSortKey, nil
	}
	switch trimmed {
	case defaultSortLabel:
		return sortByProcess, nil
	case "port":
		return sortByPort, nil
	case "pid":
		return sortByPID, nil
	case "protocol":
		return sortByProtocol, nil
	case "address":
		return sortByAddress, nil
	case "status":
		return sortByStatus, nil
	default:
		return defaultSortKey, fmt.Errorf("invalid sort key %q (valid: %s)", value, validSortKeys)
	}
}

func sortKeyLabel(key portSortKey) string {
	switch key {
	case sortByProcess:
		return "process"
	case sortByPort:
		return "port"
	case sortByPID:
		return "pid"
	case sortByProtocol:
		return "protocol"
	case sortByAddress:
		return "address"
	case sortByStatus:
		return "status"
	default:
		return defaultSortLabel
	}
}

func sortPorts(ports []portInfo, key portSortKey) {
	sort.SliceStable(ports, func(i, j int) bool {
		a, b := ports[i], ports[j]
		switch key {
		case sortByProcess:
			if cmp := compareString(a.name, b.name, true); cmp != 0 {
				return cmp < 0
			}
		case sortByPort:
			if a.port != b.port {
				return a.port < b.port
			}
		case sortByPID:
			if a.pid != b.pid {
				return a.pid < b.pid
			}
		case sortByProtocol:
			if cmp := compareString(a.protocol, b.protocol, false); cmp != 0 {
				return cmp < 0
			}
		case sortByAddress:
			if cmp := compareString(a.address, b.address, false); cmp != 0 {
				return cmp < 0
			}
		case sortByStatus:
			if cmp := compareString(a.status, b.status, false); cmp != 0 {
				return cmp < 0
			}
		}
		return tieBreak(a, b)
	})
}

func tieBreak(a, b portInfo) bool {
	if a.port != b.port {
		return a.port < b.port
	}
	if cmp := compareString(a.protocol, b.protocol, false); cmp != 0 {
		return cmp < 0
	}
	if a.pid != b.pid {
		return a.pid < b.pid
	}
	if cmp := compareString(a.address, b.address, false); cmp != 0 {
		return cmp < 0
	}
	if cmp := compareString(a.name, b.name, true); cmp != 0 {
		return cmp < 0
	}
	return strings.ToLower(a.status) < strings.ToLower(b.status)
}

func compareString(a, b string, emptyLast bool) int {
	na := strings.ToLower(a)
	nb := strings.ToLower(b)
	if na == nb {
		return 0
	}
	if emptyLast {
		if na == "" {
			return 1
		}
		if nb == "" {
			return -1
		}
	}
	if na < nb {
		return -1
	}
	return 1
}

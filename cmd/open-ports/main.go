package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

type portInfo struct {
	protocol string
	address  string
	port     uint32
	pid      int32
	name     string
	status   string
}

func loadPorts() ([]portInfo, error) {
	conns, err := net.Connections("inet")
	if err != nil {
		return nil, err
	}

	ports := make([]portInfo, 0, len(conns))
	for _, conn := range conns {
		status := strings.ToUpper(conn.Status)
		if status != "LISTEN" && status != "NONE" {
			continue
		}

		proto := "unknown"
		switch conn.Type {
		case syscall.SOCK_STREAM:
			proto = "tcp"
		case syscall.SOCK_DGRAM:
			proto = "udp"
		}

		name := ""
		if conn.Pid > 0 {
			proc, err := process.NewProcess(conn.Pid)
			if err == nil {
				if pname, err := proc.Name(); err == nil {
					name = pname
				}
			}
		}

		ports = append(ports, portInfo{
			protocol: proto,
			address:  conn.Laddr.IP,
			port:     conn.Laddr.Port,
			pid:      conn.Pid,
			name:     name,
			status:   status,
		})
	}

	sort.Slice(ports, func(i, j int) bool {
		if ports[i].port == ports[j].port {
			return ports[i].protocol < ports[j].protocol
		}
		return ports[i].port < ports[j].port
	})

	return ports, nil
}

func fillTable(table *tview.Table, ports []portInfo) {
	headers := []string{"Protocol", "Address", "Port", "PID", "Process", "Status"}
	for col, header := range headers {
		cell := tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false).
			SetAttributes(tcell.AttrBold)
		table.SetCell(0, col, cell)
	}

	for row, info := range ports {
		values := []string{
			strings.ToUpper(info.protocol),
			info.address,
			strconv.FormatUint(uint64(info.port), 10),
			strconv.FormatInt(int64(info.pid), 10),
			info.name,
			info.status,
		}
		for col, value := range values {
			cell := tview.NewTableCell(value)
			if col == 2 {
				cell.SetTextColor(tcell.ColorLightCyan)
			}
			table.SetCell(row+1, col, cell)
		}
	}
}

func main() {
	app := tview.NewApplication()
	table := tview.NewTable().
		SetFixed(1, 0).
		SetSelectable(true, false)

	statusBar := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)

	var ports []portInfo

	refresh := func() {
		table.Clear()
		var err error
		ports, err = loadPorts()
		if err != nil {
			statusBar.SetText(fmt.Sprintf("[red]Failed to load ports: %v", err))
			fillTable(table, nil)
			return
		}
		statusBar.SetText(fmt.Sprintf("[green]Found %d open ports. [white](r)efresh  (k)ill  (q)uit", len(ports)))
		fillTable(table, ports)
	}

	refresh()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q':
			app.Stop()
			return nil
		case 'r':
			refresh()
			return nil
		case 'k':
			row, _ := table.GetSelection()
			if row <= 0 || row-1 >= len(ports) {
				statusBar.SetText("[yellow]Select a port to kill.")
				return nil
			}
			selected := ports[row-1]
			if selected.pid <= 0 {
				statusBar.SetText("[yellow]Selected port has no PID to kill.")
				return nil
			}
			if err := syscall.Kill(int(selected.pid), syscall.SIGTERM); err != nil {
				statusBar.SetText(fmt.Sprintf("[red]Failed to kill PID %d: %v", selected.pid, err))
				return nil
			}
			statusBar.SetText(fmt.Sprintf("[green]Sent SIGTERM to PID %d on port %d.", selected.pid, selected.port))
			refresh()
			return nil
		}
		return event
	})

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(table, 0, 1, true).
		AddItem(statusBar, 1, 0, false)

	if err := app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

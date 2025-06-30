package views

import (
	"fmt"
	"strings"

	"github.com/admiller/ltop/internal/models"
	"github.com/admiller/ltop/internal/ui/styles"
	"github.com/admiller/ltop/pkg/utils"
)

type NetworkView struct{}

func NewNetworkView() *NetworkView {
	return &NetworkView{}
}

func (nv *NetworkView) Render(snapshot *models.MetricsSnapshot, width, height int) string {
	if snapshot == nil {
		return "No data available"
	}

	var sections []string

	sections = append(sections, nv.renderNetworkInterfaces(snapshot))

	content := strings.Join(sections, "\n\n")
	return styles.Panel().Width(width).Render(content)
}

func (nv *NetworkView) renderNetworkInterfaces(snapshot *models.MetricsSnapshot) string {
	var network []string
	network = append(network, styles.Title().Render("Network Interfaces"))

	if len(snapshot.Network.Interfaces) == 0 {
		network = append(network, styles.Muted().Render("No network interfaces found"))
		return strings.Join(network, "\n")
	}

	headers := []string{"Interface", "State", "RX Rate", "TX Rate", "RX Bytes", "TX Bytes", "RX Errors", "TX Errors"}
	network = append(network, nv.renderNetworkHeader(headers))

	for _, iface := range snapshot.Network.Interfaces {
		row := nv.renderNetworkRow(iface)
		network = append(network, row)
	}

	return strings.Join(network, "\n")
}

func (nv *NetworkView) renderNetworkHeader(headers []string) string {
	var parts []string
	widths := []int{12, 8, 12, 12, 12, 12, 10, 10}

	for i, header := range headers {
		if i < len(widths) {
			text := utils.PadString(header, widths[i], ' ')
			parts = append(parts, styles.TableHeader().Render(text))
		}
	}

	return strings.Join(parts, " ")
}

func (nv *NetworkView) renderNetworkRow(iface models.NetworkInterface) string {
	var parts []string
	widths := []int{12, 8, 12, 12, 12, 12, 10, 10}

	name := utils.TruncateString(iface.Name, widths[0])
	name = utils.PadString(name, widths[0], ' ')
	parts = append(parts, styles.TableRow().Render(name))

	state := utils.PadString(iface.State, widths[1], ' ')
	stateStyle := styles.TableRow()
	if iface.State == "up" {
		stateStyle = styles.Success()
	} else if iface.State == "down" {
		stateStyle = styles.Muted()
	}
	parts = append(parts, stateStyle.Render(state))

	rxRate := utils.PadString(utils.FormatBytesPerSecond(iface.RecvBytesPerSec), widths[2], ' ')
	parts = append(parts, styles.TableRow().Render(rxRate))

	txRate := utils.PadString(utils.FormatBytesPerSecond(iface.SentBytesPerSec), widths[3], ' ')
	parts = append(parts, styles.TableRow().Render(txRate))

	rxBytes := utils.PadString(utils.FormatBytes(iface.BytesRecv), widths[4], ' ')
	parts = append(parts, styles.TableRow().Render(rxBytes))

	txBytes := utils.PadString(utils.FormatBytes(iface.BytesSent), widths[5], ' ')
	parts = append(parts, styles.TableRow().Render(txBytes))

	rxErrors := utils.PadString(fmt.Sprintf("%d", iface.ErrorsRecv), widths[6], ' ')
	errStyle := styles.TableRow()
	if iface.ErrorsRecv > 0 {
		errStyle = styles.Warning()
	}
	parts = append(parts, errStyle.Render(rxErrors))

	txErrors := utils.PadString(fmt.Sprintf("%d", iface.ErrorsSent), widths[7], ' ')
	errStyle = styles.TableRow()
	if iface.ErrorsSent > 0 {
		errStyle = styles.Warning()
	}
	parts = append(parts, errStyle.Render(txErrors))

	return strings.Join(parts, " ")
}
package gather

import (
	"errors"
	"strconv"
)

// USChannelHtml += "<td class='moto-param-header-s'>Lock Status</td>";
// USChannelHtml += "<td class='moto-param-header-s'>Channel Type</td>";
// USChannelHtml += "<td class='moto-param-header-s'>Channel ID</td>";
// USChannelHtml += "<td class='moto-param-header-s'>Symb. Rate (Ksym/sec)</td>";
// USChannelHtml += "<td class='moto-param-header-s'>Freq. (MHz)</td>";
// USChannelHtml += "<td class='moto-param-header-s'>Pwr (dBmV)</td></tr>";

type UpstreamInfo struct {
	ID            int64
	LockStatus    string
	ChannelType   string
	Channel       int64
	SymbolRate    int64
	Frequency     float64
	DecibelMillivolts float64
}

func (info *UpstreamInfo) Parse(row []string) error {
	const upstreamDataSize = 6
	if len(row) != upstreamDataSize {
		return errors.New("invalid data size")
	}

	var err error

	info.ID, err = strconv.ParseInt(row[0], 10, 64)
	if err != nil {
		return err
	}

	info.LockStatus = row[1]
	info.ChannelType = row[2]

	info.Channel, err = strconv.ParseInt(row[3], 10, 64)

	info.SymbolRate, err = strconv.ParseInt(row[4], 10, 64)
	if err != nil {
		return err
	}

	info.Frequency, err = strconv.ParseFloat(row[5], 64)
	if err != nil {
		return err
	}

	info.DecibelMillivolts, err = strconv.ParseFloat(row[6], 64)
	if err != nil {
		return err
	}

	return nil
}

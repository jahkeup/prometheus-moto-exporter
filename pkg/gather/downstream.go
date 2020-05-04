package gather

import (
	"strings"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// var DSChannelHtml = "<tr align='center'><td class='moto-param-header-s'>&nbsp;&nbsp;&nbsp;Channel</td>";
//
// DSChannelHtml += "<td class='moto-param-header-s'>Lock Status</td>";
// DSChannelHtml += "<td class='moto-param-header-s'>Modulation</td>";
// DSChannelHtml += "<td class='moto-param-header-s'>Channel ID</td>";
//
// DSChannelHtml += "<td class='moto-param-header-s'>Freq. (MHz)</td>";
// DSChannelHtml += "<td class='moto-param-header-s'>Pwr (dBmV)</td>";
// DSChannelHtml += "<td class='moto-param-header-s'>SNR (dB)</td>";
//
// DSChannelHtml += "<td class='moto-param-header-s'>Corrected</td>";
// DSChannelHtml += "<td class='moto-param-header-s'>Uncorrected</td></tr>";

type DownstreamInfo struct {
	ID int64
	Status string
	Modulation string
	ChannelID int64
	Frequency float64
	DecibelMillivolts float64
	Signal float64
	Corrected int64
	Uncorrected int64
}

func (info *DownstreamInfo) Parse(row []string) error {
	const infoRowSize = 10
	const (
		idField = iota
		statusField
		modulationField
		channelIDField
		frequencyField
		dbmvField
		signalField
		correctedField
		uncorrectedField

	)
	if len(row) != infoRowSize {
		return errors.New("invalid data size")
	}

	getField := func(i int) string {
		return strings.TrimSpace(row[i])
	}

	var err error

	info.ID, err = strconv.ParseInt(getField(idField), 10, 64)
	if err != nil {
		return errors.Wrap(err, "unable to parse upstream ID")
	}

	info.Status = row[statusField]
	info.Modulation = row[modulationField]

	logrus.WithField("row", row).WithField("info", info).Debugf("parsing channel ID from: %q", row[2])
	info.ChannelID, err = strconv.ParseInt(getField(channelIDField), 10, 64)
	if err != nil {
		return errors.Wrap(err, "parse channel ID")
	}

	info.Frequency, err = strconv.ParseFloat(getField(frequencyField), 64)
	if err != nil {
		return errors.Wrap(err, "parse frequency")
	}

	info.DecibelMillivolts, err = strconv.ParseFloat(getField(dbmvField), 64)
	if err != nil {
		return errors.Wrap(err, "parse dBmV")
	}

	info.Signal, err = strconv.ParseFloat(getField(signalField), 64)
	if err != nil {
		return errors.Wrap(err, "parse signal")
	}

	info.Corrected, err = strconv.ParseInt(getField(correctedField), 10, 64)
	if err != nil {
		return errors.Wrap(err, "parse corrected counter")
	}

	info.Uncorrected, err = strconv.ParseInt(getField(uncorrectedField), 10, 64)
	if err != nil {
		return errors.Wrap(err, "parse uncorrected counter")
	}

	return nil
}

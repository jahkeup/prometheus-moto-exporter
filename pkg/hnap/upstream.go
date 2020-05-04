package hnap

import (
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/jahkeup/prometheus-moto-exporter/pkg/plustable"
)

type UpstreamInfo struct {
	ID                int64
	LockStatus        string
	Modulation        string
	Channel           int64
	SymbolRate        int64
	Frequency         float64
	DecibelMillivolts float64
}

func (info *UpstreamInfo) Parse(row []string) error {
	const dataSize = 8 // 1 unused field, not sure what it is!

	const (
		idField = iota
		lockStatusField
		modulationField
		channelIDField
		symbolRateField
		frequencyField
		dbmvField
	)

	if len(row) != dataSize {
		return errors.Errorf("invalid data size: expected %d but found %d", dataSize, len(row))
	}

	var err error

	info.ID, err = strconv.ParseInt(row[idField], 10, 64)
	if err != nil {
		return errors.Wrap(err, "parse ID")
	}

	info.LockStatus = row[lockStatusField]
	info.Modulation = row[modulationField]

	info.Channel, err = strconv.ParseInt(row[channelIDField], 10, 64)
	if err != nil {
		return errors.Wrap(err, "parse channel ID")
	}

	info.SymbolRate, err = strconv.ParseInt(row[symbolRateField], 10, 64)
	if err != nil {
		return errors.Wrap(err, "parse symbol rate")
	}
	// ksym -> sym
	info.SymbolRate *= 1000

	info.Frequency, err = strconv.ParseFloat(row[frequencyField], 64)
	if err != nil {
		return errors.Wrap(err, "parse frequency")
	}
	info.Frequency *= 1000 * 1000 // Mhz -> hz

	info.DecibelMillivolts, err = strconv.ParseFloat(row[dbmvField], 64)
	if err != nil {
		return errors.Wrap(err, "parse dBmV")
	}

	return nil
}

type UpstreamChannelResponse struct {
	Channels []UpstreamInfo
}

func (r *UpstreamChannelResponse) UnmarshalJSON(data []byte) error {
	var innerType struct {
		MotoConnUpstreamChannel string
	}

	err := json.Unmarshal(data, &innerType)
	if err != nil {
		return err
	}

	tbl := plustable.Parse(innerType.MotoConnUpstreamChannel)
	info := make([]UpstreamInfo, len(tbl))
	for i, row := range tbl {
		err = info[i].Parse(row)
		if err != nil {
			logrus.WithError(err).WithField("row", row).Error("could not parse data")
			return err
		}
	}

	r.Channels = info

	return nil
}

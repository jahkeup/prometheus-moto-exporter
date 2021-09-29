package main

import (
	"fmt"
	"os"
	"net/url"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	"github.com/prometheus/common/expfmt"

	"github.com/jahkeup/prometheus-moto-exporter/pkg/gather"
)

func NewCheckCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "check",
		Short: "Run a check run against the configured endpoint",
		SilenceUsage: true,
	}

	var (
		endpointURL *url.URL
	)

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		endpoint, err := cmd.Flags().GetString("endpoint")
		if err != nil {
			return err
		}
		parsedEndpoint, err := url.Parse(endpoint)
		if err != nil {
			return err
		}
		endpointURL = parsedEndpoint

		return nil
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			return err
		}
		password, err := cmd.Flags().GetString("password")
		if err != nil {
			return err
		}
		gatherer, err := gather.New(endpointURL, username, password)
		if err != nil {
			return err
		}

		srv, err := NewServer(gatherer)
		if err != nil {
			return err
		}

		reg := prometheus.NewRegistry()
		err = srv.RegisterMetrics(reg)
		if err != nil {
			return fmt.Errorf("failed to register exporter metrics: %w", err)
		}
		err = srv.Collect()
		if err != nil {
			return fmt.Errorf("unable to get metrics from endpoint: %w", err)
		}

		// Gather metrics and dump to console.

		mfs, err := reg.Gather()
		if err != nil {
			return err
		}
		for _, mf := range mfs {
			if _, err := expfmt.MetricFamilyToText(os.Stdout, mf); err != nil {
				return err
			}
		}

		return nil
	}

	return cmd
}

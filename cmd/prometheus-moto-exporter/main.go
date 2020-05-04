package main

import (
	"net/url"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	"github.com/jahkeup/prometheus-moto-exporter/pkg/gather"
)

// TODO: read these in! Maybe use viper?
const (
	envEndpoint = "MOTO_ENDPOINT"
	envUsername = "MOTO_USERNAME"
	envPassword = "MOTO_PASSWORD"
)

func main() {
	if err := App().Execute(); err != nil {
		os.Exit(1)
	}
}

func App() *cobra.Command {

	var (
		logDebug bool
		bindAddr string
		endpoint string
		username string
		password string
	)

	cmd := &cobra.Command{
		Use: "prometheus-moto-exporter",
		Short: "Exporter for Motorola modems equipped with HNAP",
		// TODO: move to a build-time var
		Version: "0.1.0",
		// Don't print usage on run errors.
		SilenceUsage: true,
	}
	cmd.AddCommand(NewCheckCommand())

	cmd.Flags().StringVar(&bindAddr, "bind", "127.0.0.1:9731", "http server bind address")

	cmd.PersistentFlags().StringVar(&endpoint, "endpoint", "https://192.168.100.1/HNAP1/", "modem HNAP endpoint")
	cmd.PersistentFlags().StringVar(&username, "username", "admin", "modem HNAP username")
	cmd.PersistentFlags().StringVar(&password, "password", "motorola", "modem HNAP password")

	cmd.PersistentFlags().BoolVar(&logDebug, "debug", false, "enable debug logging")

	var (
		endpointURL *url.URL
	)

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		logrus.SetLevel(logrus.InfoLevel)
		if logDebug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		parsedEndpoint, err := url.Parse(endpoint)
		if err != nil {
			return err
		}
		endpointURL = parsedEndpoint

		return nil
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		logrus.Debugf("Pulling HNAP metrics from %s", endpointURL.String())

		gatherer, err := gather.New(endpointURL, username, password)
		if err != nil {
			return err
		}

		server, err := NewServer(gatherer)
		if err != nil {
			logrus.WithError(err).Error("unable to setup server")
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())

		sigsent := make(chan os.Signal, 1)
		signal.Notify(sigsent, os.Interrupt)

		go func() {
			<-sigsent
			logrus.Info("SIGINT: shutting down server")
			cancel()
		}()

		err = server.Run(ctx, bindAddr)
		if err != nil {
			logrus.WithError(err).Error("server error")
			return err
		}

		logrus.Info("done")

		return nil
	}

	return cmd
}

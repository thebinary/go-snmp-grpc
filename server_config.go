package main

import (
	"fmt"
	"log"
	"os"

	"github.com/soniah/gosnmp"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/thebinary/go-snmp-grpc/server"
)

const (
	defaultAddr             = ":8161"
	defaultMetricsEnabled   = false
	defaultMetricsAddr      = ":9161"
	defaultMetricsPath      = "/metrics"
	varDefaultSnmpCommunity = "default_comm"
	varDefaultSnmpVersion   = "default_version"
	varLogSnmp              = "log_snmp"
	varLogSnmpPrefix        = "log_snmp_prefix"
	varMetricsEnable        = "metrics_enable"
	varMetricsAddr          = "metrics_addr"
	varMetricsPath          = "metrics_path"
)

// Config defines the instance configuration
type Config struct {
	MetricsEnabled bool
	MetricsAddr    string
	MetricsPath    string
}

var commandSrv = server.CommandServer{}
var config = Config{}

func init() {
	// bind ADRR env var without the prefix
	err := viper.BindEnv("addr", "ADDR")
	if err != nil {
		log.Fatalf("bindenv '%s' error: %v", "addr", err)
	}

	// bind other viper env vars
	viper.SetEnvPrefix("GSNMP")
	envs := []string{
		varDefaultSnmpCommunity, varDefaultSnmpVersion, varLogSnmp, varLogSnmpPrefix,
		varMetricsEnable, varMetricsAddr, varMetricsPath,
	}
	for _, envName := range envs {
		err := viper.BindEnv(envName)
		if err != nil {
			log.Fatalf("Binding env '%s' error: %v", envName, err)
		}
	}

	// add command line flags
	flag.String("addr", defaultAddr, "Listening Address")
	flag.String(varDefaultSnmpCommunity, gosnmp.Default.Community, "Default SNMP Community")
	flag.String(varDefaultSnmpVersion, fmt.Sprint(gosnmp.Default.Version), "Default SNMP Version")
	flag.Bool(varLogSnmp, false, "SNMP Debug logging")
	flag.String(varLogSnmpPrefix, "[SNMP]\t", "SNMP Debug logging prefix")
	flag.Bool(varMetricsEnable, defaultMetricsEnabled, "Enable Prometheus Metrics")
	flag.String(varMetricsAddr, defaultMetricsAddr, "Prometheus Metrics Listening Address")
	flag.String(varMetricsPath, defaultMetricsPath, "Prometheus Metrics Path")

	// parse command line flags
	flag.Parse()

	// bind command line flags to viper config
	err = viper.BindPFlags(flag.CommandLine)
	if err != nil {
		log.Fatalf("Binding command flags err: %v", err)
	}

	// init server with given configurations
	commandSrv.ListenAddr = viper.GetString("addr")
	commandSrv.DefaultCommunity = viper.GetString(varDefaultSnmpCommunity)
	defSnmpVersion := viper.GetString(varDefaultSnmpVersion)

	switch defSnmpVersion {
	case "1":
		commandSrv.DefaultVersion = gosnmp.Version1
	case "2", "2c":
		commandSrv.DefaultVersion = gosnmp.Version2c
	case "3":
		log.Fatal("snmp-version 3 not supported yet")
	default:
		log.Fatalf("invalid snmp-version; %s", defSnmpVersion)
	}

	if viper.GetBool(varLogSnmp) {
		commandSrv.SNMPLogger = log.New(os.Stdout, viper.GetString(varLogSnmpPrefix), log.LstdFlags)
	}

	config.MetricsEnabled = viper.GetBool(varMetricsEnable)
	config.MetricsAddr = viper.GetString(varMetricsAddr)
	config.MetricsPath = viper.GetString(varMetricsPath)

	/*
		// Print configuration
		fmt.Println("============================================================")
		fmt.Println("CONFIG")
		fmt.Println("Listening Address =", commandSrv.ListenAddr)
		fmt.Println("Default Version =", commandSrv.DefaultVersion)
		fmt.Println("Default Community =", commandSrv.DefaultCommunity)
		fmt.Println("============================================================")
		fmt.Println()
	*/
}

package main

import (
	"fmt"
	"log"

	"github.com/soniah/gosnmp"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultAddr             = ":8161"
	varDefaultSnmpCommunity = "default_comm"
	varDefaultSnmpVersion   = "default_version"
)

type CommandServer struct {
	ListenAddr       string
	DefaultVersion   gosnmp.SnmpVersion
	DefaultCommunity string
}

// C : Global Application Config
var Server = CommandServer{}

func init() {
	// bind ADRR env var without the prefix
	err := viper.BindEnv("addr", "ADDR")
	if err != nil {
		log.Fatalf("bindenv '%s' error: %v", "addr", err)
	}

	// bind other viper env vars
	viper.SetEnvPrefix("GSNMP")
	envs := []string{
		varDefaultSnmpCommunity, varDefaultSnmpVersion,
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

	// parse command line flags
	flag.Parse()

	// bind command line flags to viper config
	err = viper.BindPFlags(flag.CommandLine)
	if err != nil {
		log.Fatalf("Binding command flags err: %v", err)
	}

	// init server with given configurations
	Server.ListenAddr = viper.GetString("addr")
	Server.DefaultCommunity = viper.GetString(varDefaultSnmpCommunity)
	defSnmpVersion := viper.GetString(varDefaultSnmpVersion)

	switch defSnmpVersion {
	case "1":
		Server.DefaultVersion = gosnmp.Version1
	case "2", "2c":
		Server.DefaultVersion = gosnmp.Version2c
	case "3":
		log.Fatal("snmp-version 3 not supported yet")
	default:
		log.Fatalf("invalid snmp-version; %s", defSnmpVersion)
	}

	/*
		// Print configuration
		fmt.Println("============================================================")
		fmt.Println("CONFIG")
		fmt.Println("Listening Address =", Server.ListenAddr)
		fmt.Println("Default Version =", Server.DefaultVersion)
		fmt.Println("Default Community =", Server.DefaultCommunity)
		fmt.Println("============================================================")
		fmt.Println()
	*/
}

// Package config is used to define any configuration that isn't passed in from the command line
// or is default options that can be overridden
package config

import (
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/milligan22963/afmlog"
	"gopkg.in/yaml.v2"
)

const (
	identificationPath = "/var/cache/afm/identifier.id"
	defaultSQLPort     = 3306
	defaultSQLType     = "mysql"
	defaultMQTTPort    = 1883
	defaultFileOptions = 0600
	defaultWebPort     = 8080
	mysqlDBType        = "mysql"
	sqlliteDBType      = "sqlite"
	postGRESDBType     = "postgres"
)

type BrokerSettings struct {
	Address        string `yaml:"br_address"`
	Port           int    `yaml:"br_port"`
	UseSSL         bool   `yaml:"br_use_ssl"`
	PrivateKeyPath string `yaml:"br_private_key_path"`
	PublicKeyPath  string `yaml:"br_public_key_path"`
	CAPath         string `yaml:"br_ca_path"`
}

type DatabaseSettings struct {
	Name     string `yaml:"db_name"`
	User     string `yaml:"db_user"`
	Password string `yaml:"db_password"`
	Host     string `yaml:"db_host"`
	Port     int    `yaml:"db_port"`
	Type     string `yaml:"type"`
}

type WebServerSettings struct {
	Host     string `yaml:"ws_host"`
	FileRoot string `yaml:"ws_root"`
	Port     int    `yaml:"ws_port"`
}

type CameraConfiguration struct {
	BrokerSettings    BrokerSettings       `yaml:"broker"`
	DatabaseSettings  DatabaseSettings     `yaml:"database"`
	WebServerSettings WebServerSettings    `yaml:"website"`
	LogSettings       afmlog.Configuration `yaml:"log"`
}

// DefaultConfigPath to our default config
const DefaultConfigPath = "/etc/afm/camera.yaml"

// DatabaseQueryResponse struct representing the response to a given query
type DatabaseQueryResponse struct {
	Query    string
	Response *sqlx.Rows
	Err      error
}

// AppConfiguration is configuration
type AppConfiguration struct {
	CameraConfiguration CameraConfiguration
	AppActive           chan struct{}
	IncomingMQTT        chan [2]string
	OutgoingMQTT        chan [3]string
	ClientID            string
	Database            *sqlx.DB
}

func determineCurrentIP() (string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		// = GET LOCAL IP ADDRESS
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("unable to find ip address")
}

func determineCurrentNetworkHardwareInterface(currentIP string) (string, error) {
	// get all the system's or local machine's network interfaces
	interfaces, interfaceerr := net.Interfaces()
	for _, interf := range interfaces {

		if addrs, err := interf.Addrs(); err == nil {
			for _, addr := range addrs {
				// only interested in the name with current IP address
				if strings.Contains(addr.String(), currentIP) {
					return interf.Name, nil
				}
			}
		}
	}
	return "", interfaceerr
}

func determineDeviceMACAddress() (string, error) {

	currentIP, err := determineCurrentIP()
	if err != nil {
		return "", err
	}

	hardwareInterfaceName, err := determineCurrentNetworkHardwareInterface(currentIP)

	if err != nil {
		return "", err
	}

	// extract the hardware information base on the interface name
	// capture above
	netInterface, err := net.InterfaceByName(hardwareInterfaceName)

	if err != nil {
		return "", err
	}

	macAddress := netInterface.HardwareAddr

	// verify if the MAC address can be parsed properly
	hwAddr, err := net.ParseMAC(macAddress.String())

	if err != nil {
		return "", err
	}

	return hwAddr.String(), nil
}

func determineDeviceSerialNumber() (string, error) {
	out, err := exec.Command("/usr/sbin/dmidecode", "-t", "system").Output()
	if err == nil {
		return "", err
	}
	for _, l := range strings.Split(string(out), "\n") {
		if strings.Contains(l, "Serial Number") {
			s := strings.Split(l, ":")
			return s[len(s)-1], nil
		}
	}
	return "", fmt.Errorf("unable to find serial number")
}

func determineDeviceClientID() string {
	// Try reading well known file location
	identifier, err := ioutil.ReadFile(identificationPath)

	if err == nil {
		cleansedIdentifier := strings.TrimRight(string(identifier), "\n")
		return cleansedIdentifier
	}

	// Try reading hardware assigned serial number
	serialNumber, err := determineDeviceSerialNumber()
	if err == nil {
		return serialNumber
	}

	// Default back to mac address
	macAddress, err := determineDeviceMACAddress()
	if err == nil {
		return macAddress
	}

	return "id_failure"
}

func (configuration *CameraConfiguration) LoadConfiguration(filename string) error {
	fileContents, err := ioutil.ReadFile(filepath.Clean(filename))

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(fileContents, configuration)
	if err != nil {
		return err
	}
	return configuration.LogSettings.LoadConfiguration()
}

func (configuration *CameraConfiguration) SetupDatabase(initialDBNameConnect bool) *sqlx.DB {
	connectionString := configuration.DatabaseSettings.User + ":" + configuration.DatabaseSettings.Password
	connectionString += "@tcp(" + configuration.DatabaseSettings.Host + ":" + strconv.Itoa(configuration.DatabaseSettings.Port) + ")/"

	if initialDBNameConnect {
		connectionString += configuration.DatabaseSettings.Name
	}

	database, err := sqlx.Open(configuration.DatabaseSettings.Type, connectionString)

	if err != nil {
		panic(err.Error())
	}

	return database
}

func (appConfig *AppConfiguration) GetLogger() *afmlog.Log {
	return appConfig.CameraConfiguration.LogSettings.UserLog
}

// NewSiteConfiguration creates an instance of the site configuration struct
func NewSiteConfiguration(configFile string, initialDBNameConnect bool) *AppConfiguration {
	appConfig := &AppConfiguration{
		CameraConfiguration: CameraConfiguration{},
		AppActive:           make(chan struct{}),
		IncomingMQTT:        make(chan [2]string),
		OutgoingMQTT:        make(chan [3]string),
		ClientID:            determineDeviceClientID(),
		Database:            nil,
	}

	//	appConfig.Database = appConfig.CameraConfiguration.SetupDatabase(initialDBNameConnect)
	err := appConfig.CameraConfiguration.LoadConfiguration(configFile)
	if err != nil {
		panic(err)
	}

	return appConfig
}

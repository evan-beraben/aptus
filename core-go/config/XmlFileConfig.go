package config

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type ServerConfig struct{
	XMLName 		xml.Name `xml:"serverConfig"`
	DbConnection	DbConnectionProperties `xml:"server-db"`
	MqServer		MQServerProperties `xml:"mq-server"`
}

type DbConnectionProperties struct{
	XMLName 			xml.Name `xml:"server-db"`
	DriverClass			string `xml:"driverClass"`
	ConnectionString 	string	`xml:"connectionString"`
	Username			string	`xml:"username"`
	Password			string	`xml:"password"`
}

type MQServerProperties struct{
	XMLName 	xml.Name 	`xml:"mq-server"`
	MqUrl		string 		`xml:"mqurl,attr"`
	Username	string  	`xml:"username,attr"`
	Password	string		`xml:"password,attr"`
}


func GetConfig() (*ServerConfig){

	file,err := os.Open("./conf/server.xml")
	if err != nil{
		return nil
	}

	defer file.Close()

	bytes,err := ioutil.ReadAll(file)

	var config ServerConfig
	xml.Unmarshal(bytes, &config)

	return &config

}
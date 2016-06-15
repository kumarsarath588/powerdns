package main

import (
	"encoding/json"
	"net/http"
	"strings"

	pdns "github.com/hashicorp/terraform/builtin/providers/powerdns"
)

//Config is struct for ServerUrl and ApiKey
type Config struct {
	ServerURL string
	APIKey    string
	Zone      string
	Oper      string
	Inputs    Inputs
}

//Inputs need to given to powerdns
type Inputs struct {
	Name    string
	Content string
	Type    string
}

//PdnsWebhookHomePage homepage of powerdns
func PdnsWebhookHomePage(w http.ResponseWriter, r *http.Request) {
	config := new(Config)
	err := json.NewDecoder(r.Body).Decode(config)
	if err != nil {
		panic(err)
	}
	client, err := NewClient(config)
	if err != nil {
		panic(err)
	}
	if strings.ToLower(config.Oper) == "create" {
		CreateRecord(client, config.Zone, config.Inputs.Name, config.Inputs.Content, config.Inputs.Type)
	}
	if strings.ToLower(config.Oper) == "delete" {
		DeleteRecord(client, config.Zone, config.Inputs.Name, config.Inputs.Type)
	}
}

//CreateRecord is used to create new records in pdns
func CreateRecord(c *pdns.Client, Zone string, Name string, IPaddress string, Type string) (string, error) {
	exists, err := c.RecordExists(Zone, Name, Type)
	if exists == true {
		return "record already exists: " + Name, err
	}
	rec := pdns.Record{
		Name:     Name,
		Type:     Type,
		Content:  IPaddress,
		TTL:      1800,
		Disabled: false,
	}
	status, err := c.CreateRecord(Zone, rec)
	return status, err
}

//DeleteRecord is used to delete record in pdns
func DeleteRecord(c *pdns.Client, Zone string, Name string, Type string) (string, error) {
	if exists, err := c.RecordExists(Zone, Name, Type); exists != true {
		return "record Dosen't exists: " + Name, err
	}
	err := c.DeleteRecordSet(Zone, Name, Type)
	if err != nil {
		return "Error Deleting Record: " + Name, err
	}
	return "Record Deleted Sucessfully: " + Name, err
}

//NewClient is used to create new connection to powerdns server
func NewClient(config *Config) (*pdns.Client, error) {
	client, err := pdns.NewClient(config.ServerURL, config.APIKey)
	if err != nil {
		return client, err
	}
	return client, err
}

func main() {
	http.HandleFunc("/", PdnsWebhookHomePage)
	http.ListenAndServe(":3001", nil)
}


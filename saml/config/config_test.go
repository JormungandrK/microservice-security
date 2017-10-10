package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config := `{
	    "services": {
	    	"microservice-registration": "https://127.0.0.1:8083/users",
	    	"microservice-user": "http://127.0.0.1:8081/users"
	    }
	  }`

	cnfFile, err := ioutil.TempFile("", "tmp-config")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(cnfFile.Name())

	cnfFile.WriteString(config)

	cnfFile.Sync()

	_, err = LoadConfig("not-exists.json")
	if err == nil {
		t.Fatal("Nil error for invalid config file")
	}

	loadedCnf, err := LoadConfig(cnfFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if loadedCnf == nil {
		t.Fatal("Configuration was not read")
	}
}

package main

import (
	"log"
	"os"
	"text/template"
)

type Inventory struct {
	UserName   string
	Host       string
	Port       int
	PrivateKey string
}

func generateInventory(inventory Inventory) (err error) {
	t, err := template.New("InventoryYAML").Parse(`
all:
  hosts:
    cloud-shell:
      ansible_user: {{ .UserName }}
      ansible_ssh_private_key_file: {{ .PrivateKey }}
      ansible_host: {{ .Host }}
      ansible_port: {{ .Port }}
`)
	if err != nil {
		log.Panicf("Could not create template: %v\n", err)
		return
	}
	f, err := os.Create("./ansible-playbook/generated-inventory.yaml")
	if err != nil {
		log.Panicf("Could not open file: %v\n", err)
		return
	}
	err = t.Execute(f, inventory)
	if err != nil {
		log.Panicf("Could not execute template: %v\n", err)
		return
	}
	return nil
}

package main

import (
	"flag"
	"log"

	"docker-builder/models"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path del file di configurazione se diverso dalla cartella di esecuzione")
	flag.Parse()
	config, err := models.InitConfig(configPath)
	defer config.Close()

	if err != nil {
		log.Fatalf("Config file errero: %s", err)
	}

	docker := NewDocker(config)

	log.Printf("Folders in: %s\r\n", config.Path)
	for k, v := range config.Images {
		log.Printf("Start build: %s\r\n", k)

		docker.BuildImage(k, v)

		if v.Push {
			log.Printf("Push: %s\r\n", k)
			docker.PushImage(k)
		}
	}
}

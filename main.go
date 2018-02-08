package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/brainlabs/mocking_server/route"
	config "github.com/spf13/viper"
)

func main() {

	config.AddConfigPath("./")
	config.SetConfigName("config")

	err := config.ReadInConfig()

	fmt.Println("config errr: ", err)

	fmt.Println(config.GetString("app.listen"))

	r := new(route.Router).Make()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.GetString("app.listen")), r))
}

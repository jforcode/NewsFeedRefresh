package main

import (
	"database/sql"
	"fmt"

	"github.com/magiconair/properties"
)

func main() {
	p := properties.MustLoadFile("app.properties", properties.UTF8)
	dataSource := p.GetString("datasource", "")
	apiKey := p.GetString("apiKey", "")
	apiUrl := p.GetString("apiUrl", "")
	requestLimit := p.GetInt("requestLimit", 0)

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	fmt.Println(dataSource, apiKey, apiUrl, requestLimit, db)
}

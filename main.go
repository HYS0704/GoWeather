package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func main() {
	r := gin.Default()
	r.SetHTMLTemplate(template.Must(template.ParseFiles("templates/index.html")))

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"weather": nil})
	})

	r.POST("/weather", func(c *gin.Context) {
		city := c.PostForm("city")
		apiKey := "3ff928fda6d0b93a72d1e27406ee808b" 
		url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric&lang=zh_tw", city, apiKey)

		resp, err := http.Get(url)
		if err != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{"error": "查詢失敗"})
			return
		}
		defer resp.Body.Close()

		var data WeatherResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{"error": "解析失敗"})
			return
		}

		if len(data.Weather) > 0 {
			result := fmt.Sprintf("%s 現在氣溫是 %.1f°C，天氣狀況：%s", data.Name, data.Main.Temp, data.Weather[0].Description)
			c.HTML(http.StatusOK, "index.html", gin.H{"weather": result})
		} else {
			c.HTML(http.StatusOK, "index.html", gin.H{"error": "找不到天氣資料"})
		}
	})

	r.Run(":8000")
}

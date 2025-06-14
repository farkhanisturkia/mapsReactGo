package main

import (
	"encoding/csv"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"io"
	"math"
	"os"
	"strconv"
)

type Point struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
}

type RouteRequest struct {
	Current Point   `json:"current"`
	Points  []Point `json:"points"`
}

func haversine(p1, p2 Point) float64 {
	const R = 6371

	lat1 := p1.Lat * math.Pi / 180
	lat2 := p2.Lat * math.Pi / 180
	dLat := (p2.Lat - p1.Lat) * math.Pi / 180
	dLng := (p2.Lng - p1.Lng) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func findNearestRoute(current Point, points []Point) []Point {
	visited := make(map[string]bool)
	route := []Point{current}
	visited[current.Name] = true
	remaining := len(points)

	for remaining > 0 {
		minDistance := math.MaxFloat64
		var nearest Point

		for _, p := range points {
			if !visited[p.Name] {
				distance := haversine(route[len(route)-1], p)
				if distance < minDistance {
					minDistance = distance
					nearest = p
				}
			}
		}

		if minDistance != math.MaxFloat64 {
			route = append(route, nearest)
			visited[nearest.Name] = true
			remaining--
		}
	}

	return route
}

func main() {
	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/api/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from Go Backend!",
		})
	})

	r.POST("/api/route", func(c *gin.Context) {
		var request RouteRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		route := findNearestRoute(request.Current, request.Points)

		c.JSON(200, gin.H{
			"route": route,
		})
	})

	r.POST("/api/upload-csv", func(c *gin.Context) {
		file, _, err := c.Request.FormFile("csv")
		if err != nil {
			c.JSON(400, gin.H{"error": "No CSV file uploaded"})
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		var points []Point
		header := true

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				c.JSON(400, gin.H{"error": "Invalid CSV format"})
				return
			}

			if header {
				header = false
				continue
			}

			if len(record) != 3 {
				c.JSON(400, gin.H{"error": "CSV must have 3 columns: name,lat,lng"})
				return
			}

			lat, err := strconv.ParseFloat(record[1], 64)
			if err != nil {
				c.JSON(400, gin.H{"error": "Invalid latitude"})
				return
			}

			lng, err := strconv.ParseFloat(record[2], 64)
			if err != nil {
				c.JSON(400, gin.H{"error": "Invalid longitude"})
				return
			}

			points = append(points, Point{
				Name: record[0],
				Lat:  lat,
				Lng:  lng,
			})
		}

		jsonData, err := json.MarshalIndent(points, "", "  ")
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to convert to JSON"})
			return
		}

		err = os.WriteFile("../frontend/public/data.json", jsonData, 0644)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to save JSON file"})
			return
		}

		c.JSON(200, gin.H{
			"message": "CSV processed and saved as data.json",
		})
	})

	r.Run(":9990")
}

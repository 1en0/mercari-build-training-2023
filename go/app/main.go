package main

import (
	"fmt"

	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const (
	ImgDir = "images"
)

type Response[T Item | Items | string] struct {
	Message T `json:"message"`
}

type Item struct {
	Name          string `json:"name"`
	Category      string `json:"category"`
	ImageFilename string `json:"image_filename"`
}

type Items struct {
	Items []Item `json:"items"`
}

func root(c echo.Context) error {
	res := Response[string]{Message: "Hello, world!"}
	return c.JSON(http.StatusOK, res)
}

func addItem(c echo.Context) error {
	// Get form data
	name := c.FormValue("name")
	c.Logger().Infof("Receive item: %s", name)
	message := fmt.Sprintf("item received: %s", name)
	category := c.FormValue("category")

	img, _ := c.FormFile("image")

	//imgName := getImageHash(img) + ".jpg"
	imgName := getImageHash(img) + path.Ext(img.Filename)

	err := saveImage(img, imgName)
	if err != nil {
		return err
	}
	//err = updateFile(name, category, imgName)
	//if err != nil {
	//	return err
	//}

	err = addItemInDb(name, category, imgName)
	if err != nil {
		return err
	}
	res := Response[string]{Message: message}

	return c.JSON(http.StatusOK, res)
}

func getItemsById(c echo.Context) error {
	itemId, _ := strconv.Atoi(c.Param("id"))
	//items, err := readItemListFromFile()
	//if err != nil {
	//	return err
	//}
	items, err := getItemsInDb()
	if err != nil {
		return err
	}
	if itemId >= len(items.Items) {
		res := Response[string]{Message: "Item index out of range"}
		return c.JSON(http.StatusBadRequest, res)
	}
	//buf, err := json.Marshal(items.Items[itemId])
	res := Response[Item]{Message: items.Items[itemId]}
	return c.JSON(http.StatusOK, res)
}

func searchItems(c echo.Context) error {
	keyword := c.QueryParam("keyword")
	items, err := getItemsByKeywordInDb(keyword)
	if err != nil {
		return err
	}
	res := Response[Items]{Message: *items}
	return c.JSON(http.StatusOK, res)
}

func getItems(c echo.Context) error {
	//items, err := readItemListFromFile()
	items, err := getItemsInDb()
	if err != nil {
		return err
	}
	//message, _ := strconv.Unquote(bufStr)
	//print(bufStr + "\n")

	res := Response[Items]{Message: *items}
	return c.JSON(http.StatusOK, res)
}

func getImg(c echo.Context) error {
	// Create image path
	imgPath := path.Join(ImgDir, c.Param("imageFilename"))

	if !strings.HasSuffix(imgPath, ".jpg") {
		res := Response[string]{Message: "Image path does not end with .jpg"}
		return c.JSON(http.StatusBadRequest, res)
	}
	if _, err := os.Stat(imgPath); err != nil {
		c.Logger().Debugf("Image not found: %s", imgPath)
		imgPath = path.Join(ImgDir, "default.jpg")
	}
	return c.File(imgPath)
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.SetLevel(log.INFO)

	front_url := os.Getenv("FRONT_URL")
	if front_url == "" {
		front_url = "http://localhost:3000"
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{front_url},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Init database
	DbInit()

	defer Db.Close()

	// Routes
	e.GET("/", root)
	e.POST("/items", addItem)
	e.GET("/items", getItems)
	e.GET("/image/:imageFilename", getImg)
	e.GET("/items/:id", getItemsById)
	e.GET("/search", searchItems)

	// Start server
	e.Logger.Fatal(e.Start(":9000"))

}

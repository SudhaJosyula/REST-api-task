package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo"
)

type Documents struct {
	ID           ObjectID  `json:"_id"`
	Name         string    `json:"name"`
	ParentID     ObjectID  `json:"parentId"`
	TenantID     ObjectID  `json:"tenantId"`
	Type         string    `json:"type"`
	IsDeleted    bool      `json:"isDeleted"`
	CreatedDate  Date      `json:"createdDate"`
	ModifiedDate Date      `json:"modifiedDate"`
	CreatedBy    User      `json:"createdBy"`
	ModifiedBy   User      `json:"modifiedBy"`
	PartitionKey string    `json:"partitionKey"`
	AttributeID  *ObjectID `json:"attributeId,omitempty"`
	CategoryID   *ObjectID `json:"categoryId,omitempty"`
	Value        json.RawMessage `json:"value,omitempty"`
	DataType     *string   `json:"dataType,omitempty"`
}


type ObjectID struct {
	Oid string `json:"$oid"`
}
type Date struct {
	Date string `json:"$date"`
}
type User struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Oid   string `json:"oid"`
}

func LoadData(filename string, dest interface{}) error {
	filePath := filepath.Join("mongo-data", filename)
	fmt.Println("Loading JSON from:", filePath) // Debug statement
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return nil
}

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello Prapancham")
	})

	e.GET("/documents", getDocs)
	e.GET("/documents/core_metadata/:id", getDocsByID)
	e.GET("/documents/files", getFiles)
	e.GET("/documents/folders", getFolders)

	e.Logger.Fatal(e.Start(":8080"))
}

func getDocs(c echo.Context) error {
	var documents []Documents
	if err := LoadData("rolodex.metadata.json", &documents); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, documents)
}

func getDocsByID(c echo.Context) error {
	id := c.Param("id")

	var documents []Documents
	if err := LoadData("rolodex.metadata.json", &documents); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}
	var core_metadata []Documents
	for _, doc := range documents {
		if doc.ParentID.Oid == id {
			// return c.JSON(http.StatusOK, doc)
			core_metadata = append(core_metadata, doc)
		}
	}
	return c.JSON(http.StatusOK, core_metadata)

	// return c.JSON(http.StatusNotFound, echo.Map{"message": "Document not found"})
}

func getFiles(c echo.Context) error {
	var documents []Documents
	if err := LoadData("rolodex.metadata.json", &documents); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	var files []Documents
	for _, doc := range documents {
		if doc.Type == "file" {
			files = append(files, doc)
		}
	}
	return c.JSON(http.StatusOK, files)
}

func getFolders(c echo.Context) error {
	var documents []Documents
	if err := LoadData("rolodex.metadata.json", &documents); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	var folders []Documents
	for _, doc := range documents {
		if doc.Type == "folder" {
			folders = append(folders, doc)
		}
	}
	return c.JSON(http.StatusOK, folders)
}

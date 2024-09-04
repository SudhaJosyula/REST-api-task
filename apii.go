package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

//Handler functions to routes

// get handlers

//to get a file 
func getFileHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)

	folderId := vars["folderId"]
	fileId := vars["fileId"]
	tenantId := vars["tenantId"]

	url := os.Getenv("FOLDER_URL") + folderId + "/files/" + fileId
	tok := os.Getenv("TOKEN")
	getreq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Error creating GET request", http.StatusInternalServerError)
		return
	}
	getreq.Header.Set("x-tenant-id", tenantId)
	getreq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tok))
	getreq.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make the request
	getresp, err := client.Do(getreq)
	if err != nil {
		http.Error(w, "Error making GET request", http.StatusInternalServerError)
		return
	}
	defer getresp.Body.Close()

	// Read the response body for GET request
	getBody, err := ioutil.ReadAll(getresp.Body)
	if err != nil {
		http.Error(w, "Error reading GET response body", http.StatusInternalServerError)
		return
	}

	// Send the GET response back to the client
	w.Header().Set("Content-Type", "application/json")
	w.Write(getBody)

}

//to get a folder 
func getFolderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantId := vars["tenantId"]

	url := os.Getenv("FOLDER_URL")

	token := os.Getenv("TOKEN")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Error creating GET request", http.StatusInternalServerError)
		return
	}
	// Set headers
	req.Header.Set("x-tenant-id", tenantId) // Replace with the appropriate value
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Create a new HTTP client
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make the request
	getresp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	defer getresp.Body.Close()

	// Read the response body for GET request
	getBody, err := ioutil.ReadAll(getresp.Body)
	if err != nil {
		http.Error(w, "Error reading GET response body", http.StatusInternalServerError)
		return
	}

	// Send the GET response back to the client
	w.Header().Set("Content-Type", "application/json")
	w.Write(getBody)

}

//to get metadata of a file or folder 
func getMetadata(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	objectId := vars["id"]
	tenanatId := vars["tenantId"]
	token := os.Getenv("TOKEN")
	url := os.Getenv("META_DATA_URL") + "/metadata?objectId="+objectId
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Error while creating the GET request", http.StatusInternalServerError)

	}

	req.Header.Set("objectId", objectId)
	req.Header.Set("x-tenant-id", tenanatId)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	getresp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error while making GET request", http.StatusInternalServerError)
	}
	defer getresp.Body.Close()
	respBody, err := ioutil.ReadAll(getresp.Body)
	if err != nil {
		http.Error(w, "Error reading GET response body", http.StatusInternalServerError)
		return
	}

	// Send the back to the client
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBody)

}

// Post Request Handlers

//to create a folder 
func createFolder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	root_folder_id := vars["id"]
	tenanatId := vars["tenantId"]
	requestBodyStruct := struct {
		Name string `json:"name"`
	}{
		Name: vars["name"],
	}

	body, err := json.Marshal(requestBodyStruct)
	if err != nil {
		http.Error(w, "Error marshalling request body", http.StatusInternalServerError)
		return
	}

	url := os.Getenv("FOLDER_URL") + "/" + root_folder_id
	token := os.Getenv("TOKEN")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "Error creating POST request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("id", root_folder_id)
	req.Header.Set("x-tenant-id", tenanatId)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{
		Timeout: 10 * time.Second, // Set the timeout duration
	}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error making POST request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Send the response back to the client
	postBody, err := ioutil.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	w.Write(postBody)

}

//add metadata to a file or folder 

func addMetadata(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantId := vars["tenantId"]
	objectId := vars["objectId"]
	attributeName := vars["attrName"]
	attributeValue := vars["attrValue"]
	url := os.Getenv("META_DATA_URL") + "/attributes"

	token := os.Getenv("TOKEN")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("x-tenant-id", tenantId)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error making GET request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	type CreatedBy struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		OID   string `json:"oid"`
	}

	type ModifiedBy struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		OID   string `json:"oid"`
	}

	type ResponseData struct {
		ID           string      `json:"id"`
		Name         string      `json:"name"`
		DataType     string      `json:"dataType"`
		Description  string      `json:"description"`
		Required     bool        `json:"required"`
		TenantID     string      `json:"tenantId"`
		DefaultValue interface{} `json:"defaultValue"`
		CreatedDate  string      `json:"createdDate"`
		ModifiedDate string      `json:"modifiedDate"`
		CreatedBy    CreatedBy   `json:"createdBy"`
		ModifiedBy   ModifiedBy  `json:"modifiedBy"`
		Type         string      `json:"type"`
	}

	type attributeData struct {
		Data []ResponseData `json:"data"`
	}

	type Metadata struct {
		AttributeID string `json:"attributeId"`
		Value       string `json:"value"`
	}

	type BodyStruct struct {
		ObjectID string     `json:"objectId"`
		Metadata []Metadata `json:"metadata"`
	}

	// Inside your handler
	var attributeDataResponse attributeData
	err = json.Unmarshal(respBody, &attributeDataResponse)
	if err != nil {
		http.Error(w, "Error unmarshalling GET response body", http.StatusInternalServerError)
		return
	}

	var body BodyStruct
	for _, attribute := range attributeDataResponse.Data {
		if attribute.Name == attributeName {
			body = BodyStruct{
				ObjectID: objectId,
				Metadata: []Metadata{
					{
						AttributeID: attribute.ID,
						Value:       attributeValue,
					},
				},
			}
			break
		}
	}

	requestBody, err := json.Marshal(body)
	if err != nil {
		http.Error(w, "Error marshalling request body", http.StatusInternalServerError)
		return
	}

	//  POST request
	URL := os.Getenv("META_DATA_URL") + "/metadata"
	request, err := http.NewRequest("POST", URL, bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("x-tenant-id", tenantId)
	request.Header.Set("Content-Type", "application/json")

	client2 := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp2, err := client2.Do(request)
	if err != nil {
		http.Error(w, "Error making POST request", http.StatusInternalServerError)
		return
	}
	defer resp2.Body.Close()

	// return the response from the POST request
	responseBody2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		http.Error(w, "Error reading POST response body", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody2)

}

func main() {

	r := mux.NewRouter()

	// Set the HTTP handlers for specific routes
	r.HandleFunc("/getfile/{folderId}/{fileId}/{tenantId}", getFileHandler).Methods("GET") //to download a file
	r.HandleFunc("/getfolder/{tenantId}", getFolderHandler)                                //to get details of folder
	r.HandleFunc("/createFolder/{tenantId}/{id}/{name}", createFolder)                     //to create a folder
	r.HandleFunc("/addMetadata/{tenantId}/{objectId}/{attrName}/{attrValue}", addMetadata) //to add metadata to a file or folder
	r.HandleFunc("/getMetadata/{tenantId}/{id}", getMetadata)                               // to get metadata of a file or folder
	// Start the server on port 8080
	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("Error starting server:", err)
	}

}

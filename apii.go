package main

import (
	
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"io"
	"net/http"
	"os"
	"time"
	"path"
	"mime"
	"regexp"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

//Handler functions to routes

// get handlers

//to download a file from a specific folder
func getFileHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	folderId := vars["folderId"]
	fileId := vars["fileId"]
	tenantId := vars["tenantId"]

	url := os.Getenv("FOLDER_URL")  + "/"+ folderId + "/files/" + fileId
	token := os.Getenv("TOKEN")
	request, err := http.NewRequest("GET" , url , nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	request.Header.Set("x-tenant-id", tenantId) // Replace with the appropriate value
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Accept",  "application/json"  )
	client := http.Client{Timeout: time.Second * 10}

	response, err := client.Do(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	filename, err := getFilenameFromResponse(response, url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	outFile, err := os.Create(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, response.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Write( []byte("File successfully downloaded"))


}
//extracts the filename from the HTTP response.
func getFilenameFromResponse(resp *http.Response, url string) (string, error) {
    contentDisposition := resp.Header.Get("Content-Disposition")
    if contentDisposition != "" {
        // Look for 'filename=' in the header
        _, params, err := mime.ParseMediaType(contentDisposition)
        if err == nil && params["filename"] != "" {
            return params["filename"], nil
        }
        regEx := regexp.MustCompile(`(?i)filename\*?=["']?([^;"']+)["']?`)
        matches := regEx.FindStringSubmatch(contentDisposition)
        if len(matches) > 1 {
            return matches[1], nil
        }
        return "", fmt.Errorf("unable to extract filename from Content-Disposition header")
    }
 
    
    filename := path.Base(url)
    if filename == "/" || filename == "" {
        return "", fmt.Errorf("unable to determine filename")
    }
    return filename, nil
}

//get the folder based on tenant Id
func getFolderMetadataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantId := vars["tenantId"]

	url := os.Getenv("FOLDER_URL")

	token := os.Getenv("TOKEN")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Error creating GET request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("x-tenant-id", tenantId) // Replace with the appropriate value
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	getresp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	defer getresp.Body.Close()

	getBody, err := ioutil.ReadAll(getresp.Body)
	if err != nil {
		http.Error(w, "Error reading GET response body", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(getBody)

}

//get detaiils of a specific folder 

func getFolderHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]
	tenanatId := vars["tenantId"]
	token := os.Getenv("TOKEN")
	url := os.Getenv("FOLDER_URL") + "/" + id

	request,err := http.NewRequest("GET" , url , nil)
	if err != nil{
		http.Error(w, err.Error() , http.StatusInternalServerError)
	}
	request.Header.Set("x-tenant-id", tenanatId)
	request.Header.Set("id", id)
	request.Header.Set("Authorization" , fmt.Sprintf("Bearer %s", token))

	client := &http.Client{
		Timeout: time.Second *10,
	}

	response , err  := client.Do(request)
	if err != nil{
		http.Error(w, err.Error() , http.StatusInternalServerError)
	}
	defer response.Body.Close()
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil{
		http.Error(w, err.Error() , http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBody)
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

func addMetadata(w http.ResponseWriter, r *http.Request) {   //fetches the sttributes 
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

func sharePermissions(w http.ResponseWriter,  r *http.Request){
	token := os.Getenv("TOKEN")
	vars := mux.Vars(r)
	tenantId := vars["tenantId"]
	objectId := vars["objectId"]
	objectType := vars["objectType"]
	entityId := vars["entityId"]
	entityType := vars["entityType"]
	relation := vars["relation"]

	type RelationStruct struct {
        ObjectId   string `json:"objectId"`
        ObjectType string `json:"objectType"`
        EntityId   string `json:"entityId"`
        EntityType string `json:"entityType"`
        Relation   string `json:"relation"`
    }
	type BodyStruct struct {
        Relations []RelationStruct `json:"relations"`
    }

    // Create the request body with the correct format
    reqBody := BodyStruct{
        Relations: []RelationStruct{
            {
                ObjectId:   objectId,
                ObjectType: objectType,
                EntityId:   entityId,
                EntityType: entityType,
                Relation:   relation,
            },
        },
    }

    // Marshal the request body to JSON
    body, err := json.Marshal(reqBody)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	url := os.Getenv("SHARE_URL")
	request , err := http.NewRequest( "POST" ,url , bytes.NewBuffer(body))

	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}
	request.Header.Set("Authorization" , fmt.Sprintf("Bearer %s", token))
	request.Header.Set("x-tenant-id", tenantId)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	respBody , err := client.Do(request)
	if err != nil {
		http.Error(w, "Error making POST request", http.StatusInternalServerError)
		return
	}
	defer respBody.Body.Close()

	// return the response from the POST request
	response, err := ioutil.ReadAll(respBody.Body)
	if err != nil {
		http.Error(w, "Error reading POST response body", http.StatusInternalServerError)
		return
	}

	w.Write(response )
	
	



}

func onBoardingHandler(w http.ResponseWriter,  r *http.Request){
	vars := mux.Vars(r)
	ad_group_id := vars["ad_group_id"]
	token := os.Getenv("TOKEN")
	// url := os.Getenv("ONBOARDING_URL") + "?/ad_group_id="+ ad_group_id
	url := "https://rolodex.dev.maersk-digital.net/api/v1/storage/tenant?ad_group_id=" + ad_group_id
	request,err := http.NewRequest( "GET",url , nil)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
	// request.Header.Set("Authentication" , fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("ad_group_id" , ad_group_id)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	respBody , err := client.Do(request)
	if err != nil{
		http.Error(w, err.Error() , http.StatusInternalServerError)
		return
	}
	defer respBody.Body.Close()
	response, err := ioutil.ReadAll(respBody.Body)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(response)

}

func deleteFile(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	fileId := vars["fileId"]
	folderId := vars["folderId"]
	tenantId := vars["tenantId"]
	token := os.Getenv("TOKEN")
	url := os.Getenv("FOLDER_URL") + "/" + folderId + "/files/" + fileId
	request,err := http.NewRequest("DELETE" , url , nil)

	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}
	request.Header.Set("Authorization" , fmt.Sprintf("Bearer %s", token))
	request.Header.Set("x-tenant-id", tenantId)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	respBody , err := client.Do(request)
	if err != nil {
		http.Error(w, "Error making POST request", http.StatusInternalServerError)
		return
	}
	defer respBody.Body.Close()

	// return the response from the POST request
	response, err := ioutil.ReadAll(respBody.Body)
	if err != nil {
		http.Error(w, "Error reading POST response body", http.StatusInternalServerError)
		return
	}
	fmt.Println(response)

	w.Write([]byte("deleted successfully"))
	
}


func main() {

	r := mux.NewRouter()

	// Set the HTTP handlers for specific routes
	r.HandleFunc("/getfile/{folderId}/{fileId}/{tenantId}", getFileHandler)//to download a file
	r.HandleFunc("/getfolder/{tenantId}", getFolderMetadataHandler) // to get a folder
	r.HandleFunc("/getfolder/{tenantId}/{id}" , getFolderHandler)                               //to get folder metadata by tenanatId
	
	r.HandleFunc("/createFolder/{tenantId}/{id}/{name}", createFolder)                     //to create a folder
	r.HandleFunc("/addMetadata/{tenantId}/{objectId}/{attrName}/{attrValue}", addMetadata) //to add metadata to a file or folder
	r.HandleFunc("/getMetadata/{tenantId}/{id}", getMetadata)                               // to get metadata of a file or folder
	r.HandleFunc("/sharePermission/{tenantId}/{objectId}/{objectType}/{entityId}/{entityType}/{relation}" , sharePermissions) 
	r.HandleFunc("/onBoarding/{ad_group_id}" , onBoardingHandler) //getting details of a tenant
	r.HandleFunc("/deleteFile/{fileId}/{folderId}/{tenantId}" , deleteFile) //delete file
	// Start the server on port 8080
	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("Error starting server:", err)
	}

}

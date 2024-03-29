package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-martini/martini"
)

type inBody struct {
	In string `json:"in"`
}

type Docker struct {
	Name    string `json:"Name"`
	Image   string `json:"Image"`
	Command string `json:"Command"`
}

type Tarball struct {
	Data      []byte `json:"Data"`
	Container Docker `json:"Container"`
}

func main() {
	var inputCollection []string

	uploads := make(map[string]Tarball)

	m := martini.Classic()
	m.Get("/", func() string {
		collection := "Collected inputs:\n"
		for i, input := range inputCollection {
			collection += fmt.Sprintf("index: %v \t|\tinput was: {%v}\n", i, input)
		}
		return collection
	})

	m.Post("/in", func(res http.ResponseWriter, req *http.Request) string {
		fmt.Printf("Received a message!")
		body, err := ioutil.ReadAll(req.Body)
		fmt.Printf("You gave me (raw): %s", string(body))
		if err != nil {
			fmt.Printf("something went wrong")
			return "something went wrong"
		}
		var in inBody
		err = json.Unmarshal(body, &in)
		if err != nil {
			fmt.Printf("something else went wrong")
			inputCollection = append(inputCollection, "Couldnt parse it, but here's what I got:\n"+string(body))
			return "something else went wrong"
		}
		inputCollection = append(inputCollection, in.In)
		fmt.Printf("You gave me: %s", in.In)
		return fmt.Sprintf("here's what you gave me: "+in.In+"\ncurrent collection size: %v", len(inputCollection))
	})
	m.Post("/upload_container", func(res http.ResponseWriter, req *http.Request) string {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return "something went wrong"
		}
		var tarball Tarball
		err = json.Unmarshal(body, &tarball)
		if err != nil {
			return fmt.Sprintf("data did not read as Container/Data pair %s", string(body))
		}
		//		path := "./uploads/"+tarball.Container.Name+".tgz"
		//		if _, err := os.Stat(path); !os.IsNotExist(err) {
		//			return "this container has already been uploaded"
		//		}
		//		err = writeFileToDisk(path, tarball.Data)
		//		if err != nil {
		//			return "error writing file to disk" + err.Error()
		//		}
		if _, haskey := uploads[tarball.Container.Name]; haskey {
			return fmt.Sprintf("error: %s has already been uploaded.", tarball.Container.Name)
		}
		uploads[tarball.Container.Name] = tarball
		return fmt.Sprintf("%s saved successfully", tarball.Container.Name)
	})

	m.Get("/download_container/:container_name", func(params martini.Params) ([]byte) {
		containerName := params["container_name"]
		if _, haskey := uploads[containerName]; !haskey {
			return []byte("500: error: " + containerName + " does not exist.")
		}
		tarball := uploads[containerName]
		response, err := json.Marshal(tarball)
		if err != nil {
			return []byte("500: error converting tarball into a response")
		}
		delete(uploads, containerName)
		return response
	})

	m.Get("/list_uploaded_containers", func() ([]byte) {
		uploaded_containers := "Containers I've got:"
		for key, _ := range uploads {
			uploaded_containers = fmt.Sprintf("%s\n%s", uploaded_containers, key)
		}
		return []byte(uploaded_containers)
	})

	m.Run()
}

//func writeFileToDisk(path string, data []byte) error {
//	err := ioutil.WriteFile(path, data, 0666)
//	return err
//}
//
//func readFileFromDisk(path string) ([]byte, error) {
//	data, err := ioutil.ReadFile(path)
//	return data, err
//}

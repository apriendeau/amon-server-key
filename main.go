package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type CreateReqBody struct {
	Name string `json:"name"`
}

type Server struct {
	Provider  string `json:"provider"`
	Name      string `json:"name"`
	Key       string `json:"key"`
	ID        string `json:"id"`
	LastCheck string `json:"last_check"`
}
type ListResp struct {
	Status  int      `json:"status"`
	Servers []Server `json:"servers"`
}

func main() {
	key := os.Getenv("AMON_KEY")
	name := os.Getenv("SERVER_NAME")
	host := os.Getenv("AMON_HOST")
	u := constructUrl(host, "/api/v1/servers/create/", key)
	reqBody, err := json.Marshal(&CreateReqBody{Name: name})
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post(u, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("Failed to Register with AMON")
	}

	listUrl := constructUrl(host, "/api/v1/servers/list/", key)
	resp, err = http.Get(listUrl)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	listResp := &ListResp{}
	err = json.Unmarshal(b, &listResp)
	if err != nil {
		log.Fatal(err)
	}
	for _, server := range listResp.Servers {
		if name == server.Name {
			fmt.Fprintln(os.Stdout, server.Key)
			return
		}
	}
}

func constructUrl(host, path, key string) string {
	return "http://" + host + path + "?api_key=" + key
}

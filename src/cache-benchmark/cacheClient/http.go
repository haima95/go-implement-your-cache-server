package cacheClient

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type httpClient struct {
	*http.Client
	server string
}

func (c *httpClient) get(key string) string {
	resp, err := c.Get(c.server + key)
	if err != nil {
		log.Println(key)
		panic(err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return ""
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (c *httpClient) set(key, value string) {
	req, err := http.NewRequest(http.MethodPut,
		c.server+key, strings.NewReader(value))
	if err != nil {
		log.Println(key)
		panic(err)
	}
	resp, err := c.Do(req)
	if err != nil {
		log.Println(key)
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
}

func (c *httpClient) Run(cmd *Cmd) {
	if cmd.Name == "get" {
		cmd.Value = c.get(cmd.Key)
		return
	}
	if cmd.Name == "set" {
		c.set(cmd.Key, cmd.Value)
		return
	}
	panic("unknown cmd name " + cmd.Name)
}

func (c *httpClient) PipelinedRun([]*Cmd) {
	panic("httpClient pipelined run not implement")
}

func newHTTPClient(server string) *httpClient {
	client := &http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: 1}}
	return &httpClient{client, "http://" + server + ":12345/cache/"}
}

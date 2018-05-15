package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestTempDir(t *testing.T) {
	temp, _ := ioutil.TempDir("", "appengine-aetest")
	t.Log(temp)
}

func TestHandler(t *testing.T) {
	//
	var out bytes.Buffer
	//log
	t.Logf("Start server")
	//start server
	c := exec.Command("goapp", "serve") //default port=8080, adminPort=8000
	c.Stdout = &out
	c.Stderr = &out
	c.Start()
	defer c.Process.Kill()
	//
	time.Sleep(10 * time.Second)
	//log
	t.Logf("Create request for http://localhost:8080/")
	//create request
	req, err := http.NewRequest("GET", "http://localhost:8080/", nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	//log
	t.Logf("GET http://localhost:8080/")
	//do GET
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	//log
	t.Logf("GET http://localhost:8080/")
	//delay to wait worker to run
	time.Sleep(10 * time.Second)
	//quit
	quitReq, err := http.NewRequest("GET", "http://localhost:8000/quit", nil)
	quitResp, err := client.Do(quitReq)
	if err != nil {
		fmt.Errorf("GET /quit handler error: %v", err)
	}
	defer quitResp.Body.Close()
	//log
	t.Logf("GET http://localhost:8000/quit")
	//
	//serve para testar Handler... não precisa de instância
	//http.HandlerFunc(Handler).ServeHTTP(resp, req)
	//read response
	b, _ := ioutil.ReadAll(resp.Body)
	resp_content := string(b)
	//checking
	if !strings.Contains(resp_content, "Handler Success") {
		t.Errorf("Handler not working")
	}
	//log server content
	logserver := out.String()
	if !strings.Contains(logserver, "Worker succeeded") {
		t.Errorf("Worker not working")
	}
	//log response
	t.Logf(logserver)
}

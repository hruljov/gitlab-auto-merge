package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var private_token = ""
var groups_id = 

var source_branch = flag.String("source-branch", "dev", "source branch")
var target_branch = flag.String("target-branch", "staging", "target branch")

func main() {
	flag.Parse()
	project, err := doSmth()
	if err != nil {
		println(err.Error())
	}
	for _, p := range project.Projects {
		//   openMR(p.ID, p.Name)
		fmt.Println(p.ID, p.Name)
	}
	var name = "test"
	openMR(19211206, name)
	//fmt.Printf(*source_branch, *target_branch)
}
func doSmth() (Project, error) {
	getid, err := http.Get("https://gitlab.com/api/v4/groups/"+ groups_id +"order_by=name&access_token=" + private_token)
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(getid.Body)

	if err != nil {
		panic(err.Error())
	}

	var Project Project
	json.Unmarshal(body, &Project)

	return Project, nil

}
func openMR(id int, name string) {
	url := "https://gitlab.com/api/v4/projects/" + strconv.Itoa(id) + "/merge_requests?source_branch=" + *source_branch + "&target_branch=" + *target_branch + "&title=auto_merge"
	fmt.Printf(url)

	method := "POST"
	client := &http.Client{}
	merge_requests, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	merge_requests.Header.Add("Private-Token", private_token)

	res, err := client.Do(merge_requests)
	defer res.Body.Close()
	merge_requests_res, err := ioutil.ReadAll(res.Body)

	var MergeRequests mergeRequests
	json.Unmarshal(merge_requests_res, &MergeRequests)
	if MergeRequests.ChangesCount == nil {
		var msg = "No changes"
		closeMR(id, name, MergeRequests.Iid, msg)
	} else {
		merge(id, name, MergeRequests.Iid)
	}

}

func merge(id int, name string, iid int) {
	url := "https://gitlab.com/api/v4/projects/" + strconv.Itoa(id) + "/merge_requests/" + strconv.Itoa(iid) + "/merge?merge_when_pipeline_succeeds=true"
	method := "PUT"
	client := &http.Client{}
	merge, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}

	merge.Header.Add("Private-Token", private_token)

	res, err := client.Do(merge)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if res.StatusCode != 200 {
		fmt.Println(name, string(body))
	} else {
		fmt.Println(name, "ok")
	}

}
func closeMR(id int, name string, iid int, msg string) {
	url := "https://gitlab.com/api/v4/projects/" + strconv.Itoa(id) + "/merge_requests/" + strconv.Itoa(iid)
	method := "DELETE"
	fmt.Println(url)
	client := &http.Client{}
	del_merge_requests, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
	}
	del_merge_requests.Header.Add("Private-Token", private_token)

	res, err := client.Do(del_merge_requests)
	defer res.Body.Close()
	fmt.Println(name, "close mr:", msg)
}

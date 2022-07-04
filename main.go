package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
	"path/filepath"
)

var kratosUrl = "http://kratos-admin.dev.ukama.com"

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Use 'kratos-id-cleanup some@example.com' or with wildcard 'kratos-id-cleanup *@example.com'")
		fmt.Println("Set KRATOS_URL env ver to override default: http://kratos-admin")
		return
	}

	ku := os.Getenv("KRATOS_URL")
	if ku != "" {
		kratosUrl = ku
	}

	email := os.Args[1]

	ls := getAllIdentities()

	for _, i := range ls {
		tr := i["traits"].(map[string]interface{})
		match, err := filepath.Match(email, tr["email"].(string))
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		if match {
			fmt.Println("Deleting: ", tr["email"], "  id:", i["id"])
			deleteIdentity(i["id"].(string))
		}

		//fmt.Println("ID:", i["id"])
	}

	println("Done")
}

func getAllIdentities() []map[string]interface{} {
	c := resty.New()
	res := []map[string]interface{}{}

	i := 0
	for {
		ls := []map[string]interface{}{}
		fmt.Println("Fetching identities page:", i)
		resp, err := c.R().SetResult(&ls).Get(fmt.Sprintf("%s/identities?per_page=500&page=%d", kratosUrl, i))
		if err != nil {
			panic(err)
		}

		if resp.StatusCode() != 200 {
			panic(resp.Status())
		}

		res = append(res, ls...)
		if len(ls) == 0 {
			break
		}
		i++
	}

	return res
}

func deleteIdentity(id string) {
	c := resty.New()
	resp, err := c.R().Delete(kratosUrl + "/identities/" + id)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode() != 200 && resp.StatusCode() != 204 {
		if resp.StatusCode() < 500 {
			fmt.Println("Identity nod deletes. Response code: ", resp.StatusCode())
		} else {
			panic(resp.Status())
		}
	}

}

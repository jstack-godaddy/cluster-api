package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type SnowConnection struct {
	server   string
	user     string
	password string
}

type SnowGenericGet struct {
	Result []map[string]string `json:"Result"`
}

type SnowUserGroupsGet struct {
	Result []SnowUserGroupsGetResult `json:"Result"`
}

type SnowUserGroupsGetResult struct {
	Group          SnowUserGroupsGetResultUG `json:"Group"`
	User           SnowUserGroupsGetResultUG `json:"User"`
	U_temporary    string                    `json:"u_temporary"`
	Sys_id         string                    `json:"sys_id"`
	Sys_updated_by string                    `json:"sys_updated_by"`
	Sys_created_on string                    `json:"sys_created_on"`
	Sys_mod_count  string                    `json:"sys_mod_count"`
	Sys_updated_on string                    `json:"sys_updated_on"`
	Sys_tags       string                    `json:"sys_tags"`
	U_expiration   string                    `json:"u_expiration"`
}

type SnowUserGroupsGetResultUG struct {
	Link  string `json:"link"`
	Value string `json:"value"`
}

type SnowGroupGet struct {
	Result []map[string]string `json:"Result"`
}

func (s SnowConnection) queryTable(table string, field string, filter string) (bodyText []byte, err error) {
	client := &http.Client{}
	fullURL := s.server + table + "?" + field + "=" + filter

	req, err := http.NewRequest("GET", fullURL, nil)

	if err != nil {
		return
	}
	req.SetBasicAuth(s.user, s.password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	//fmt.Println(req)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	bodyText, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return
}

func GetTeams(u string) (teams []string, httpStatus int, err error) {
	godotenv.Load("./.snow_creds.env")
	httpStatus = http.StatusOK

	snow := SnowConnection{
		server:   os.Getenv("SNOW_SERVER"),
		user:     os.Getenv("SNOW_USER"),
		password: os.Getenv("SNOW_PASSWORD"),
	}

	userInfo, err := snow.queryTable("sys_user", "user_name", u)
	if err != nil {
		return
	}

	var uqr SnowGenericGet
	json.Unmarshal([]byte(userInfo), &uqr)

	if len(uqr.Result) == 0 {
		httpStatus = http.StatusNotFound
		teams = append(teams, "User not found.")
		return
	}
	sys_id := uqr.Result[0]["sys_id"]

	userGroupInfo, err := snow.queryTable("sys_user_grmember", "user", sys_id)
	if err != nil {
		httpStatus = http.StatusBadRequest
		return
	}

	var ugqr SnowUserGroupsGet
	err = json.Unmarshal([]byte(userGroupInfo), &ugqr)
	if err != nil {
		httpStatus = http.StatusBadRequest
		fmt.Println("Unmarshall err", err)
	}

	for _, result := range ugqr.Result {
		groupInfo, err := snow.queryTable("sys_user_group", "sys_id", result.Group.Value)
		if err != nil {
			httpStatus = http.StatusBadRequest
			fmt.Println("Query Error ", err)
		}
		//fmt.Println(groupInfo)
		var gqr SnowGroupGet
		_ = json.Unmarshal([]byte(groupInfo), &gqr)
		//if err != nil {
		//	fmt.Println("Unmarshall err ", err)
		//}
		//fmt.Println(gqr.Result[0]["name"])
		teams = append(teams, gqr.Result[0]["name"])
	}

	return
}

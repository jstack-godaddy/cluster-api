package helpers

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type ClusterDB struct {
	*sql.DB
}

// Structs to match DB Schema
type Datacenter struct {
	Id   int    `json:"dc_id"`
	Name string `json:"dc_name"`
	Abbr string `json:"dc_abbr"`
}

type Datastore struct {
	Id   int    `json:"datastore_id"`
	Name string `json:"datastore_name"`
	Abbr string `json:"datastore_abbr"`
}

type Environment struct {
	Id            int    `json:"env_id"`
	Name          string `json:"env_name"`
	Abbr          string `json:"env_abbr"`
	Single_letter string `json:"env_single_letter"`
}

type NetworkZone struct {
	Id   int    `json:"nz_id"`
	Name string `json:"nz_name"`
	Abbr string `json:"nz_abbr"`
}

type Team struct {
	Id            int    `json:"team_id"`
	Slack_channel string `json:"slack_channel"`
	Name          string `json:"team_friendly"`
	Snow_group    string `json:"team_snow"`
}

type Project struct {
	Id          int    `json:"project_id"`
	Name        string `json:"project_name"`
	Datacenter  string `json:"dc_id"`
	Team        string `json:"team_id"`
	Environment string `json:"env_id"`
}

type Cluster struct {
	Id           int    `json:"cluster_id"`
	Name         string `json:"cluster_name"`
	Project_id   int    `json:"project_id"`
	Network_zone int    `json:"nz_id"`
	Datastore    int    `json:"datastore_id"`
	Floater      string `json:"floater"`
	Floater_v6   string `json:"floater_v6"`
	Public       bool   `json:"public"`
	Created_on   string `json:"created_on"`
}

type Server struct {
	Id         int    `json:"server_id"`
	Cluster    string `json:"cluster_id"`
	Project_id int    `json:"project_id"`
	Hostname   string `json:"hostname"`
	IP         string `json:"ip"`
	IP_v6      string `json:"ipv6"`
	Created_on string `json:"created_on"`
}

func NewDBClient() (clusterDB *ClusterDB, err error) {
	godotenv.Load("./.db_creds.env")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/cluster_api", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"))
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	return &ClusterDB{db}, nil
}

func (db *ClusterDB) GetProjectsByTeam(team string) (projects []Project, httpStatus int, err error) {
	httpStatus = 200

	team_query := fmt.Sprintf("SELECT team_id FROM teams WHERE team_snow='%s';", team)
	var team_id int
	row := db.QueryRow(team_query)
	switch err := row.Scan(&team_id); err {
	case sql.ErrNoRows:
		httpStatus = 404
	case nil:
		httpStatus = 200
	default:
		panic(err)
	}

	projects_query := fmt.Sprintf("SELECT * FROM projects WHERE team_id=%d;", team_id)
	fmt.Println(projects_query)
	rows, err := db.Query(projects_query)
	fmt.Println(rows)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	for rows.Next() {
		var p Project
		err = rows.Scan(&p.Id, &p.Name, &p.Datacenter, &p.Team, &p.Environment)
		if err != nil {
			// handle this error
			panic(err)
		}
		fmt.Println(p)
		projects = append(projects, p)
	}

	return
}

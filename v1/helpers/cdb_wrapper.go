package helpers

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

// Struct for database object
type ClusterDB struct {
	*sqlx.DB
}

// Structs to match DB Schema
type Datacenter struct {
	Id   int    `db:"dc_id"`
	Name string `db:"dc_name"`
	Abbr string `db:"dc_abbr"`
}

type Datastore struct {
	Id   int    `db:"datastore_id"`
	Name string `db:"datastore_name"`
	Abbr string `db:"datastore_abbr"`
}

type Environment struct {
	Id            int    `db:"env_id"`
	Name          string `db:"env_name"`
	Abbr          string `db:"env_abbr"`
	Single_letter string `db:"env_single_letter"`
}

type NetworkZone struct {
	Id   int    `db:"nz_id"`
	Name string `db:"nz_name"`
	Abbr string `db:"nz_abbr"`
}

type Flavor struct {
	Id   int `db:"flavor_id"`
	Cpus int `db:"cpus"`
	Ram  int `db:"ram"`
	Disk int `db:"disk"`
	Cost int `db:"cost"`
}

type Team struct {
	Id            int    `db:"team_id"`
	Slack_channel string `db:"slack_channel"`
	Name          string `db:"team_friendly"`
	Snow_group    string `db:"team_snow"`
}

type Project struct {
	Id          int    `db:"project_id"`
	Name        string `db:"project_name"`
	Datacenter  string `db:"dc_id"`
	Team        string `db:"team_id"`
	Environment string `db:"env_id"`
}

type Cluster struct {
	Id           int    `db:"cluster_id"`
	Name         string `db:"cluster_name"`
	Project_id   int    `db:"project_id"`
	Network_zone int    `db:"nz_id"`
	Datastore    int    `db:"datastore_id"`
	Floater      string `db:"floater"`
	Floater_v6   string `db:"floater_v6"`
	Public       bool   `db:"public"`
	Created_on   string `db:"created_on"`
}

type Server struct {
	Id         int    `db:"server_id"`
	Cluster    string `db:"cluster_id"`
	Project_id int    `db:"project_id"`
	Hostname   string `db:"hostname"`
	IP         string `db:"ip"`
	IP_v6      string `db:"ipv6"`
	Created_on string `db:"created_on"`
}

const (
	dsn_const                  = "%s:%s@tcp(%s:3306)/cluster_api"
	all_datacenters_query      = `SELECT * FROM datacenters`
	all_datastores_query       = `SELECT * FROM datastores`
	all_environments_query     = `SELECT * FROM environments`
	all_network_zones_query    = `SELECT * FROM network_zones`
	all_flavors_query          = `SELECT * FROM flavors`
	team_id_from_name_query    = `SELECT team_id FROM teams WHERE team_snow=?`
	projects_by_team_query     = `SELECT * FROM projects WHERE team_id=?`
	project_id_from_name_query = `SELECT project_id FROM projects WHERE project_name=?`
	servers_by_project_query   = `SELECT * FROM servers WHERE project_id=?`
)

func NewClusterDBConn() (clusterDB ClusterDB, err error) {
	godotenv.Load("./.db_creds.env")
	dsn := fmt.Sprintf(dsn_const, os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"))
	db, err := sqlx.Open("mysql", dsn)
	clusterDB = ClusterDB{db}

	return
}

func (db *ClusterDB) GetAllDataCenters() (dcs []Datacenter, err error) {

	db.Select(&dcs, all_datacenters_query)

	return
}

func (db *ClusterDB) GetAllEnvironments() (envs []Environment, err error) {

	db.Select(&envs, all_environments_query)

	return
}

func (db *ClusterDB) GetAllNetworkZones() (nzs []NetworkZone, err error) {

	db.Select(&nzs, all_network_zones_query)

	return
}

func (db *ClusterDB) GetAllFlavors() (flavors []Flavor, err error) {

	db.Select(&flavors, all_flavors_query)

	return
}

func (db *ClusterDB) GetAllDatastores() (datastores []Datastore, err error) {

	db.Select(&datastores, all_datastores_query)

	return
}

func (db *ClusterDB) GetProjectsByTeam(team string) (projects []Project, err error) {

	var team_id int
	row := db.QueryRowx(team_id_from_name_query, team)
	_ = row.Scan(&team_id)

	db.Select(&projects, projects_by_team_query, team_id)

	return
}

func (db *ClusterDB) GetServersByProject(project_id string) (servers []Server, err error) {

	db.Select(&servers, servers_by_project_query, project_id)

	return
}
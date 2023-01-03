package helpers

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
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
	Environment  int    `db:"environment_id"`
	Datastore    int    `db:"datastore_id"`
	Floater      string `db:"floater"`
	Floater_v6   string `db:"floater_v6"`
	Created_on   string `db:"created_on"`
}
type Server struct {
	Id         int    `db:"server_id"`
	Cluster    int    `db:"cluster_id"`
	Flavor     int    `db:"flavor_id"`
	Hostname   string `db:"hostname"`
	OS         string `db:"operating_system"`
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
	server_by_shortname        = `SELECT * FROM servers WHERE hostname=?`
	select_project_by_name     = `SELECT * FROM projects WHERE project_name=?`
	select_cluster_by_name     = `SELECT * FROM clusters where cluster_name=?`
	select_nz_by_name          = `SELECT * FROM network_zones WHERE nz_name=?`
	select_datastore_by_name   = `SELECT * FROM datastores WHERE datastore_name=?`
	select_flavor_by_name      = `SELECT * FROM flavor WHERE flavor_name=?`
	select_env_by_sl           = `SELECT * FROM environments WHERE environment_single_letter=?`
	insert_cluster             = `INSERT INTO cluster (cluster_name,project_id,nz_id,datastore_id,floater,floater_v6) floater (:cluster_name,
		:project_id,:nz_id,:datastore_id,:floater,:floater_v6)`
	insert_server = `INSERT INTO server (cluster_id,hostname,nz_id,operating_system,ip,ipv6) floater (:cluster_id,:hostname,:operating_system,
		:datastore_id,:floater,:floater_v6)`
)

// Create new connection to OPSDB
func NewClusterDBConn() (cdb *ClusterDB) {
	godotenv.Load("./.db_creds.env")
	dsn := fmt.Sprintf(dsn_const, os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"))
	db, _ := sqlx.Open("mysql", dsn)

	return &ClusterDB{db}
}

// Some SELECTS to get information if needed.
func (db *ClusterDB) GetAllDataCenters() (dcs []Datacenter, err error) {

	err = db.Select(&dcs, all_datacenters_query)

	return
}

func (db *ClusterDB) GetAllEnvironments() (envs []Environment, err error) {

	err = db.Select(&envs, all_environments_query)

	return
}

func (db *ClusterDB) GetAllNetworkZones() (nzs []NetworkZone, err error) {

	err = db.Select(&nzs, all_network_zones_query)

	return
}

func (db *ClusterDB) GetAllFlavors() (flavors []Flavor, err error) {

	err = db.Select(&flavors, all_flavors_query)

	return
}

func (db *ClusterDB) GetAllDatastores() (datastores []Datastore, err error) {

	err = db.Select(&datastores, all_datastores_query)

	return
}

func (db *ClusterDB) GetProjectsByTeam(team string) (projects []Project, err error) {

	var team_id int
	row := db.QueryRowx(team_id_from_name_query, team)
	_ = row.Scan(&team_id)

	err = db.Select(&projects, projects_by_team_query, team_id)

	return
}

func (db *ClusterDB) GetServersByProject(project_id string) (servers []Server, err error) {

	err = db.Select(&servers, servers_by_project_query, project_id)

	return
}

// METADATA METHODS
// These methods insert or delete metadata from the OPSDB
func (db *ClusterDB) InsertMetadata(clusterName string, projectName string, os string, networkZone string, dbName string,
	envSL string, flavorName string, floatingIP *floatingips.FloatingIP, serversCreated []*servers.Server) (clusterResult Cluster, serversResult []Server, err error) {

	fmt.Println("attempting metadata insertion")
	// query for the needful
	var project Project
	err = db.Select(&project, select_project_by_name, projectName)
	if err != nil {
		return
	}

	var nz NetworkZone
	err = db.Select(&nz, select_nz_by_name, networkZone)
	if err != nil {
		return
	}

	var ds Datastore
	err = db.Select(&ds, select_datastore_by_name, dbName)
	if err != nil {
		return
	}

	var flavor Flavor
	err = db.Select(&flavor, select_flavor_by_name, flavorName)
	if err != nil {
		return
	}

	var env Environment
	err = db.Select(&env, select_env_by_sl, envSL)
	if err != nil {
		return
	}
	fmt.Println("got all IDs")

	// process raw input above into structs that match schema
	clusterRaw := Cluster{
		Name:         clusterName,
		Project_id:   project.Id,
		Network_zone: nz.Id,
		Environment:  env.Id,
		Datastore:    ds.Id,
		Floater:      floatingIP.IP,
	}

	clusterResult, err = db.insertClusterMetadata(clusterRaw)
	if err != nil {
		return
	}

	var serversRaw []Server
	for i := 1; i <= len(serversCreated); i++ {
		//msg += fmt.Sprintf("%d: %s\n", i, serversCreated[i-1].Name)
		addMe := Server{
			Cluster:  clusterResult.Id,
			Hostname: serversCreated[i-1].Name,
			Flavor:   flavor.Id,
			OS:       os,
			IP:       serversCreated[i-1].AccessIPv4,
			IP_v6:    serversCreated[i-1].AccessIPv6,
		}
		serversRaw = append(serversRaw, addMe)
	}
	err = db.insertServerMetadata(serversRaw)

	return
}

func (db *ClusterDB) insertClusterMetadata(cluster Cluster) (clusterResult Cluster, err error) {

	_, err = db.NamedExec(insert_cluster, cluster)
	if err != nil {
		return
	}

	err = db.Select(&clusterResult, select_cluster_by_name, cluster.Name)
	return
}

func (db *ClusterDB) insertServerMetadata(servers []Server) (err error) {

	_, err = db.NamedExec(insert_server, servers)

	return
}

package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"slashbase.com/backend/daos"
	"slashbase.com/backend/middlewares"
	"slashbase.com/backend/models"
	"slashbase.com/backend/utils"
	"slashbase.com/backend/views"
)

type DBConnectionController struct{}

var dbConnDao daos.DBConnectionDao

func (dbcc DBConnectionController) CreateDBConnection(c *gin.Context) {
	var createBody struct {
		ProjectID   string `json:"projectId"`
		Name        string `json:"name"`
		Host        string `json:"host"`
		Port        string `json:"port"`
		Password    string `json:"password"`
		User        string `json:"user"`
		DBName      string `json:"dbname"`
		UseSSH      string `json:"useSSH"`
		SSHHost     string `json:"sshHost"`
		SSHUser     string `json:"sshUser"`
		SSHPassword string `json:"sshPassword"`
		SSHKeyFile  string `json:"sshKeyFile"`
	}
	c.BindJSON(&createBody)
	authUser := middlewares.GetAuthUser(c)
	// TODO: check for admin permssion in project before creating new db
	dbConn, err := models.NewPostgresDBConnection(authUser.ID, createBody.ProjectID, createBody.Name, createBody.Host, createBody.Port,
		createBody.User, createBody.Password, createBody.DBName, createBody.UseSSH, createBody.SSHHost, createBody.SSHUser, createBody.SSHPassword, createBody.SSHKeyFile)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	err = dbConnDao.CreateDBConnection(dbConn)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    views.BuildDBConnection(dbConn),
	})
}

func (dbcc DBConnectionController) GetDBConnections(c *gin.Context) {
	authUserProjectIds := middlewares.GetAuthUserProjectIds(c)

	dbConns, err := dbConnDao.GetDBConnectionsByProjectIds(*authUserProjectIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	dbConnViews := []views.DBConnectionView{}
	for _, dbConn := range dbConns {
		dbConnViews = append(dbConnViews, views.BuildDBConnection(dbConn))
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dbConnViews,
	})
}

func (dbcc DBConnectionController) GetSingleDBConnection(c *gin.Context) {
	dbConnID := c.Param("dbConnId")
	_ = middlewares.GetAuthUser(c)
	dbConn, err := dbConnDao.GetDBConnectionByID(dbConnID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	// TODO: check if authUser is member of project
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    views.BuildDBConnection(dbConn),
	})
}

func (dbcc DBConnectionController) GetDBConnectionsByProject(c *gin.Context) {
	projectID := c.Param("projectId")
	authUserProjectIds := middlewares.GetAuthUserProjectIds(c)
	if !utils.ContainsString(*authUserProjectIds, projectID) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   errors.New("not allowed"),
		})
		return
	}

	dbConns, err := dbConnDao.GetDBConnectionsByProject(projectID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	dbConnViews := []views.DBConnectionView{}
	for _, dbConn := range dbConns {
		dbConnViews = append(dbConnViews, views.BuildDBConnection(dbConn))
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dbConnViews,
	})
}

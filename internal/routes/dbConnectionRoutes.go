package routes

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"slashbase.com/backend/internal/controllers"
	"slashbase.com/backend/internal/middlewares"
	"slashbase.com/backend/internal/models"
	"slashbase.com/backend/internal/utils"
	"slashbase.com/backend/internal/views"
)

type DBConnectionRoutes struct{}

var dbConnController controllers.DBConnectionController

func (dbcr DBConnectionRoutes) CreateDBConnection(c *gin.Context) {
	var createBody struct {
		ProjectID   string `json:"projectId"`
		Name        string `json:"name"`
		Type        string `json:"type"`
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

	if isAllowed, err := controllers.GetAuthUserHasRolesForProject(authUser, createBody.ProjectID, []string{models.ROLE_ADMIN}); err != nil || !isAllowed {
		return
	}

	dbConn, err := dbConnController.CreateDBConnection(authUser, createBody.ProjectID, createBody.Name, createBody.Type, createBody.Host, createBody.Port,
		createBody.User, createBody.Password, createBody.DBName, createBody.UseSSH, createBody.SSHHost, createBody.SSHUser, createBody.SSHPassword, createBody.SSHKeyFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
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

func (dbcr DBConnectionRoutes) GetDBConnections(c *gin.Context) {
	authUserProjectIds := middlewares.GetAuthUserProjectIds(c)

	dbConns, err := dbConnController.GetDBConnections(authUserProjectIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
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

func (dbcr DBConnectionRoutes) DeleteDBConnection(c *gin.Context) {
	dbConnID := c.Param("dbConnId")
	authUser := middlewares.GetAuthUser(c)
	err := dbConnController.DeleteDBConnection(authUser, dbConnID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func (dbcr DBConnectionRoutes) GetSingleDBConnection(c *gin.Context) {
	dbConnID := c.Param("dbConnId")
	authUser := middlewares.GetAuthUser(c)
	dbConn, err := dbConnController.GetSingleDBConnection(authUser, dbConnID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
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

func (dbcr DBConnectionRoutes) GetDBConnectionsByProject(c *gin.Context) {
	projectID := c.Param("projectId")
	authUserProjectIds := middlewares.GetAuthUserProjectIds(c)
	if !utils.ContainsString(*authUserProjectIds, projectID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   errors.New("not allowed"),
		})
		return
	}

	dbConns, err := dbConnController.GetDBConnectionsByProject(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
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

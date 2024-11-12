package handlers

import (
	graphApi "github.com/g3ksa/lab5_otrpo/internal/graph_api/service"
	"github.com/g3ksa/lab5_otrpo/internal/graph_api/service/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var secretToken string

type HTTPServer struct {
	Service *graphApi.Service
}

func NewHttpServer(storage graphApi.Storage, authToken string) *HTTPServer {
	secretToken = authToken
	return &HTTPServer{
		Service: graphApi.NewService(storage),
	}
}

func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != secretToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}

func (h *HTTPServer) Router() http.Handler {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/", h.getAllNodes)
		api.GET("/relations/:id", h.getNodeWithRelations)
		api.POST("/", AuthMiddleware, h.Insert)
		api.DELETE("/:id", AuthMiddleware, h.Delete)
	}

	return router
}

func (h *HTTPServer) getAllNodes(c *gin.Context) {
	nodes, err := h.Service.GetAllNodes(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, nodes)
}

func (h *HTTPServer) getNodeWithRelations(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	relations, err := h.Service.GetNodeWithRelations(c, id)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, relations)
}

func (h *HTTPServer) Insert(c *gin.Context) {
	var request model.InsertRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err := h.Service.Insert(c, request)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, map[string]bool{"success": true})
}

func (h *HTTPServer) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.Service.Delete(c, id)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, map[string]bool{"success": true})
}

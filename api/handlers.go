package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *API) getUser(c *gin.Context) {
	name := c.Param("name")
	for i, user := range a.data {
		if user.Name == name {
			a.data[i].AddLog(fmt.Sprintf("user requsted from IP: %s", c.ClientIP()))
			c.JSON(http.StatusOK, user)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

func (a *API) addPet(c *gin.Context) {
	name := c.Param("name")
	var pet Pet
	err := c.BindJSON(&pet)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}
	for i, user := range a.data {
		if user.Name == name {
			a.data[i].Pets = append(a.data[i].Pets, pet)
			a.data[i].AddLog(fmt.Sprintf("pet '%s' (%s) added from IP: %s", pet.Name, pet.Type, c.ClientIP()))
			c.JSON(http.StatusOK, a.data[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

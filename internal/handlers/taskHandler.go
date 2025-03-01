package handlers

import "github.com/gin-gonic/gin"

type Taskhandler struct{}

func DefaultTaskHandler() *Taskhandler {
	return &Taskhandler{}
}

func (h *Taskhandler) NewTaskHandler(c *gin.Context) {
}

func (h *Taskhandler) UpdateTaskHandler(c *gin.Context) {}

func (h *Taskhandler) DeleteTaskHandler(c *gin.Context) {}

func (h *Taskhandler) FinishTaskHandler(c *gin.Context) {}

func (h *Taskhandler) TaskListsHandler(c *gin.Context) {}

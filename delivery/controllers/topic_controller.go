package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"remedymate-backend/domain/AppError"
	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/interfaces"

	"github.com/gin-gonic/gin"
)

const defaultControllerTimeout = 5 * time.Second

type TopicController struct {
	topicUsecase interfaces.TopicUsecase
}

func NewTopicController(topicUsecase interfaces.TopicUsecase) *TopicController {
	return &TopicController{topicUsecase: topicUsecase}
}

// CreateTopicHandler validates JSON and creates a topic (temporary hard-coded user)
func (tc *TopicController) CreateTopicHandler(c *gin.Context) {
	var req dto.TopicCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c, defaultControllerTimeout)
	defer cancel()

	topic, err := tc.topicUsecase.CreateTopic(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, AppError.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, AppError.ErrTopicAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, topic)
}

// UpdateTopicHandler updates an existing topic (uses test user id for updated_by)
func (tc *TopicController) UpdateTopicHandler(c *gin.Context) {
	topicKey := c.Param("topic_key")
	var req dto.TopicUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c, defaultControllerTimeout)
	defer cancel()

	updatedTopic, err := tc.topicUsecase.UpdateTopic(ctx, topicKey, req)
	if err != nil {
		switch {
		case errors.Is(err, AppError.ErrTopicNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, AppError.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, updatedTopic)
}

// GetTopicHandler retrieves a topic by key
func (tc *TopicController) GetTopicHandler(c *gin.Context) {
	topicKey := c.Param("topic_key")
	if topicKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "topicKey is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c, defaultControllerTimeout)
	defer cancel()

	topic, err := tc.topicUsecase.GetTopicByKey(ctx, topicKey)
	if err != nil {
		if errors.Is(err, AppError.ErrTopicNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "topic not found"})
		} else {
			log.Printf("GetTopicHandler failed err=%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	c.JSON(http.StatusOK, topic)
}

// DeleteTopicHandler soft-deletes a topic (uses test user id when no user in context)
func (tc *TopicController) DeleteTopicHandler(c *gin.Context) {
	topicKey := c.Param("topic_key")
	ctx, cancel := context.WithTimeout(c, defaultControllerTimeout)
	defer cancel()

	if err := tc.topicUsecase.SoftDeleteTopic(ctx, topicKey); err != nil {
		switch {
		case errors.Is(err, AppError.ErrTopicNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

// ListAllTopicsHandler returns paginated topics
func (tc *TopicController) ListAllTopicsHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.DefaultQuery("search", "")
	sortBy := c.DefaultQuery("sort_by", "name_en")
	sortOrder := c.DefaultQuery("sort_order", "asc")

	params := dto.TopicListQueryParams{
		PaginationQueryParams: dto.PaginationQueryParams{Page: page, Limit: limit},
		FilterQueryParams:     dto.FilterQueryParams{Search: search},
		SortQueryParams:       dto.SortQueryParams{SortBy: sortBy, Order: sortOrder},
	}

	ctx, cancel := context.WithTimeout(c, defaultControllerTimeout)
	defer cancel()

	log.Printf("ListAllTopics called page=%d limit=%d", params.Page, params.Limit)

	res, err := tc.topicUsecase.ListAllTopics(ctx, params)
	if err != nil {
		log.Printf("ListAllTopics failed err=%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, res)
}

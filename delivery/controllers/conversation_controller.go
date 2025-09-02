package controllers

import (
	"net/http"

	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"

	"github.com/gin-gonic/gin"
)

type ConversationController struct {
	conversationUsecase interfaces.ConversationUsecase
}

func NewConversationController(conversationUsecase interfaces.ConversationUsecase) *ConversationController {
	return &ConversationController{
		conversationUsecase: conversationUsecase,
	}
}

// ConversationRequest represents a unified request for both starting and continuing conversations
type ConversationRequest struct {
	ConversationID string `json:"conversation_id,omitempty"` // Required for continuing, optional for starting
	Symptom        string `json:"symptom,omitempty"`         // Required for starting, optional for continuing
	Language       string `json:"language,omitempty"`        // Required for starting, optional for continuing
	Answer         string `json:"answer,omitempty"`          // Required for continuing, optional for starting
	UserID         string `json:"user_id,omitempty"`         // Optional for both
}

// ConversationResponse represents a unified response for both starting and continuing conversations
type ConversationResponse struct {
	ConversationID    string                 `json:"conversation_id"`
	Question          *entities.Question     `json:"question,omitempty"` // Next question if available
	Message           string                 `json:"message,omitempty"`  // Feedback message
	IsComplete        bool                   `json:"is_complete"`        // Whether all questions are answered
	CurrentStep       int                    `json:"current_step"`
	TotalSteps        int                    `json:"total_steps"`
	Report            *entities.HealthReport `json:"report,omitempty"`    // Final report if complete
	IsNewConversation bool                   `json:"is_new_conversation"` // Whether this is a new conversation
}

// HandleConversation handles both starting and continuing conversations in a single endpoint
// POST /api/v1/conversation
func (cc *ConversationController) HandleConversation(c *gin.Context) {
	var req ConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Determine if this is a new conversation or continuing an existing one
	if req.ConversationID == "" {
		// Starting a new conversation
		if req.Symptom == "" || req.Language == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Missing required fields",
				Details: "symptom and language are required for starting a new conversation",
			})
			return
		}

		startReq := dto.StartConversationRequest{
			Symptom:  req.Symptom,
			Language: req.Language,
			UserID:   req.UserID,
		}

		response, err := cc.conversationUsecase.StartConversation(c.Request.Context(), startReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Failed to start conversation",
				Details: err.Error(),
			})
			return
		}

		// Convert to unified response format
		unifiedResponse := ConversationResponse{
			ConversationID:    response.ConversationID,
			Question:          &response.Question,
			IsComplete:        false,
			CurrentStep:       response.CurrentStep,
			TotalSteps:        response.TotalSteps,
			IsNewConversation: true,
		}

		c.JSON(http.StatusOK, unifiedResponse)
	} else {
		// Continuing an existing conversation
		if req.Answer == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Missing answer",
				Details: "answer is required for continuing a conversation",
			})
			return
		}

		answerReq := dto.SubmitAnswerRequest{
			ConversationID: req.ConversationID,
			Answer:         req.Answer,
		}

		response, err := cc.conversationUsecase.SubmitAnswer(c.Request.Context(), answerReq)
		if err != nil {
			// Check for specific error types
			if err.Error() == "conversation not found" {
				c.JSON(http.StatusNotFound, dto.ErrorResponse{
					Error:   "Conversation not found",
					Details: "The specified conversation ID does not exist",
				})
				return
			}

			if err.Error() == "conversation is not active" {
				c.JSON(http.StatusBadRequest, dto.ErrorResponse{
					Error:   "Conversation not active",
					Details: "This conversation has already been completed or expired",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Failed to submit answer",
				Details: err.Error(),
			})
			return
		}

		// Convert to unified response format
		unifiedResponse := ConversationResponse{
			ConversationID:    response.ConversationID,
			Question:          response.Question,
			Message:           response.Message,
			IsComplete:        response.IsComplete,
			CurrentStep:       response.CurrentStep,
			TotalSteps:        response.TotalSteps,
			IsNewConversation: false,
		}

		// If conversation is complete, get the report
		if response.IsComplete {
			reportResponse, err := cc.conversationUsecase.GetReport(c.Request.Context(), req.ConversationID)
			if err == nil {
				unifiedResponse.Report = reportResponse.Report
			}
		}

		c.JSON(http.StatusOK, unifiedResponse)
	}
}

// StartConversation handles POST /api/v1/conversation/start
func (cc *ConversationController) StartConversation(c *gin.Context) {
	var req dto.StartConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Note: Conversation service is designed for unauthenticated users
	// UserID is optional and can be provided in the request body if needed

	response, err := cc.conversationUsecase.StartConversation(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to start conversation",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// SubmitAnswer handles POST /api/v1/conversation/answer
func (cc *ConversationController) SubmitAnswer(c *gin.Context) {
	var req dto.SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	response, err := cc.conversationUsecase.SubmitAnswer(c.Request.Context(), req)
	if err != nil {
		// Check for specific error types
		if err.Error() == "conversation not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Conversation not found",
				Details: "The specified conversation ID does not exist",
			})
			return
		}

		if err.Error() == "conversation is not active" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Conversation not active",
				Details: "This conversation has already been completed or expired",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to submit answer",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetReport handles GET /api/v1/conversation/{id}/report
func (cc *ConversationController) GetReport(c *gin.Context) {
	conversationID := c.Param("id")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Conversation ID is required",
			Details: "Please provide a valid conversation ID in the URL path",
		})
		return
	}

	response, err := cc.conversationUsecase.GetReport(c.Request.Context(), conversationID)
	if err != nil {
		// Check for specific error types
		if err.Error() == "conversation not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Conversation not found",
				Details: "The specified conversation ID does not exist",
			})
			return
		}

		if err.Error() == "conversation is not complete" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Conversation not complete",
				Details: "This conversation has not been completed yet",
			})
			return
		}

		if err.Error() == "final report not found" {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Report not found",
				Details: "The final health report could not be generated",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to retrieve report",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

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

// HandleConversation handles both starting and continuing conversations in a single endpoint
// POST /api/v1/conversation
func (cc *ConversationController) HandleConversation(c *gin.Context) {
	var req dto.ConversationRequest
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
		unifiedResponse := dto.ConversationResponse{
			ConversationID:    response.ConversationID,
			Heading:           "Let's assess your symptoms",
			Subheading:        "I'll ask you a few questions to understand your condition better",
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
		unifiedResponse := dto.ConversationResponse{
			ConversationID:    response.ConversationID,
			Heading:           fmt.Sprintf("Question %d of %d", response.CurrentStep, response.TotalSteps),
			Subheading:        "Please provide more details",
			Question:          response.Question,
			Message:           response.Message,
			IsComplete:        response.IsComplete,
			CurrentStep:       response.CurrentStep,
			TotalSteps:        response.TotalSteps,
			IsNewConversation: false,
		}

		// If conversation is complete, get the report and remedy
		if response.IsComplete {
			reportResponse, err := cc.conversationUsecase.GetReport(c.Request.Context(), req.ConversationID)

			// Always create a remedy response
			var remedy *dto.RemedyResponse

			if err != nil {
				// If report generation failed, create a basic remedy
				fmt.Printf("Error getting report: %v\n", err)
				remedy = &dto.RemedyResponse{
					SessionID: generateSessionID(),
					Triage: dto.TriageResponse{
						Level:    "YELLOW",
						RedFlags: []string{},
						Message:  "Please consult a healthcare provider for personalized advice",
					},
					Content: &entities.GuidanceCard{
						SelfCare:      []string{"Rest", "Stay hydrated", "Monitor your symptoms"},
						OTCCategories: []entities.OTCCategory{},
						SeekCareIf:    []string{"Symptoms worsen", "New symptoms appear"},
						Disclaimer:    "This is general advice. Please consult a healthcare provider for personalized care.",
					},
				}
			} else if reportResponse.Report != nil && reportResponse.Report.Remedy != nil {
				// Use remedy from report
				remedy = &dto.RemedyResponse{
					SessionID: generateSessionID(),
					Triage: dto.TriageResponse{
						Level:    "GREEN", // Default level
						RedFlags: []string{},
						Message:  "Based on your symptoms",
					},
					Content: &entities.GuidanceCard{
						SelfCare:      reportResponse.Report.Remedy.SelfCare,
						OTCCategories: reportResponse.Report.Remedy.OTCCategories,
						SeekCareIf:    reportResponse.Report.Remedy.SeekCareIf,
						Disclaimer:    reportResponse.Report.Remedy.Disclaimer,
					},
				}
			} else {
				// No remedy in report, create a basic one
				fmt.Printf("No remedy found in report, creating basic remedy\n")
				remedy = &dto.RemedyResponse{
					SessionID: generateSessionID(),
					Triage: dto.TriageResponse{
						Level:    "GREEN",
						RedFlags: []string{},
						Message:  "Please consult a healthcare provider for personalized advice",
					},
					Content: &entities.GuidanceCard{
						SelfCare:      []string{"Rest", "Stay hydrated", "Monitor your symptoms"},
						OTCCategories: []entities.OTCCategory{},
						SeekCareIf:    []string{"Symptoms worsen", "New symptoms appear"},
						Disclaimer:    "This is general advice. Please consult a healthcare provider for personalized care.",
					},
				}
			}

			// Set the remedy in response
			unifiedResponse.Remedy = remedy

			// Set heading based on triage level
			if remedy.Triage.Level == "RED" {
				unifiedResponse.Heading = "⚠️ Medical Emergency Detected"
				unifiedResponse.Subheading = "Your symptoms require immediate medical attention"
				unifiedResponse.Message = "Please seek immediate medical care. Do not delay."
			} else {
				unifiedResponse.Heading = "Your Personalized Remedy"
				unifiedResponse.Subheading = "Based on your symptoms, here's what we recommend"
			}

			c.JSON(http.StatusOK, unifiedResponse)
		} else {
			c.JSON(http.StatusOK, unifiedResponse)
		}
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

// generateSessionID generates a unique session ID
func generateSessionID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// fallback to timestamp if random fails
		return fmt.Sprintf("session_%d", time.Now().UnixNano())
	}
	return "session_" + hex.EncodeToString(b)
}

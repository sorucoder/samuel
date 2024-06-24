package server

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sorucoder/samuel/internal/configuration"
	"github.com/sorucoder/samuel/internal/payloads"
	"github.com/sorucoder/samuel/internal/samuel"
)

// Responses
func respondAPISuccess(context *gin.Context, code int, response map[string]any) {
	if response == nil {
		context.Status(code)
		return
	}
	context.JSON(code, response)
}

func respondAPIError(context *gin.Context, code int, message string, err error) {
	response := make(map[string]any)
	response["error"] = message
	if configuration.Application.GetString("mode") == "development" {
		if err != nil {
			response["details"] = err.Error()
		}
	}

	context.AbortWithStatusJSON(code, response)
}

// Middleware
func handleAuthorizedAPIGroup(context *gin.Context) {
	authorization := context.GetHeader("Authorization")
	authorizationScheme, authorizationPayload, authorizationValid := strings.Cut(authorization, " ")
	if !authorizationValid {
		respondAPIError(context, http.StatusBadRequest, "missing authorization", nil)
		return
	} else if authorizationScheme != "Bearer" {
		respondAPIError(context, http.StatusBadRequest, "invalid authentication scheme", nil)
		return
	}

	sessionToken, errParseSessionToken := uuid.Parse(authorizationPayload)
	if errParseSessionToken != nil {
		respondAPIError(context, http.StatusBadRequest, "invalid session token", errParseSessionToken)
		return
	}
	user, session, errAuthenticateSession := samuel.AuthenticateSession(context, sessionToken)
	if errAuthenticateSession != nil {
		respondAPIError(context, http.StatusUnauthorized, "session expired", errAuthenticateSession)
		return
	}

	context.Set("user", user)
	context.Set("session", session)
}

func handleAdministratorAPIGroup(context *gin.Context) {
	currentUser := context.MustGet("user").(*samuel.User)

	if !currentUser.Is("administrator") {
		respondAPIError(context, http.StatusForbidden, "user not administrator", nil)
		return
	}

	currentAdministrator, errGetAdministrator := samuel.GetAdministratorByUser(context, currentUser)
	if errGetAdministrator != nil {
		respondAPIError(context, http.StatusInternalServerError, "cannot get administrator", errGetAdministrator)
		return
	}

	context.Set("administrator", currentAdministrator)
}

// Routes
func handleLogin(context *gin.Context) {
	authorization := context.GetHeader("Authorization")
	authorizationScheme, authorizationPayload, authorizationValid := strings.Cut(authorization, " ")
	if !authorizationValid {
		respondAPIError(context, http.StatusBadRequest, "missing authorization", nil)
		return
	} else if authorizationScheme != "Basic" {
		respondAPIError(context, http.StatusBadRequest, "invalid authorization scheme", nil)
		return
	}

	credentialsBytes, errDecode := base64.StdEncoding.DecodeString(authorizationPayload)
	if errDecode != nil {
		respondAPIError(context, http.StatusBadRequest, "invalid authentication data", errDecode)
		return
	}

	identity, password, credentialsValid := strings.Cut(string(credentialsBytes), ":")
	if !credentialsValid {
		respondAPIError(context, http.StatusBadRequest, "malformed credentials", nil)
		return
	}

	user, session, errAuthenticate := samuel.LoginUser(context, identity, password)
	if errAuthenticate != nil {
		respondAPIError(context, http.StatusUnauthorized, "invalid credentials", errAuthenticate)
		return
	}

	respondAPISuccess(context, http.StatusOK, map[string]any{
		"user":    user,
		"session": session,
	})
}

func handleLogout(context *gin.Context) {
	user := context.MustGet("user").(*samuel.User)
	session := context.MustGet("session").(*samuel.Session)

	errEndSession := samuel.LogoutUser(context, user, session)
	if errEndSession != nil {
		respondAPIError(context, http.StatusInternalServerError, "cannot end session", errEndSession)
		return
	}

	respondAPISuccess(context, http.StatusOK, nil)
}

func handleCreatePasswordChange(context *gin.Context) {
	var payload payloads.CreatePasswordChange
	errBindPayload := context.ShouldBind(&payload)
	if errBindPayload != nil {
		respondAPIError(context, http.StatusBadRequest, "malformed password change create data", errBindPayload)
		return
	}

	errCreatePasswordChange := samuel.CreatePasswordChange(context, payload.Email)
	if errCreatePasswordChange != nil {
		if errors.Is(errCreatePasswordChange, samuel.ErrPasswordChangeExists) {
			respondAPISuccess(context, http.StatusOK, nil)
		} else {
			respondAPIError(context, http.StatusInternalServerError, "cannot create password change", errCreatePasswordChange)
		}
		return
	}

	respondAPISuccess(context, http.StatusCreated, nil)
}

func handleFulfillPasswordChange(context *gin.Context) {
	passwordChangeToken, errParsePasswordChangeToken := uuid.Parse(context.Param("token"))
	if errParsePasswordChangeToken != nil {
		respondAPIError(context, http.StatusBadRequest, "malformed password change token", errParsePasswordChangeToken)
		return
	}

	var payload payloads.FulfillPasswordChange
	errBindPayload := context.ShouldBind(&payload)
	if errBindPayload != nil {
		respondAPIError(context, http.StatusBadRequest, "malformed fulfill password change data", errBindPayload)
		return
	}

	errChangePassword := samuel.FulfillPasswordChange(context, passwordChangeToken, payload.NewPassword)
	if errChangePassword != nil {
		respondAPIError(context, http.StatusInternalServerError, "cannot fulfill password change", errChangePassword)
		return
	}

	respondAPISuccess(context, http.StatusOK, nil)
}

func handlePing(context *gin.Context) {
	respondAPISuccess(context, http.StatusOK, nil)
}

func handleDashboard(context *gin.Context) {
	user := context.MustGet("user").(*samuel.User)
	session := context.MustGet("session").(*samuel.Session)

	response := map[string]any{
		"user":    user,
		"session": session,
	}
	switch user.Role().ID() {
	case "administrator":
		administrator, errGetAdministrator := samuel.GetAdministratorByUser(context, user)
		if errGetAdministrator != nil {
			respondAPIError(context, http.StatusInternalServerError, "cannot get administrator", errGetAdministrator)
			return
		}
		response["administrator"] = administrator
	case "instructor":
		instructor, errGetInstructor := samuel.GetInstructorByUser(context, user)
		if errGetInstructor != nil {
			respondAPIError(context, http.StatusInternalServerError, "cannot get instructor", errGetInstructor)
			return
		}
		response["instructor"] = instructor
	case "supervisor":
		supervisor, errGetSupervisor := samuel.GetSupervisorByUser(context, user)
		if errGetSupervisor != nil {
			respondAPIError(context, http.StatusInternalServerError, "cannot get supervisor", errGetSupervisor)
			return
		}
		response["supervisor"] = supervisor
	case "student":
		student, errGetStudent := samuel.GetStudentByUser(context, user)
		if errGetStudent != nil {
			respondAPIError(context, http.StatusInternalServerError, "cannot get student", errGetStudent)
			return
		}
		response["student"] = student
	}

	respondAPISuccess(context, http.StatusOK, response)
}

func handleViewAudit(context *gin.Context) {
	user := context.MustGet("user").(*samuel.User)
	session := context.MustGet("session").(*samuel.Session)
	administrator := context.MustGet("administrator").(*samuel.Administrator)

	var payload payloads.ViewAudit
	errBindPayload := context.ShouldBind(&payload)
	if errBindPayload != nil {
		respondAPIError(context, http.StatusBadRequest, "malformed audit view data", errBindPayload)
		return
	}

	auditBatch, errGetAuditBatch := samuel.GetAuditBatchByDate(context, payload.Date, payload.Page, payload.Count, payload.Sort, payload.Descending)
	if errGetAuditBatch != nil {
		respondAPIError(context, http.StatusInternalServerError, "cannot get audits", errGetAuditBatch)
		return
	}

	respondAPISuccess(context, http.StatusOK, map[string]any{
		"user":          user,
		"session":       session,
		"administrator": administrator,
		"batch":         auditBatch,
	})
}

func newRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	API := router.Group("/api")
	{
		API.GET("/login", handleLogin)

		passwordChangeAPI := API.Group("/password_change")
		{
			passwordChangeAPI.POST("/create", handleCreatePasswordChange)
			passwordChangeAPI.PUT("/fulfill/:token", handleFulfillPasswordChange)
		}

		authorizedAPI := API.Group("/", handleAuthorizedAPIGroup)
		{
			authorizedAPI.GET("/ping", handlePing)
			authorizedAPI.GET("/dashboard", handleDashboard)
			authorizedAPI.GET("/logout", handleLogout)

			administratorAPI := authorizedAPI.Group("/", handleAdministratorAPIGroup)
			{
				administratorAPI.GET("/audit/view", handleViewAudit)
			}
		}
	}

	return router
}

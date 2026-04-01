package rest

import (
	"ReservationsService/internal/core"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"log/slog"
	"net/http"
	"strconv"
)

type ServiceInterface interface {
	AddReservation(ctx context.Context, request *core.ReservationRequest) (*core.ReservationResponse, error)
	CloseReservation(ctx context.Context, request *core.ReturnRequest) (*core.ReturnResponse, error)
	CheckReservation(ctx context.Context, userId int) (*core.CheckReservResponse, error)
}

type RestReservServer struct {
	log     *slog.Logger
	service ServiceInterface
}

func NewRestReservServer(log *slog.Logger, service ServiceInterface) *RestReservServer {
	return &RestReservServer{
		log:     log,
		service: service,
	}
}

func (s *RestReservServer) RegisterRoutes(engine *gin.Engine) {
	engine.POST("/reservations", s.AddReservation)
	engine.PATCH("/reservations/:id/close", s.CloseReservation)
	engine.GET("/reservations", s.CheckReservation)
}

func (s *RestReservServer) AddReservation(ctx *gin.Context) {
	var req core.ReservationRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.log.Error("Bind json err: ", err)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  core.ErrInvalidInput,
		})
		return
	}

	resp, err := s.service.AddReservation(ctx.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{
				"status": "error",
				"error":  "books not found",
			})
			return

		case errors.Is(err, core.ErrInvalidInput):
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  "invalid input",
			})
			return

		case errors.Is(err, core.ErrUnauthorized):
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  "unauthorized",
			})
			return

		case errors.Is(err, core.ErrNoRights):
			ctx.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error":  "no rights",
			})
			return

		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  "internal server error",
			})
		}

		s.log.Error("book not added", "error", err)

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   resp,
	})
}

func (s *RestReservServer) CloseReservation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	req := core.ReturnRequest{
		ReservationID: id,
	}

	resp, err := s.service.CloseReservation(ctx.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrInvalidInput):
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  "invalid input",
			})
			return

		case errors.Is(err, core.ErrUnauthorized):
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  "unauthorized",
			})
			return

		case errors.Is(err, core.ErrNoRights):
			ctx.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error":  "no rights",
			})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  "internal server error",
			})
		}

		s.log.Error("book not added", "error", err)

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   resp,
	})
}

func (s *RestReservServer) CheckReservation(ctx *gin.Context) {
	var req core.CheckReservRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.log.Error("Bind json err: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  core.ErrInvalidInput,
		})
		return
	}

	resp, err := s.service.CheckReservation(ctx.Request.Context(), req.UserID)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrInvalidInput):
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  "invalid input",
			})
			return

		case errors.Is(err, core.ErrUnauthorized):
			s.log.Error("unauthorized user", "error", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  "unauthorized",
			})
			return

		case errors.Is(err, core.ErrNoRights):
			ctx.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error":  "no rights",
			})
			return

		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  "internal server error",
			})
		}

		s.log.Error("book not added", "error", err)

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   resp,
	})
}

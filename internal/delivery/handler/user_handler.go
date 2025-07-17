package handler

import (
	"net/http"

	"github.com/AndikaPrasetia/wash-shoe/internal/model"
	"github.com/AndikaPrasetia/wash-shoe/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserHandler struct {
    userUC usecase.UserUsecase
}

func NewUserHandler(userUC usecase.UserUsecase) *UserHandler {
    return &UserHandler{userUC: userUC}
}

func (h *UserHandler) Delete(c *gin.Context) {
    // get user from Context (set by middleware)
    currentUser, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "user not found in context",
        })
        return
    }

    authUser, ok := currentUser.(model.User)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "invalid user type",
        })
        return
    }

    // get a deleted user ID from path params
    userID := c.Param("id")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "user ID is required",
        })
    }

    // parse ID to UUID
    targetUUID, err := uuid.Parse(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "invalid user ID format",
        })
        return
    }
    // check only user or admin can delete 
    currentUUID, err := uuid.Parse(authUser.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "invalid current use ID",
        })
        return
    }

    // if not admin and not user to delete its own
    if authUser.Role != "admin" && targetUUID != currentUUID {
        c.JSON(http.StatusForbidden, gin.H{
            "error": "you don't have permission to delete this user",
        })
        return
    }

    // delete user from DB
    err = h.userUC.Delete(c.Request.Context(), pgtype.UUID{
        Bytes: targetUUID,
        Valid: true,
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "user deleted successfully",
    })
}


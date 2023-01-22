package handler

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	_ "github.com/vitalygudza/duty-app/docs"
	"github.com/vitalygudza/duty-app/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// Frontend: serve static files from the frontend/dist directory.
	router.Use(static.Serve("/", static.LocalFile("./frontend/dist", true)))

	// Backend: REST API
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/register", h.register)
		auth.POST("/login", h.login)
		auth.POST("/logout", h.logout)
	}

	api := router.Group("/api", h.userIdentity)
	{
		teams := api.Group("/teams")
		{
			teams.POST("/", h.createTeam)      // Create team
			teams.GET("/:id", h.readTeam)      // Read team
			teams.PUT("/:id", h.updateTeam)    // Update team
			teams.DELETE("/:id", h.deleteTeam) // Delete team
			teams.GET("/", h.listTeams)        // List teams

			teams.GET("/:id/history", h.historyDuties) // Read team duties history

			teammates := teams.Group(":id/teammates")
			{
				teammates.POST("/", h.createTeammate) // Create teammate
				teammates.GET("/", h.listTeammates)   // List teammates
			}

			duties := teams.Group(":id/duties")
			{
				duties.POST("/", h.createDuty)   // Create duty for teammate
				duties.GET("/", h.listDuties)    // List duties
				duties.GET("/now", h.readDuties) // Read current duties (daily and weekly)
				// duties.POST("/notify-daily", h.notifyDailyDuty) // Create notification for daily duty: TODO: for what this method?
			}

			notifications := teams.Group(":id/notifications")
			{
				notifications.POST("/test", h.sendTestNotification) // Send test notification to team's Space channel
			}
		}

		teammates := api.Group("/teammates")
		{
			teammates.GET("/:id", h.readTeammate)      // Read teammate
			teammates.PUT("/:id", h.updateTeammate)    // Update teammate
			teammates.DELETE("/:id", h.deleteTeammate) // Delete teammate
		}

		duties := api.Group("/duties")
		{
			duties.PUT("/:id", h.updateDuty)    // Update duty
			duties.DELETE("/:id", h.deleteDuty) // Delete duty
		}
	}
	return router
}

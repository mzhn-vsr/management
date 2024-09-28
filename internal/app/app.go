package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"mzhn/management/internal/config"
	"mzhn/management/internal/handlers"
	"mzhn/management/internal/services/chatservice"
	"mzhn/management/internal/services/faqservice"
	"mzhn/management/internal/services/feedbackservice"

	"github.com/labstack/echo/v4"
	emw "github.com/labstack/echo/v4/middleware"
)

type App struct {
	app *echo.Echo
	cfg *config.Config

	faqsvc      *faqservice.FaqService
	chatsvc     *chatservice.ChatService
	feedbacksvc *feedbackservice.FeedbackService
}

func newApp(cfg *config.Config, faqsvc *faqservice.FaqService, chatsvc *chatservice.ChatService, fdsvc *feedbackservice.FeedbackService) *App {
	return &App{
		app:         echo.New(),
		cfg:         cfg,
		faqsvc:      faqsvc,
		chatsvc:     chatsvc,
		feedbacksvc: fdsvc,
	}
}

func (a *App) initApp() {
	a.app.Use(emw.Logger())
	a.app.Use(emw.CORSWithConfig(emw.CORSConfig{
		AllowOrigins:     strings.Split(a.cfg.CORS.Origins, ","),
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
		AllowCredentials: true,
	}))

	// endpoints
	a.app.POST("/faq", handlers.CreateFaq(a.faqsvc))
	a.app.PUT("/faq", handlers.UpdateFaq(a.faqsvc))
	a.app.GET("/faq", handlers.ListFaq(a.faqsvc))
	a.app.GET("/faq/:id", handlers.FindFaq(a.faqsvc))
	a.app.DELETE("/faq/:id", handlers.DeleteFaq(a.faqsvc))

	a.app.POST("/invoke", handlers.Invoke(a.chatsvc))
	a.app.POST("/predict", handlers.Predict(a.chatsvc))

	a.app.PUT("/feedback", handlers.SendFeedback(a.feedbacksvc))
	a.app.GET("/feedback", handlers.FeedbackStats(a.feedbacksvc))
}

func (a *App) Run() {

	a.initApp()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(
		sigChan,
		os.Interrupt,
		syscall.SIGTERM,
	)

	go func() {
		port := a.cfg.App.Port
		addr := fmt.Sprintf(":%d", port)
		slog.Info("running server", slog.String("addr", addr))
		a.app.Start(addr)
	}()

	sig := <-sigChan
	slog.Info(fmt.Sprintf("Signal %v received, stopping server...\n", sig))
	a.app.Shutdown(context.Background())
}

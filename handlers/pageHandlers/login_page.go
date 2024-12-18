package pageHandlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/derekadair/go-templ-htmx-example-app/handlers"
	"github.com/derekadair/go-templ-htmx-example-app/services"
	"github.com/derekadair/go-templ-htmx-example-app/ui/components"
	"github.com/derekadair/go-templ-htmx-example-app/ui/pages"
	"github.com/labstack/echo/v4"
)

type LoginPageHandler struct {
	ctx    context.Context
	logger *slog.Logger

	*services.ServiceMesh
}

func NewLoginPageHandler(ctx context.Context, logger *slog.Logger, serviceMesh *services.ServiceMesh) *LoginPageHandler {
	return &LoginPageHandler{
		ctx:         ctx,
		logger:      logger,
		ServiceMesh: serviceMesh,
	}
}

func (obj *LoginPageHandler) RegisterPublicRoutes(echoHandler *echo.Echo) {
	echoHandler.GET("/login", obj.BasePage)
	echoHandler.POST("/login/validate-inputs", obj.ValidateInputs)
	echoHandler.POST("/login/submit", obj.Submit)
}

func (obj *LoginPageHandler) BasePage(echoCtx echo.Context) error {
	page := pages.LoginPage(pages.LoginInputFormParams{})
	handlers.Render(echoCtx, &page)
	return nil
}

func (obj *LoginPageHandler) FormToParams(echoCtx echo.Context) pages.LoginInputFormParams {
	email := echoCtx.FormValue("email")
	password := echoCtx.FormValue("password")

	submitButtonDisabled := false
	if email == "" || password == "" {
		submitButtonDisabled = true
	}

	params := pages.LoginInputFormParams{
		Email:                email,
		Password:             password,
		SubmitButtonDisabled: submitButtonDisabled,
	}
	return params
}

func (obj *LoginPageHandler) ValidateInputs(echoCtx echo.Context) error {
	params := obj.FormToParams(echoCtx)
	component := components.Button("Login", params.SubmitButtonDisabled)

	handlers.Render(echoCtx, &component)
	return nil
}

func (obj *LoginPageHandler) Submit(echoCtx echo.Context) error {
	params := obj.FormToParams(echoCtx)

	if params.FormAppearsValid() {
		token, err := obj.UserAuthenticationService.Signin(
			echoCtx.Request().Context(),
			params.Email,
			params.Password,
		)

		if err == nil {
			echoCtx.SetCookie(&http.Cookie{
				Name:     "token",
				Value:    token,
				Expires:  time.Now().Add(services.TOKEN_DURATION),
				HttpOnly: true,
				Path:     "/r/",
				Secure:   false,
			})
			echoCtx.Response().Header().Add("HX-Redirect", "/r/")
			return nil

		} else {
			params.ShowFailedLoginFlag = true

		}
	}

	component := pages.LoginInputForm(params)
	handlers.Render(echoCtx, &component)
	return nil
}

package configure

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/derekadair/go-templ-htmx-example-app/handlers"
	"github.com/derekadair/go-templ-htmx-example-app/handlers/componentHandlers"
	"github.com/derekadair/go-templ-htmx-example-app/handlers/pageHandlers"
	"github.com/derekadair/go-templ-htmx-example-app/services"
	"github.com/derekadair/go-templ-htmx-example-app/settings"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func ConfigureUIHandlers(ctx context.Context, echoHandler *echo.Echo, logger *slog.Logger, configuration *ServiceConfiguration) {

	// alive check
	echoHandler.GET("/alive", func(echoCtx echo.Context) error {
		return echoCtx.String(http.StatusOK, "Service Alive...")
	})

	var publicHandler []handlers.PublicUIHandler = []handlers.PublicUIHandler{
		pageHandlers.NewHomePageHandler(ctx, logger),
		pageHandlers.NewSignupPageHandler(ctx, logger, &configuration.ServiceMesh),
		pageHandlers.NewLoginPageHandler(ctx, logger, &configuration.ServiceMesh),
	}

	for _, handler := range publicHandler {
		handler.RegisterPublicRoutes(echoHandler)
	}

	// restricted endpoint handlers
	echoRestrictedHandler := echoHandler.Group("/r")

	jwtConfig := echojwt.Config{
		SigningKey: []byte(settings.JWT_SECRET_KEY),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(services.JWTScopeClaims)
		},
		TokenLookup: "cookie:token",
	}
	echoRestrictedHandler.Use(
		echojwt.WithConfig(jwtConfig),
	)

	var componentHandlerMesh componentHandlers.ComponentHandlerMesh = componentHandlers.NewComponentHandlerMesh(ctx, logger, &configuration.ServiceMesh)
	componentHandlerMesh.RegisterProtectedRoutes(echoRestrictedHandler)

	var protectedHandler []handlers.ProtectedUIHandler = []handlers.ProtectedUIHandler{
		pageHandlers.NewAppHandler(ctx, logger, &configuration.ServiceMesh, &componentHandlerMesh),
	}

	for _, handler := range protectedHandler {
		handler.RegisterProtectedRoutes(echoRestrictedHandler)
	}

}

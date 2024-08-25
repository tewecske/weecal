package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	project "weecal"
	"weecal/internal/config"
	"weecal/internal/handlers"
	"weecal/internal/hash/passwordhash"
	m "weecal/internal/middleware"
	database "weecal/internal/store/db"
	"weecal/internal/store/session"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/unrolled/secure"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg := config.LoadConfig()

	passwordHash := passwordhash.NewHPasswordHash()

	dbAccess := database.SetupDB(cfg.DatabaseName, passwordHash)

	dbAccess.UserStore.CreateUser("a@a.a", "aaa")
	dbAccess.SessionStore.CreateSession(&session.Session{
		UserID: 1,
	})

	// TODO: Check: base-uri 'none'; object-src 'none';
	// TODO: Check: script-src 'strict-dynamic' 'unsafe-inline' 'unsafe-eval'
	// TODO: Use: script-src 'report-sample' and report-uri /_/_/csp_report
	// w.Header().Set("Content-Security-Policy", cspHeader)
	// Use this for testing CSP
	// w.Header().Set("Content-Security-Policy-Report-Only", cspHeader)
	secureMiddleware := secure.New(secure.Options{
		AllowedHosts:            []string{"gothhost"},                            // AllowedHosts is a list of fully qualified domain names that are allowed. Default is empty list, which allows any and all host names.
		AllowedHostsAreRegex:    false,                                           // AllowedHostsAreRegex determines, if the provided AllowedHosts slice contains valid regular expressions. Default is false.
		AllowRequestFunc:        nil,                                             // AllowRequestFunc is a custom function type that allows you to determine if the request should proceed or not based on your own custom logic. Default is nil.
		HostsProxyHeaders:       []string{"X-Forwarded-Hosts"},                   // HostsProxyHeaders is a set of header keys that may hold a proxied hostname value for the request.
		SSLRedirect:             true,                                            // If SSLRedirect is set to true, then only allow HTTPS requests. Default is false.
		SSLTemporaryRedirect:    false,                                           // If SSLTemporaryRedirect is true, the a 302 will be used while redirecting. Default is false (301).
		SSLHost:                 "gothhost",                                      // SSLHost is the host name that is used to redirect HTTP requests to HTTPS. Default is "", which indicates to use the same host.
		SSLHostFunc:             nil,                                             // SSLHostFunc is a function pointer, the return value of the function is the host name that has same functionality as `SSLHost`. Default is nil. If SSLHostFunc is nil, the `SSLHost` option will be used.
		SSLProxyHeaders:         map[string]string{"X-Forwarded-Proto": "https"}, // SSLProxyHeaders is set of header keys with associated values that would indicate a valid HTTPS request. Useful when using Nginx: `map[string]string{"X-Forwarded-Proto": "https"}`. Default is blank map.
		STSSeconds:              31536000,                                        // STSSeconds is the max-age of the Strict-Transport-Security header. Default is 0, which would NOT include the header.
		STSIncludeSubdomains:    true,                                            // If STSIncludeSubdomains is set to true, the `includeSubdomains` will be appended to the Strict-Transport-Security header. Default is false.
		STSPreload:              true,                                            // If STSPreload is set to true, the `preload` flag will be appended to the Strict-Transport-Security header. Default is false.
		ForceSTSHeader:          false,                                           // STS header is only included when the connection is HTTPS. If you want to force it to always be added, set to true. `IsDevelopment` still overrides this. Default is false.
		FrameDeny:               true,                                            // If FrameDeny is set to true, adds the X-Frame-Options header with the value of `DENY`. Default is false.
		CustomFrameOptionsValue: "SAMEORIGIN",                                    // CustomFrameOptionsValue allows the X-Frame-Options header value to be set with a custom value. This overrides the FrameDeny option. Default is "".
		ContentTypeNosniff:      true,                                            // If ContentTypeNosniff is true, adds the X-Content-Type-Options header with the value `nosniff`. Default is false.
		BrowserXssFilter:        true,                                            // If BrowserXssFilter is true, adds the X-XSS-Protection header with the value `1; mode=block`. Default is false.
		// CustomBrowserXssValue:   "1; report=https://gothhost/xss-report",      // CustomBrowserXssValue allows the X-XSS-Protection header value to be set with a custom value. This overrides the BrowserXssFilter option. Default is "".
		ContentSecurityPolicy: "default-src 'self'; " +
			"style-src $NONCE; " + // 'unsafe-inline' would be needed for htmx but it doesn't work if a nonce or a hash is already defined
			"script-src $NONCE 'sha256-bUqqSw0+i0yR+Nl7kqNhoZsb1FRN6j9mj9w+YqY5ld8=' 'sha256-EtDJKiu1jHe6jtwOCABcdSkppIaCP/+vBbsOPG/numY=' 'sha256-NY2a+7GrW++i9IBhowd25bzXcH9BCmBrqYX5i8OxwDQ=' 'unsafe-eval'", // ContentSecurityPolicy allows the Content-Security-Policy header value to be set with a custom value. Default is "". Passing a template string will replace `$NONCE` with a dynamic nonce value of 16 bytes for each request which can be later retrieved using the Nonce function. alpine.js requires unsafe-eval by default without the csp build and that's used for now
		ReferrerPolicy: "same-origin", // ReferrerPolicy allows the Referrer-Policy header with the value to be set with a custom value. Default is "".
		// FeaturePolicy:           "vibrate 'none';",                               // Deprecated: this header has been renamed to PermissionsPolicy. FeaturePolicy allows the Feature-Policy header with the value to be set with a custom value. Default is "".
		// PermissionsPolicy:       "fullscreen=(), geolocation=()",                 // PermissionsPolicy allows the Permissions-Policy header with the value to be set with a custom value. Default is "".
		CrossOriginOpenerPolicy: "same-origin", // CrossOriginOpenerPolicy allows the Cross-Origin-Opener-Policy header with the value to be set with a custom value. Default is "".

		IsDevelopment: project.IsDevelopment, // This will cause the AllowedHosts, SSLRedirect, and STSSeconds/STSIncludeSubdomains options to be ignored during development. When deploying to production, be sure to set this to false.
	})

	authMiddleware := m.NewAuthMiddleware(dbAccess.SessionStore, cfg.SessionCookieName)

	router := chi.NewRouter()

	router.Handle("/static/*", http.StripPrefix("/static", project.Public()))

	router.Group(func(r chi.Router) {
		r.Use(
			middleware.Logger,
			secureMiddleware.Handler,
			m.TextHTMLMiddleware, // NOTE: it probably won't always be text/html
			authMiddleware.AddUserToContext,
		)

		r.NotFound(handlers.Make(handlers.HandleNotFound))

		r.Get("/", handlers.Make(handlers.HandleHome))

		r.Get("/login", handlers.Make(handlers.HandleLogin))

		r.Post("/login", handlers.Make(handlers.HandlePostLogin(
			dbAccess.UserStore,
			dbAccess.SessionStore,
			passwordHash,
			cfg.SessionCookieName,
		)))

		r.Post("/logout", handlers.Make(handlers.HandlePostLogout(
			cfg.SessionCookieName,
		)))

		r.Get("/calendar", handlers.Make(handlers.HandleCalendar))
	})
	// slog.Info("HTTP server started", "listenAddr", cfg.Port)
	// http.ListenAndServe(cfg.Port, router)

	killSig := make(chan os.Signal, 1)

	signal.Notify(killSig, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: router,
	}

	go func() {
		err := srv.ListenAndServe()

		// HTTPS
		// To generate a development cert and key, run the following from your *nix terminal:
		// go run $GOROOT/src/crypto/tls/generate_cert.go --host="localhost"
		// err := srv.ListenAndServeTLS(":8443", "cert.pem", "key.pem", app)

		if errors.Is(err, http.ErrServerClosed) {
			logger.Info("Server shutdown complete")
		} else if err != nil {
			logger.Error("Server error", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	logger.Info("Server started", slog.String("port", cfg.Port))
	<-killSig

	logger.Info("Shutting down server")

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbCloseError := dbAccess.DB.Close()
	if dbCloseError != nil {
		logger.Error("DB close failed", slog.Any("err", dbCloseError))
	} else {
		logger.Info("DB closed")
	}

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", slog.Any("err", err))
		os.Exit(1)
	}

	logger.Info("Server shutdown complete")
}

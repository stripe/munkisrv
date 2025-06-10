package main

import (
	"context"
	"crypto"
	"embed"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/stripe/munkisrv/config"
	"github.com/stripe/munkisrv/keyutils"
	"github.com/stripe/munkisrv/munkirepo"

	"github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Load the config
	configPath := flag.String("c", "../../config/config.yaml", "path to config.yaml")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("failed to load %s, err: %s\n", *configPath, err)
	}

	// Create signer for CloudFront signed URLs
	privKey, err := keyutils.ParsePrivateKey([]byte(cfg.Cloudfront.PrivateKey), "URL signer private key")
	if err != nil {
		log.Fatalf("Failed to load private key, err: %s\n", err.Error())
	}
	signerKey, ok := privKey.(crypto.Signer)
	if !ok {
		log.Fatal("private key is not a crypto.Signer\n")
	}
	signer := sign.NewURLSigner(cfg.Cloudfront.KeyID, signerKey)

	// Setup http server
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/healthz", healthz(munkirepo.Repo))
	r.Get("/repo/*", munkiRepoFunc)
	r.Get("/repo/pkgs/*", munkiPkgFunc(cfg.Cloudfront.URL, signer))
	r.Head("/repo/pkgs/*", munkiPkgFunc(cfg.Cloudfront.URL, signer))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// Start http server
	server := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: r,
	}

	go func() {
		fmt.Printf("Starting server on port %s...\n", cfg.Server.Port)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Server shutdown complete.")
}

func munkiRepoFunc(w http.ResponseWriter, r *http.Request) {
	fs := http.StripPrefix("/repo/", http.FileServer(http.FS(munkirepo.Repo)))
	fs.ServeHTTP(w, r)
}

func munkiPkgFunc(cloudFrontURL string, signer *sign.URLSigner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(cloudFrontURL)
		if err != nil {
			http.Error(w, "failed to parse base url", http.StatusInternalServerError)
			return
		}
		u.Path = path.Join(u.Path, r.URL.Path)
		finalURL := u.String()

		signedURL, err := signer.Sign(finalURL, time.Now().Add(1*time.Hour))
		if err != nil {
			http.Error(w, "Failed to sign url", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, signedURL, http.StatusTemporaryRedirect)
	}
}

func healthz(repo embed.FS) http.HandlerFunc {
	var healthy bool
	if _, err := repo.Open("catalogs/all"); err == nil {
		healthy = true
	} else {
		log.Printf("healthcheck failed with %s\n", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if !healthy {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

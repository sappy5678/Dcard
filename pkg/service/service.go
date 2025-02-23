// Copyright 2017 Emir Ribic. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// dcard - Go(lang) restful starter kit
//
// API Docs for dcard v1
//
//		 Terms Of Service:  N/A
//	    Schemes: http
//	    Version: 2.0.0
//	    License: MIT http://opensource.org/licenses/MIT
//	    Contact: Emir Ribic <ribice@gmail.com> https://ribice.ba
//	    Host: localhost:8080
//
//	    Consumes:
//	    - application/json
//
//	    Produces:
//	    - application/json
//
//	    Security:
//	    - bearer: []
//
//	    SecurityDefinitions:
//	    bearer:
//	         type: apiKey
//	         name: Authorization
//	         in: header
//
// swagger:meta
package service

import (
	"net/http"
	"os"

	"github.com/labstack/echo"

	"github.com/sappy5678/dcard/pkg/service/shorturl"
	sl "github.com/sappy5678/dcard/pkg/service/shorturl/logservice"
	st "github.com/sappy5678/dcard/pkg/service/shorturl/transport"
	"github.com/sappy5678/dcard/pkg/utl/config"
	"github.com/sappy5678/dcard/pkg/utl/postgres"
	"github.com/sappy5678/dcard/pkg/utl/redis"
	"github.com/sappy5678/dcard/pkg/utl/server"
	"github.com/sappy5678/dcard/pkg/utl/zlog"
)

// Start starts the API service
func Start(cfg *config.Configuration) error {
	db, err := postgres.New(os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}

	redisClient, err := redis.New(os.Getenv("REDIS_URL"))
	if err != nil {
		return err
	}

	log := zlog.New()

	host := "http://localhost:8080" // should get from central config service
	machineID := uint64(1)          // should get from central config service

	e := server.New()
	rootGroup := e.Group("")
	st.NewHTTP(sl.New(shorturl.Initialize(machineID, host, db, redisClient), log), rootGroup)

	rootGroup.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "health",
		})
	})
	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})

	return nil
}

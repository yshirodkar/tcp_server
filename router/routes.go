/**
 * Copyright 2018 InsideSales.com Inc.
 * All Rights Reserved.
 *
 * NOTICE: All information contained herein is the property of InsideSales.com, Inc. and its suppliers, if
 * any. The intellectual and technical concepts contained herein are proprietary and are protected by
 * trade secret or copyright law, and may be covered by U.S. and foreign patents and patents pending.
 * Dissemination of this information or reproduction of this material is strictly forbidden without prior
 * written permission from InsideSales.com Inc.
 *
 * Requests for permission should be addressed to the Legal Department, InsideSales.com,
 * 1712 South East Bay Blvd. Provo, UT 84606.
 *
 * The software and any accompanying documentation are provided "as is" with no warranty.
 * InsideSales.com, Inc. shall not be liable for direct, indirect, special, incidental, consequential, or other
 * damages, under any theory of liability.
 */

package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"philo_server/common"
	"philo_server/controllers"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//The CustomRouter struct extends Gorilla's Mux Router to allow for customizing a router
//Additional functions/data members can be added to CustomRouter if you want additional functionality
type CustomRouter struct {
	*mux.Router
}

type Routes []Route

func NewRouter(
	configGetter common.IConfigGetter,
	healthController common.IHealthController,
	stackPopController controllers.IStackPopController,
) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	// This is server profiler.
	router.PathPrefix("/debug/").Handler(http.DefaultServeMux)
	// This is server health check end point
	DefineHealthRoutes(router, healthController)
	// This is server controller end point
	// DefineKpmRoutes(router, kpmController)
	return router
}

func DefineHealthRoutes(router *mux.Router, healthController common.IHealthController) {
	router.Methods("GET").Path("/health").HandlerFunc(healthController.GetHealth)
}

// func DefineKpmRoutes(
// 	jwtMiddleware IJwtMiddleware,
// 	router *mux.Router,
// 	kpmController controllers.IKpmController,
// ) *mux.Router {

// 	subRouter := CustomRouter{
// 		router.PathPrefix("/kpms").Subrouter(),
// 	}
// 	subRouter.RegisterRoute(Route{"Get stack end point", "GET", "", kpmController.GetAllKpms})
// 	return router
// }

//The RegisterRoute function is a custom function that wraps the passed in route in any middleware that has been
//added to the router
//The order that the middleware will be applied to an endpoint is Last In First Out
//For example: The first MiddleWare in the array passed into the Use function, will be the last middleware
//applied to a given route
func (r *CustomRouter) RegisterRoute(route Route) {
	handler := route.HandlerFunc
	// for _, mw := range r.mw {
	// 	handler = mw(handler)
	// }

	r.Methods(route.Method).Path(route.Pattern).Handler(handler)
}
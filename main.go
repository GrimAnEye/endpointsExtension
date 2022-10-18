/*
EndpointsExtension - an extension over Telegram-Bot-Api v5
that adds endpoints and middleware, similar to https://gopkg.in/telebot.v3
*/
package endpointsExtension

import (
	tgb "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (

	// Router - a structure that contains handlers and performs routing of received updates.
	Router struct {

		// Array of endpoints and actions
		endpoints []endpoint

		// An array of middleware.
		// The top-level context is considered the first group
		group *Group
	}

	// Context - a structure containing a copy of the update,
	// a link to the bot and custom fields passed inside the handlers
	Context struct {
		// Copy of data from the update channel
		U tgb.Update

		// Link to the bot to send messages
		B *tgb.BotAPI

		// Custom Field - A map of custom fields to pass to endpoints via Context.
		// Because it is an interface, you must always check the typeof a
		// variable when retrieving and storing data.
		CustomField map[string]interface{}
	}

	// endpoint - structure from an action and a condition that forms an endpoint
	endpoint struct {
		condition interface{} // Text or function to check endpoint execution conditions
		action    HandleFunc  // Endpoint action to be taken when the condition is successfully checked
	}

	// Group - wrapper for grouping endpoints and middleware
	Group struct {
		r          *Router
		middleware []MiddlewareFunc
	}

	// HandleFunc - action function to be called if the condition is successfully checked
	HandleFunc func(*Context) error

	// MiddlewareFunc - wrapper over HandleFunc,
	// executed after the endpoint condition, but before the main function
	MiddlewareFunc func(HandleFunc) HandleFunc
)

// NewRouter - creates a new, empty router
func NewRouter() Router {
	return Router{
		endpoints: make([]endpoint, 0),
		group:     &Group{},
	}
}

/*
Handler - checks the type and adds it to the endpoints
array for further processing. In case of error, causes panic
*/
func (r *Router) Handler(condition interface{}, h HandleFunc, m ...MiddlewareFunc) {

	// Add global middleware to the general list
	if len(r.group.middleware) > 0 {
		m = append(r.group.middleware, m...)
	}

	// Packing an action into layers of middleware
	handler := func(c *Context) error {
		return applyMiddleWare(h, m...)(c)
	}

	// Condition type check
	switch condition.(type) {

	case string, func(r *Context) bool:
		r.endpoints = append(r.endpoints, endpoint{
			condition: condition,
			action:    handler,
		})

	default:
		panic("invalid endpoint")
	}
}

// applyMiddleWare - gets the action function and wraps it with middleware functions in order of addition
func applyMiddleWare(h HandleFunc, middleware ...MiddlewareFunc) HandleFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

// Use - adds middleware to the bot's global chain.
func (r *Router) Use(middleware ...MiddlewareFunc) {
	r.group.Use(middleware...)
}

// Use - adds middleware to the group.
// Executed in order of addition, but before the main function
func (g *Group) Use(middleware ...MiddlewareFunc) {
	g.middleware = append(g.middleware, middleware...)
}

// Handler - checks the type and adds it to the endpoints
// array for further processing. In case of error, causes panic
func (g *Group) Handler(c interface{}, h HandleFunc, m ...MiddlewareFunc) {
	g.r.Handler(c, h, append(g.middleware, m...)...)
}

// Group - creates a new routing group, delegating upper layer middleware to it
func (r *Router) Group() *Group {
	return &Group{r: r}
}

// Group - like *router.Group, creates a new routing group by delegating top and parent level middleware to it.
func (g *Group) Group() *Group {
	return &Group{
		r:          g.r,
		middleware: g.middleware,
	}
}

/*
Route - starts the verification process for all endpoints,
and if a match is found, it stops the enumeration of conditions,
for performing the assigned function.
*/
func (r Router) Route(c Context) error {
	for _, h := range r.endpoints {

		// Check endpoint condition
		var allow bool
		switch endpoint := h.condition.(type) {
		case string:
			if c.U.Message != nil {
				allow = c.U.Message.Text == endpoint
			}
		case func(c *Context) bool:
			allow = endpoint(&c)
		}

		// If the condition is true, execute the endpoint and break the loop
		if allow {
			if err := h.action(&c); err != nil {
				return err
			}
			break
		}
	}
	return nil
}

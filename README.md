# EndpointExtension
This package is an extension for [telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api), adding to it the ability to use endpoints, as in [go-telebot](https://github.com/go-telebot/telebot) (it also served as the code base for the package).

## Expansion Composition
The package adds the following components:
1. 	**Context** - a structure containing a copy of the update, a link to the bot and custom fields passed inside the handlers
```
Context struct {
	U tgb.Update
	B *tgb.BotAPI
	CustomField map[string]interface{}
}
```
After creating a bot, every time you receive an update from Telegram, you need to update the context and start routing to endpoints.

Custom fields are used to transfer custom fields and data, through the context, into handlers for further use.
The most important thing in this case is to perform type checking, so as not to be caught when performing a fatalpanic:
```
// Custom structure
type TUser struct{
	Name string
	Age int
}

// Create a context
	context := e.Context{U: update,B: bot, CustomFields: make(map[string]interface{})}

// Passing custom struct to context
	context.CustomField["data"] = TUser{Name: "Bob", Age: 18}

// Passing the context to the router
	go func(){
		router.Route(context)
	}()

...

// Inside the handler I use this data
func(c *ee.Context) error {
	fmt.Println(c.CustomField["data"].(TUser))
	return nil
}
```

2. **Endpoints** - when starting routing, the router goes through the list of endpoints, checks the conditions. If the condition evaluates to true, stops the iteration of the conditions and performs the assigned action. A condition can be either a function or a message(update.Message.Text).
```
// Creating an empty router
router:=ee.NewRouter()

// Adding an endpoint.
// Instead of a function, there can be a string and then the check will be with update.Message.Text
router.Handler(

// Endpoint execution condition
func(c *ee.Context)bool{
	return c.Message!=nil
},

// Action to take
func(c *ee.Context)error{
	msg:=tgbotapi.NewMessage(c.U.FromChat().ID,"Hello!")
	c.B.Send(msg)
	return nil
	},
)

// Requesting updates from Telegram
updates:=bot.GetUpdateChan(tgbotapi.UpdateConfig{Offset:0,Timeout:60})
for update:=range updates{

	// create a context
	context:=ee.Context{U: update,B:bot,CustomFields: make(map[string]interface{})}

	// Starting routing in a goroutine so as not to slow down other updates
	go func(c ee.Context){
		router.Route(c)
	}(context)
}

```

3. Middleware - executed after the condition but before the main action. Execution occurs in the order of adding it to the handlers - first the earliest, then the latest.
```
func(next e.HandleFunc) e.HandleFunc {
	return func(c *e.Context) error {
		fmt.Println(c.U.Message.Text)

		return next(c)
	}
}
```
Has multiple places to add:
```
// In the router:
router.Use(e.MiddlewareFunc)

// When adding a handler:
router.Handler(
	func(c *ee.Context),
	ee.HandlerFunc,
	ee.MiddlewareFunc,
)

// When grouping handlers by analogy with a router:
group.Use(ee.MiddlewareFunc)

router.Handler(
	func(c *ee.Context),
	ee.HandlerFunc,
	ee.MiddlewareFunc,
)
```


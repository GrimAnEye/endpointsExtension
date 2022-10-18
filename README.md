# EndpointExtension
This package is an extension for [telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api), adding to it the ability to use endpoints, as in [go-telebot](https://github.com/go-telebot/telebot) (it also served as the code base for the package).

## Expansion Composition
The package adds the following components:
1. Context - a general structure containing a link to the bot, a copy of the current update and an array of custom fields to pass inside the handlers.
```
func NewContext(u tgb.Update, b *tgb.BotAPI, customFields map[string]interface{}) *Context {
	return &Context{
		U:           u,
		B:           b,
		CustomField: customFields,
		endpoints:   make([]endpoint, 0),
		group:       &Group{},
	}
```
After creating a bot, every time you receive an update from Telegram, you need to update the context and start routing to endpoints. For example:
```
// Creating a simple bot
bot, _ := tgb.NewBotAPI(Settings.Bot.Token)

// Creating an empty router (in fact - a context)
router :=endpointExtension.NewContext(tgb.Update{},bot,make(map[string]interface{}))
```

2. 

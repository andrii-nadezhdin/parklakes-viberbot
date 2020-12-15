package handlers

import (
	"context"

	"github.com/triviy/parklakes-viberbot/application/commands"
	"github.com/triviy/parklakes-viberbot/application/integrations/google"
	"github.com/triviy/parklakes-viberbot/infrastructure/persistance"
	"github.com/triviy/parklakes-viberbot/web/config"
)

// Handlers collection
type Handlers struct {
	*MigrateCarOwnersHandler
	*CallbackHandler
	*SetWebhookHandler
	*HealthCheckHandler
}

// InitializeHandlers creates handlers with all dependencies
func InitializeHandlers(ctx context.Context, cfg *config.APIConfig) (h *Handlers, err error) {
	datastore, err := persistance.NewMongoDatastore(ctx, cfg.GetDBConnectionString())
	if err != nil {
		return
	}

	carOwnersRepo := persistance.NewCarOwnersRepo(datastore)
	carOwnerPropsRepo := persistance.NewCarOwnerPropsRepo(datastore)
	subscribersRepo := persistance.NewSubscribersRepo(datastore)

	gSpreadsheet, err := google.NewSpreadsheet(ctx, cfg.GetSheetsAPIKey(), cfg.GetSheetsAPISpreadsheetID())
	if err != nil {
		return
	}

	migrateCarOwnerCmd := commands.NewMigrateCarOwnersCmd(carOwnersRepo, carOwnerPropsRepo, gSpreadsheet)
	migrateCarOwnerHandler := NewMigrateCarOwnersHandler(migrateCarOwnerCmd)

	getCarOwnerByTextCmd := commands.NewGetCarOwnerByTextCmd(cfg, carOwnersRepo)
	updateSubscriberCmd := commands.NewUpdateSubscriberCmd(subscribersRepo)
	unsubscribeCmd := commands.NewUnsubscribeCmd(subscribersRepo)
	welcomeCmd := commands.NewWelcomeCmd()
	callbackHandler := NewCallbackHandler(getCarOwnerByTextCmd, updateSubscriberCmd, unsubscribeCmd, welcomeCmd)

	setWebhookCmd := commands.NewSetWebhookCmd(cfg)
	setWebhookHandler := NewSetWebhookHandler(setWebhookCmd)

	healthCheckHandler := NewHealthCheckHandler(datastore)

	h = &Handlers{
		migrateCarOwnerHandler,
		callbackHandler,
		setWebhookHandler,
		healthCheckHandler,
	}
	return
}

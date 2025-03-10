package cmd

import (
	"short/dep"
	"short/dep/provider"

	"github.com/byliuyang/app/fw"
)

type GithubConfig struct {
	ClientID     string
	ClientSecret string
}

func start(
	dbConfig fw.DBConfig,
	migrationRoot string,
	wwwRoot string,
	recaptchaSecret string,
	githubConfig GithubConfig,
	jwtSecret string,
	dbConnector fw.DBConnector,
	dbMigrationTool fw.DBMigrationTool,
) {
	db, err := dbConnector.Connect(dbConfig)
	if err != nil {
		panic(err)
	}

	err = dbMigrationTool.Migrate(db, migrationRoot)
	if err != nil {
		panic(err)
	}

	graphqlAPI := dep.InjectGraphQlService(
		"GraphQL API",
		db,
		"/graphql",
		provider.ReCaptchaSecret(recaptchaSecret),
		provider.JwtSecret(jwtSecret),
	)
	graphqlAPI.Start(8080)

	httpAPI := dep.InjectRoutingService(
		"Routing API",
		db,
		provider.WwwRoot(wwwRoot),
		provider.GithubClientID(githubConfig.ClientID),
		provider.GithubClientSecret(githubConfig.ClientSecret),
		provider.JwtSecret(jwtSecret),
	)
	httpAPI.StartAndWait(80)
}

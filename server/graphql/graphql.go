package graphql

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/blinfoldking/blockchain-go-node/model"
	"github.com/blinfoldking/blockchain-go-node/resolver"
	"github.com/graph-gophers/graphql-go"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type GraphQLHandler struct{}

// GET server graphql playground
func (handler *GraphQLHandler) Playground(c echo.Context) error {
	c.HTML(http.StatusOK, page)
	return nil
}

// POST graphql query
func (handler *GraphQLHandler) Query(c echo.Context) error {
	s, err := getSchema("schema/schema.gql")
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	err = GqlResponse(c, graphql.MustParseSchema(s, resolver.New()))
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return nil
}

func GqlResponse(ctx echo.Context, Schema *graphql.Schema) error {

	r := ctx.Request()

	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}

	var graphqlContext context.Context
	authorization := ctx.Request().Header.Get("Authorization")
	if authorization != "" {
		user, err := model.ValidateToken(authorization)
		if err == nil {
			graphqlContext = context.WithValue(context.Background(), "user", user)
		} else {
			graphqlContext = r.Context()
		}
	} else {
		graphqlContext = r.Context()
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return err
	}

	response := Schema.Exec(graphqlContext, params.Query, params.OperationName, params.Variables)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return err
	}

	ctx.Response().Write(responseJSON)

	return nil
}

func getSchema(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

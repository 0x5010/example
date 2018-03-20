package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var data map[string]user = map[string]user{
	"1": user{
		ID:   "1",
		Name: "Dan",
	},
	"2": user{
		ID:   "2",
		Name: "Lee",
	},
	"3": user{
		ID:   "3",
		Name: "Nick",
	},
}

var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					idQuery, isOK := p.Args["id"].(string)
					if isOK {
						if v, ok := data[idQuery]; ok {
							return v, nil
						}
					}
					return nil, nil
				},
			},
		},
	},
)

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)

func executeQuery(query string, schema graphql.Schema, vars map[string]interface{}) *graphql.Result {
	res := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  query,
		VariableValues: vars,
	})
	if len(res.Errors) > 0 {
		fmt.Printf("error: %v", res.Errors)
	}
	return res
}

func ginHandler() gin.HandlerFunc {
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	router := gin.Default()
	router.Any("/graphql", ginHandler())
	router.Run(":8080")
}

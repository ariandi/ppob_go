package api

import (
	db "github.com/ariandi/ppob_go/db/sqlc"
	"github.com/ariandi/ppob_go/token"
	"github.com/ariandi/ppob_go/util"
	"github.com/gin-gonic/gin"
	"testing"
)

func TestServer_softDeleteCategory(t *testing.T) {
	type fields struct {
		store      db.Store
		TokenMaker token.Maker
		Router     *gin.Engine
		config     util.Config
	}
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &Server{
				store:      tt.fields.store,
				TokenMaker: tt.fields.TokenMaker,
				Router:     tt.fields.Router,
				config:     tt.fields.config,
			}
			server.softDeleteCategory(tt.args.ctx)
		})
	}
}

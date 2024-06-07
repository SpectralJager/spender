package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SpectralJager/spender/db"
	"github.com/SpectralJager/spender/types"
	"github.com/SpectralJager/spender/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testDBURI    = "mongodb://localhost:27017"
	testDBNAME   = "spender_test"
	testUSERCOLL = "users_test"
)

type testDB struct {
	db.UserStore
}

func (tdb testDB) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testDB {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(testDBURI))
	if err != nil {
		t.Fatal(err)
	}
	// defer client.Disconnect(ctx)

	return &testDB{
		UserStore: db.NewMongoUserStore(client, testDBNAME, testUSERCOLL),
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	uh := NewUserHandler(tdb.UserStore)
	app := echo.New()
	testCases := []struct {
		param types.CreateUserParams
		res   bool
	}{
		{
			param: types.CreateUserParams{
				FirstName: "test",
				LastName:  "test",
				Email:     "some.test@test.com",
				Password:  "testpassword",
			},
			res: true,
		},
		{
			param: types.CreateUserParams{
				FirstName: "Madelynn",
				LastName:  "Lind",
				Email:     "hadley72@gmail.com",
				Password:  "g0qCPpf4xXyUNII",
			},
			res: true,
		},
		{
			param: types.CreateUserParams{
				FirstName: "Jaeden",
				LastName:  "Schulist",
				Email:     "blair_kuhlman@hotmail.com",
				Password:  "G9TmLORMHSMa3ga",
			},
			res: true,
		},
		{
			param: types.CreateUserParams{
				FirstName: "Mervin",
				LastName:  "Kling",
				Email:     "freddie_braun77@gmail.com",
				Password:  "Vlb7dBAZbBGmGn_",
			},
			res: true,
		},
		{
			param: types.CreateUserParams{
				FirstName: "Nannie",
				LastName:  "Price",
				Email:     "queen_wuckert34@gmail.com",
				Password:  "ynPb3542ZLliBrE",
			},
			res: true,
		},
		{
			param: types.CreateUserParams{
				FirstName: "Lucile",
				LastName:  "Rippin",
				Email:     "armando_nienow85@hotmail.com",
				Password:  "1eZ1tL1WQx_coL4",
			},
			res: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.param.Email, func(t *testing.T) {
			data, _ := json.Marshal(tC.param)

			req := httptest.NewRequest("", "/", bytes.NewReader(data))
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			ctx := app.NewContext(req, rec)
			uh.PostUser(ctx)
			res, err := utils.DecodeBody[map[string]any](rec.Body)
			if err != nil {
				t.Fatal(err)
			}
			if tC.res && rec.Code == http.StatusInternalServerError {
				t.Fatalf("%s", res)
			} else if !tC.res && rec.Code == http.StatusOK {
				t.Fatalf("expected error, got: %s", res)
			}
		})
	}
}

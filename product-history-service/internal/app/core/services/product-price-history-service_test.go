package services

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"

	mongorepo "github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/adapters/secondary/mongodb"
	"github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/infrastructure"
	"github.com/juanpabloavilan/meli-interview-exercise/product-history-service/internal/app/ports"
	"github.com/juanpabloavilan/meli-interview-exercise/product-history-service/mocks"
)

// TestIntegrationImportFromCSVFile tests the integration of loading prices in batches to the DB
func TestIntegrationImportFromCSVFile(t *testing.T) {
	ctx := context.Background()

	mongoDBContainer, err := mongodb.Run(ctx, "mongodb/mongodb-community-server")
	defer func() {
		if err := testcontainers.TerminateContainer(mongoDBContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	assert.NoError(t, err)

	connString, err := mongoDBContainer.ConnectionString(ctx)
	assert.NoError(t, err)

	db, err := infrastructure.NewMongoDB(ctx, connString, "test")
	assert.NoError(t, err)

	type fields struct {
		repo ports.ProductHistoryRepo
	}
	type args struct {
		ctx      context.Context
		filePath string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  bool
		wantRows int64
	}{
		{
			name: "should create 1000 documents from file",
			fields: fields{
				mongorepo.NewProductPriceHistoryRepo(db),
			},
			args: args{
				ctx,
				"./testdata/test-product-price-history-10-rows.csv",
			},
			wantRows: 10,
		},
		{
			name: "should create 175000 documents from file",
			fields: fields{
				mongorepo.NewProductPriceHistoryRepo(db),
			},
			args: args{
				ctx,
				"./testdata/test-product-price-history-175000-rows.csv",
			},
			wantRows: 175914,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &productPriceHistoryService{
				repo: tt.fields.repo,
			}
			t.Cleanup(func() {
				log.Printf("CLEANING UP")
				_, err := db.Collection(mongorepo.ProducPriceHistoryCollection).DeleteMany(ctx, bson.D{})
				assert.NoError(t, err)
			})

			file, err := os.Open(tt.args.filePath)
			assert.NoError(t, err)
			defer file.Close()

			if err := s.ImportFromCSVFile(tt.args.ctx, file); (err != nil) != tt.wantErr {
				t.Errorf("productPriceHistoryService.ImportFromCSVFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			totalRows, err := db.Collection(mongorepo.ProducPriceHistoryCollection).CountDocuments(ctx, bson.D{}, nil)
			assert.NoError(t, err)

			assert.Equal(t, totalRows, tt.wantRows)

		})
	}
}

func BenchmarkIntegrationImportFromCSVFile(b *testing.B) {
	ctx := context.Background()

	repo := mocks.NewProductHistoryRepo(b)
	repo.EXPECT().AddMany(ctx, mock.Anything).Return(nil).After(10 * time.Millisecond)

	s := &productPriceHistoryService{repo, nil}

	for i := 0; i < b.N; i++ {
		file, err := os.Open("./testdata/test-product-price-history-10-rows.csv")
		assert.NoError(b, err)
		err = s.ImportFromCSVFile(ctx, file)
		assert.NoError(b, err)
		file.Close()
	}
}

package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	groupitementity "github.com/willjrcom/sales-backend-go/internal/domain/group_item"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	processentity "github.com/willjrcom/sales-backend-go/internal/domain/process"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	schemaentity "github.com/willjrcom/sales-backend-go/internal/domain/schema"
	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
)

var (
	username = "admin"
	password = "admin"
	host     = "localhost"
	port     = "5432"
	dbName   = "sales-db"
)

func NewPostgreSQLConnection(ctx context.Context) (*bun.DB, error) {
	// Prepare connection string parameterized
	connectionParams := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username,
		password,
		host,
		port,
		dbName,
	)

	// Connect to database doing a PING
	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connectionParams), pgdriver.WithTimeout(time.Second*30)))

	// Verifique se o banco de dados já existe.
	if err := db.Ping(); err != nil {
		log.Printf("erro ao conectar ao banco de dados: %v", err)
	}

	// set connection settings
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Duration(60) * time.Minute)

	dbBun := bun.NewDB(db, pgdialect.New())

	if err := LoadAllSchemas(ctx, dbBun); err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, schemaentity.Schema("schema"), schemaentity.LOST_SCHEMA)
	if err := CreateSchema(ctx, dbBun); err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, schemaentity.Schema("schema"), schemaentity.DEFAULT_SCHEMA)
	if err := CreateSchema(ctx, dbBun); err != nil {
		return nil, err
	}

	if err := defaultTables(ctx, dbBun); err != nil {
		return nil, err
	}

	if err := LoadAllSchemas(ctx, dbBun); err != nil {
		return nil, err
	}

	fmt.Println("Db connected")
	return dbBun, nil
}

func ChangeSchema(ctx context.Context, db *bun.DB) error {
	schemaName, err := GetSchema(ctx)

	if err != nil {
		return err
	}

	_, err = db.Exec("SET search_path=?", schemaName)
	return err
}

func ChangeToPublicSchema(ctx context.Context, db *bun.DB) error {
	_, err := db.Exec("SET search_path=?", schemaentity.DEFAULT_SCHEMA)
	return err
}

func GetSchema(ctx context.Context) (string, error) {
	schemaName := ctx.Value(schemaentity.Schema("schema"))
	if schemaName == nil {
		return "", errors.New("schema not found")
	}
	return schemaName.(string), nil
}

func CreateSchema(ctx context.Context, db *bun.DB) error {
	schemaName, err := GetSchema(ctx)

	if err != nil {
		schemaName = schemaentity.LOST_SCHEMA
	}

	if _, err := db.Exec("CREATE SCHEMA IF NOT EXISTS " + schemaName); err != nil {
		return err
	}

	return nil
}

func defaultTables(ctx context.Context, db *bun.DB) error {
	db.RegisterModel((*companyentity.User)(nil))
	db.RegisterModel((*companyentity.CompanyToUsers)(nil))
	db.RegisterModel((*companyentity.CompanyWithUsers)(nil))

	if err := RegisterModels(ctx, db); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*companyentity.User)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*companyentity.CompanyToUsers)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*companyentity.CompanyWithUsers)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, "CREATE EXTENSION IF NOT EXISTS pgcrypto;"); err != nil {
		return err
	}

	return nil
}

func LoadAllSchemas(ctx context.Context, db *bun.DB) error {
	results, err := db.QueryContext(ctx, "SELECT schema_name FROM information_schema.schemata;")

	if err != nil {
		return err
	}

	for results.Next() {
		var schemaName string
		if err := results.Scan(&schemaName); err != nil {
			return err
		}

		if !strings.Contains(schemaName, "loja_") {
			continue
		}

		ctx = context.WithValue(ctx, schemaentity.Schema("schema"), schemaName)

		if err := RegisterModels(ctx, db); err != nil {
			return err
		}

		if err := LoadCompanyModels(ctx, db); err != nil {
			return err
		}

		if _, err := db.QueryContext(ctx, "CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_contact ON contacts (ddd, number);"); err != nil {
			return err
		}

		if err := SetupContactFtSearch(ctx, db); err != nil {
			return err
		}
	}

	return nil
}

func RegisterModels(ctx context.Context, db *bun.DB) error {
	mu := sync.Mutex{}

	mu.Lock()
	if err := CreateSchema(ctx, db); err != nil {
		mu.Unlock()
		return err
	}

	if err := ChangeSchema(ctx, db); err != nil {
		mu.Unlock()
		return err
	}

	db.RegisterModel((*entity.Entity)(nil))

	db.RegisterModel((*productentity.CategoryToAdditional)(nil))
	db.RegisterModel((*productentity.Size)(nil))
	db.RegisterModel((*productentity.Quantity)(nil))
	db.RegisterModel((*productentity.Category)(nil))
	db.RegisterModel((*productentity.ProcessRule)(nil))
	db.RegisterModel((*productentity.Product)(nil))

	db.RegisterModel((*addressentity.Address)(nil))
	db.RegisterModel((*personentity.Contact)(nil))
	db.RegisterModel((*cliententity.Client)(nil))
	db.RegisterModel((*employeeentity.Employee)(nil))

	db.RegisterModel((*processentity.Process)(nil))
	db.RegisterModel((*itementity.ItemToAdditional)(nil))
	db.RegisterModel((*itementity.Item)(nil))
	db.RegisterModel((*groupitementity.GroupItem)(nil))

	db.RegisterModel((*orderentity.PickupOrder)(nil))
	db.RegisterModel((*orderentity.DeliveryOrder)(nil))
	db.RegisterModel((*orderentity.TableOrder)(nil))
	db.RegisterModel((*orderentity.PaymentOrder)(nil))
	db.RegisterModel((*orderentity.Order)(nil))

	db.RegisterModel((*tableentity.Table)(nil))
	db.RegisterModel((*shiftentity.Shift)(nil))
	db.RegisterModel((*companyentity.Company)(nil))

	return nil
}

func LoadCompanyModels(ctx context.Context, db *bun.DB) error {
	mu := sync.Mutex{}

	mu.Lock()
	if err := CreateSchema(ctx, db); err != nil {
		mu.Unlock()
		return err
	}

	if err := ChangeSchema(ctx, db); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*entity.Entity)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*productentity.Size)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*productentity.CategoryToAdditional)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*productentity.Category)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*productentity.Quantity)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*productentity.ProcessRule)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*productentity.Product)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*addressentity.Address)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*personentity.Contact)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*personentity.Person)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*cliententity.Client)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*employeeentity.Employee)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*processentity.Process)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*itementity.Item)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*itementity.ItemToAdditional)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*groupitementity.GroupItem)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*orderentity.PickupOrder)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*orderentity.DeliveryOrder)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*orderentity.TableOrder)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*orderentity.PaymentOrder)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*orderentity.Order)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*tableentity.Table)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*shiftentity.Shift)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*companyentity.Company)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	return nil
}

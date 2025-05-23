package repositories_test

import (
	"context"
	"errors"
	"testing"
	"user-service/domain/dto"
	repositories "user-service/repositories/user"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUserRepository_Register(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Creates PostgreSQL Mock
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		// Create DB instances with the mock
		dialector := postgres.New(postgres.Config{
			Conn:       sqlDB,
			DriverName: "postgres",
		})
		db, err := gorm.Open(dialector, &gorm.Config{})
		require.NoError(t, err)

		// Create Repository
		repo := repositories.NewUserRepository(db)

		req := &dto.RegisterRequest{
			Name:            "faisal",
			Username:        "faisalabu",
			Email:           "faisal@mail.com",
			Password:        "strongpassword",
			ConfirmPassword: "strongpassword",
			PhoneNumber:     "082313113",
			RoleID:          1,
		}

		// Mock SQL Query Execution
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "users".*RETURNING .*`).
			WillReturnRows(sqlmock.NewRows([]string{
				"uuid", "username", "email", "name", "password", "phone_number", "role_id",
			}).AddRow(
				uuid.New(), req.Username, req.Email, req.Name, req.Password, req.PhoneNumber, req.RoleID,
			))
		mock.ExpectCommit()

		response, err := repo.Register(context.Background(), req)
		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, req.Username, response.Username)
		assert.Equal(t, req.Email, response.Email)
		assert.Equal(t, req.Name, response.Name)
		assert.Equal(t, req.Password, response.Password)
		assert.Equal(t, req.PhoneNumber, response.PhoneNumber)
		assert.Equal(t, req.RoleID, response.RoleID)
		assert.NotEqual(t, uuid.Nil, response.UUID)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	})

	t.Run("failed", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		dialector := postgres.New(postgres.Config{
			Conn:       sqlDB,
			DriverName: "postgres",
		})

		db, err := gorm.Open(dialector, &gorm.Config{})
		require.NoError(t, err)

		repo := repositories.NewUserRepository(db)

		req := &dto.RegisterRequest{
			Name:            "faisal",
			Username:        "faisalabu",
			Email:           "faisal@mail.com",
			Password:        "strongpassword",
			ConfirmPassword: "strongpassword",
			PhoneNumber:     "082313113",
			RoleID:          1,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "users"`).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		response, err := repo.Register(context.Background(), req)
		// Assertions
		assert.Error(t, err)
		assert.Nil(t, response)

		// Verify that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}

	})
}

func TestUserRepository_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		dialector := postgres.New(postgres.Config{
			Conn:       sqlDB,
			DriverName: "postgres",
		})

		db, err := gorm.Open(dialector, &gorm.Config{})
		require.NoError(t, err)

		repo := repositories.NewUserRepository(db)

		password := "newpassword123"
		req := &dto.UpdateRequest{
			Username:        "faisalupdate",
			Name:            "faisalabuupdate",
			Email:           "faisalupdate@mail.com",
			Password:        &password,
			ConfirmPassword: &password,
			PhoneNumber:     "0928318239",
		}

		uuid := uuid.New()
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users" SET "name"=\$1,"username"=\$2,"password"=\$3,"phone_number"=\$4,"updated_at"=\$5 WHERE uuid = \$6`).
			WithArgs(
				req.Name,
				req.Username,
				*req.Password,
				req.PhoneNumber,
				sqlmock.AnyArg(),
				uuid,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		response, err := repo.Update(context.Background(), req, uuid.String())

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, req.Name, response.Name)
		assert.Equal(t, req.Username, response.Username)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}

	})

	t.Run("failed", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		dialector := postgres.New(postgres.Config{
			Conn:       sqlDB,
			DriverName: "postgres",
		})

		db, err := gorm.Open(dialector, &gorm.Config{})
		require.NoError(t, err)

		repo := repositories.NewUserRepository(db)
		password := "newpassword123"
		req := &dto.UpdateRequest{
			Username:        "faisalupdate",
			Name:            "faisalabuupdate",
			Email:           "faisalupdate@mail.com",
			Password:        &password,
			ConfirmPassword: &password,
			PhoneNumber:     "0928318239",
		}

		uuid := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users" SET "name"=\$1,"username"=\$2,"password"=\$3,"phone_number"=\$4,"updated_at"=\$5 WHERE uuid = \$6`).
			WithArgs(
				req.Name,
				req.Username,
				*req.Password,
				req.PhoneNumber,
				sqlmock.AnyArg(),
				uuid,
			).WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		response, err := repo.Update(context.Background(), req, uuid.String())
		require.Error(t, err)
		assert.Nil(t, response)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}

	})
}

func TestUserRepository_FindByEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		dialector := postgres.New(postgres.Config{
			Conn:       sqlDB,
			DriverName: "postgres",
		})
		db, err := gorm.Open(dialector, &gorm.Config{})
		require.NoError(t, err)

		repo := repositories.NewUserRepository(db)

		email := "faisal@mail.com"

		mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 ORDER BY "users"."id" LIMIT \$2`).
			WithArgs(email, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).AddRow(1, "faisal", "faisal@mail.com"))

		response, err := repo.FindByEmail(context.Background(), email)
		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, email, response.Email)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	})

	t.Run("failed", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		dialector := postgres.New(postgres.Config{
			Conn:       sqlDB,
			DriverName: "postgres",
		})
		db, err := gorm.Open(dialector, &gorm.Config{})
		require.NoError(t, err)

		repo := repositories.NewUserRepository(db)

		email := "faisal@mail.com"

		mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 ORDER BY "users"."id" LIMIT \$2`).
			WithArgs(email, 1).
			WillReturnError(errors.New("database error"))

		response, err := repo.FindByEmail(context.Background(), email)
		require.Error(t, err)
		assert.Nil(t, response)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	})

	t.Run("data not found", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		dialector := postgres.New(postgres.Config{
			Conn:       sqlDB,
			DriverName: "postgres",
		})
		db, err := gorm.Open(dialector, &gorm.Config{})
		require.NoError(t, err)

		repo := repositories.NewUserRepository(db)

		email := "faisal@mail.com"

		mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 ORDER BY "users"."id" LIMIT \$2`).
			WithArgs(email, 1).
			WillReturnError(errors.New("user not found"))

		response, err := repo.FindByEmail(context.Background(), email)
		require.Error(t, err)
		assert.Nil(t, response)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	})
}

func TestUserRepository_FindByUsername(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		dialector := postgres.New(postgres.Config{
			Conn:       sqlDB,
			DriverName: "postgres",
		})
		db, err := gorm.Open(dialector, &gorm.Config{})
		require.NoError(t, err)

		repo := repositories.NewUserRepository(db)

		username := "faisalabu"

		mock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1 ORDER BY "users"."id" LIMIT \$2`).
			WithArgs(username, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "username"}).AddRow(1, "faisal", "faisalabu"))

		response, err := repo.FindByUsername(context.Background(), username)
		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, username, response.Username)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	})

	t.Run("failed", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		dialector := postgres.New(postgres.Config{
			Conn:       sqlDB,
			DriverName: "postgres",
		})
		db, err := gorm.Open(dialector, &gorm.Config{})
		require.NoError(t, err)

		repo := repositories.NewUserRepository(db)

		username := "faisalabu"

		mock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1 ORDER BY "users"."id" LIMIT \$2`).
			WithArgs(username, 1).
			WillReturnError(errors.New("database error"))

		response, err := repo.FindByUsername(context.Background(), username)
		require.Error(t, err)
		assert.Nil(t, response)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	})

	t.Run("data not found", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		dialector := postgres.New(postgres.Config{
			Conn:       sqlDB,
			DriverName: "postgres",
		})
		db, err := gorm.Open(dialector, &gorm.Config{})
		require.NoError(t, err)

		repo := repositories.NewUserRepository(db)

		username := "faisalabu"

		mock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1 ORDER BY "users"."id" LIMIT \$2`).
			WithArgs(username, 1).
			WillReturnError(errors.New("user not found"))

		response, err := repo.FindByUsername(context.Background(), username)
		require.Error(t, err)
		assert.Nil(t, response)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	})
}

func TestUserRepository_FindByUUID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		dialector := postgres.New(postgres.Config{
			Conn:       sqlDB,
			DriverName: "postgres",
		})
		db, err := gorm.Open(dialector, &gorm.Config{})
		require.NoError(t, err)

		repo := repositories.NewUserRepository(db)

		uuid := uuid.New()

		mock.ExpectQuery(`SELECT \* FROM "users" WHERE uuid = \$1 ORDER BY "users"."id" LIMIT \$2`).
			WithArgs(uuid, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "uuid"}).AddRow(1, "faisal", uuid.String()))

		response, err := repo.FindByUUID(context.Background(), uuid.String())
		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, uuid, response.UUID)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	})

	t.Run("failed", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		dialector := postgres.New(postgres.Config{
			Conn:       sqlDB,
			DriverName: "postgres",
		})
		db, err := gorm.Open(dialector, &gorm.Config{})
		require.NoError(t, err)

		repo := repositories.NewUserRepository(db)

		uuid := uuid.New()

		mock.ExpectQuery(`SELECT \* FROM "users" WHERE uuid = \$1 ORDER BY "users"."id" LIMIT \$2`).
			WithArgs(uuid, 1).
			WillReturnError(errors.New("database error"))

		response, err := repo.FindByUUID(context.Background(), uuid.String())
		require.Error(t, err)
		assert.Nil(t, response)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	})

	// test

	t.Run("data not found", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		dialector := postgres.New(postgres.Config{
			Conn:       sqlDB,
			DriverName: "postgres",
		})
		db, err := gorm.Open(dialector, &gorm.Config{})
		require.NoError(t, err)

		repo := repositories.NewUserRepository(db)

		uuid := uuid.New()

		mock.ExpectQuery(`SELECT \* FROM "users" WHERE uuid = \$1 ORDER BY "users"."id" LIMIT \$2`).
			WithArgs(uuid, 1).
			WillReturnError(errors.New("user not found"))

		response, err := repo.FindByUUID(context.Background(), uuid.String())
		require.Error(t, err)
		assert.Nil(t, response)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	})
}

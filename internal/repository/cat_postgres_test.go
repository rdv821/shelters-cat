package repository

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/catService/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var (
	repository SheltersCatRepository
)

var cat = &model.Cat{
	ID:         uuid.New(),
	Name:       "Cat 1",
	Age:        2,
	Vaccinated: true,
}

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		logrus.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14",
		Env: []string{
			"POSTGRES_PASSWORD=qwerty",
			"POSTGRES_USER=postgres",
			"POSTGRES_DB=catsdb",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		logrus.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseURL := fmt.Sprintf("postgres://postgres:qwerty@%s/catsdb?sslmode=disable", hostAndPort)

	logrus.Info("Connecting to database on url: ", databaseURL)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds
	if err = pool.Retry(func() error {
		url := fmt.Sprintf(`-url=jdbc:postgresql://%s/catsdb`, hostAndPort)
		cmd := exec.Command("flyway", url, "-user=postgres", "-password=qwerty", "-locations=filesystem:../../migrations", "migrate")
		logrus.Info(cmd)
		flywayError := cmd.Run()
		if flywayError != nil {
			logrus.Fatalf("Flyway error: %v", flywayError.Error())
			return flywayError
		}

		poolPgx, _ := pgxpool.Connect(context.Background(), databaseURL)
		repository = NewPostgresRepository(poolPgx)
		return nil
	}); err != nil {
		logrus.Fatalf("Could not connect to docker: %s", err.Error())
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		logrus.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestCreate(t *testing.T) {
	err := repository.Create(context.Background(), cat)
	require.NoError(t, err)
}

func TestGet(t *testing.T) {
	repository.Create(context.Background(), cat)
	testCat, err := repository.Get(context.Background(), cat.ID)
	require.NotEmpty(t, testCat)
	require.Equal(t, cat.ID, testCat.ID)
	require.Equal(t, cat.Name, testCat.Name)
	require.Equal(t, cat.Age, testCat.Age)
	require.Equal(t, cat.Vaccinated, testCat.Vaccinated)
	require.NoError(t, err, cat)
}

func TestGetNonExistingCat(t *testing.T) {
	var cats = &model.Cat{
		Name: "Cat 1",
	}

	_, err := repository.Get(context.Background(), cats.ID)
	require.Error(t, err)
}

func TestDelete(t *testing.T) {
	repository.Create(context.Background(), cat)
	err := repository.Delete(context.Background(), cat.ID)
	require.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	repository.Create(context.Background(), cat)
	err := repository.Update(context.Background(), cat)
	require.NoError(t, err)
}

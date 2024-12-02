package migrations

import (
	"github.com/pkg/errors"
	"os"
	"os/exec"
)

func RunGooseMigrations(dbname string) error {
	// Создаем команду для запуска goose
	cmd := exec.Command(
		"goose",
		"-dir", "./migrations",
		"up",
	)
	// Устанавливаем переменные окружения
	cmd.Env = append(os.Environ(),
		"GOOSE_DRIVER=postgres", "GOOSE_DBSTRING=host=localhost user=postgres password=postgres database=medodsTest sslmode=dosable")

	// Выполняем команду и захватываем вывод
	_, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "Error running goose migrations")
	}

	return nil
}

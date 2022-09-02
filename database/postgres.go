package database

import (
	"context"
	"database/sql"
	"log"
	"w00k/go/rest-ws/models"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository: conexión a la base de datos
func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db}, nil
}

// InsertUser: inserción de un usuario a la base de datos
// los casos que soporta son:
// - inserta el user, retorna nil
// - error al insertar el usuario, retorna el error
func (repo *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)", user.Id, user.Email, user.Password)
	return err
}

// GetUserById: obtiene el ususario por el id el usuario,
// los casos que soporta son:
// - obtiene el usuario, lo retorna
// - en caso de error, retorna un objeto usuario vacio y el error
// - en caso de no encontrar el usuario, retorna un objeto usuario vacio y el error en nil
func (repo *PostgresRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, email FROM users WHERE id = $1", id)

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var user = models.User{}
	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Email); err == nil {
			return &user, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail: obtiene el ususario por el id el usuario,
// los casos que soporta son:
// - obtiene el usuario, lo retorna
// - en caso de error, retorna un objeto usuario vacio y el error
// - en caso de no encontrar el usuario, retorna un objeto usuario vacio y el error en nil
func (repo *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, email, password FROM users WHERE email = $1", email)

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var user = models.User{}
	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Email, &user.Password); err == nil {
			return &user, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}

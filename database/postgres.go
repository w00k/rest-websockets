package database

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strconv"
	"w00k/go/rest-ws/models"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var configMap = make(map[string]int)

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

// InsertPost: inserción de un post a la base de datos
// los casos que soporta son:
// - inserta el post, retorna nil
// - error al insertar el post, retorna el error
func (repo *PostgresRepository) InsertPost(ctx context.Context, post *models.Post) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO posts (id, post_content, user_id) VALUES ($1, $2, $3)", post.Id, post.PostContent, post.UserId)
	return err
}

func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}

// GetPostById: obtiene el post por el id el post,
// los casos que soporta son:
// - obtiene el post, lo retorna
// - en caso de error, retorna un objeto post vacio y el error
// - en caso de no encontrar el post, retorna un objeto post vacio y el error en nil
func (repo *PostgresRepository) GetPostById(ctx context.Context, id string) (*models.Post, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, post_content, created_at, user_id FROM posts WHERE id = $1", id)

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var post = models.Post{}
	for rows.Next() {
		if err = rows.Scan(&post.Id, &post.PostContent, &post.CreateAt, post.UserId); err == nil {
			return &post, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &post, nil
}

// UpdatePost: update de un post a la base de datos
// los casos que soporta son:
// - actualiza el post, retorna nil
// - error al actualizar el post, retorna el error
func (repo *PostgresRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	_, err := repo.db.ExecContext(ctx, "UPDATE posts SET post_content = $2 WHERE id = $1 and user_id = $3", post.Id, post.PostContent, post.UserId)
	return err
}

// DeletePost: borra un post a la base de datos
// los casos que soporta son:
// - borra el post, retorna nil
// - error al borrar el post, retorna el error
func (repo *PostgresRepository) DeletePost(ctx context.Context, id string, userId string) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM posts WHERE id = $1 and user_id = $2", id, userId)
	return err
}

func (repo *PostgresRepository) ListPost(ctx context.Context, page uint64) ([]*models.Post, error) {
	pageSize := getPageSize()
	rows, err := repo.db.QueryContext(ctx, "SELECT id, post_content, user_id, created_at FROM posts LIMIT $1 OFFSET $2", pageSize, page*uint64(pageSize))
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var posts []*models.Post
	for rows.Next() {
		var post = models.Post{}
		if err = rows.Scan(&post.Id, &post.PostContent, &post.UserId, &post.CreateAt); err == nil {
			posts = append(posts, &post)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func getPageSize() int {
	if value, ok := configMap["page"]; ok {
		return value
	}
	page := 2
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file, value PAGE\n")
		return page
	}
	PAGE := os.Getenv("PAGE")
	page, err = strconv.Atoi(PAGE)
	if err != nil {
		log.Printf("Error with PAGE value %s, must be number\n", PAGE)
		page = 2
	}
	configMap["page"] = page
	return page
}

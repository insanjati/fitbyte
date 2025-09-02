package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/insanjati/fitbyte/internal/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll() ([]model.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users ORDER BY id`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) RegisterNewUser(c *gin.Context, payload model.User) (model.User, error){
	newId := uuid.Must(uuid.NewV7())
	var user model.User

	//Insert user
	query := `INSERT INTO Users(id, name, email, password) VALUES($1, $2, $3, $4) RETURNING id, name, email`

	err := r.db.QueryRowContext(c, query, newId, payload.Name, payload.Email, payload.Password).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil{
		if c.Err() != nil{
			return model.User{}, c.Err()
		}
		return user, err
	}
 
	//return email and token
	return model.User{}, nil
}

// func isEmailValid(e string) bool{
// 	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
// 	return emailRegex.MatchString(e)
// }

func (r *UserRepository) GetUserByEmail(c *gin.Context, email string) (model.User, error){
	var user model.User

	query := `SELECT id, name, email, password FROM users WHERE email=$1`
	err := r.db.QueryRowContext(c, query).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil{
		if c.Err() != nil{
			return model.User{}, c.Err()
		}
		return model.User{}, err
	}
	return user, nil
}
package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 3 // 数据库操作超时时间为 3 秒

var db *sql.DB // 全局变量，用于存储数据库连接池

// New 是用于创建 data 包实例的函数。它返回 Models 类型，该类型包含了我们想要在整个应用程序中使用的各种类型。
func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: User{}, // 初始化 Models 结构体中的 User 字段
	}
}

// Models 是 data 包的类型。请注意，任何在此类型中作为成员的模型都可以在整个应用程序中使用，只要使用 app 变量，同时也需要在 New 函数中添加相应的模型。
type Models struct {
	User User // 用户模型
}

// User 是从数据库中获取的用户结构体。
type User struct {
	ID        int       `json:"id"`                   // 用户ID
	Email     string    `json:"email"`                // 邮箱
	FirstName string    `json:"first_name,omitempty"` // 名字
	LastName  string    `json:"last_name,omitempty"`  // 姓氏
	Password  string    `json:"-"`                    // 密码，使用 "-" 表示不在 JSON 中显示
	Active    int       `json:"active"`               // 活动状态
	CreatedAt time.Time `json:"created_at"`           // 创建时间
	UpdatedAt time.Time `json:"updated_at"`           // 更新时间
}

// GetAll 返回按姓氏排序的所有用户的切片。
func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout) // 创建上下文并设置超时时间
	defer cancel()                                                      // 延迟取消上下文

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
	from users order by last_name` // 查询语句，按姓氏排序

	rows, err := db.QueryContext(ctx, query) // 执行查询并获取结果集
	if err != nil {
		return nil, err
	}
	defer rows.Close() // 延迟关闭结果集

	var users []*User // 存储用户的切片

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		users = append(users, &user) // 将用户添加到切片中
	}

	return users, nil
}

// GetByEmail 根据邮箱返回一个用户。
func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where email = $1` // 根据邮箱查询用户

	var user User
	row := db.QueryRowContext(ctx, query, email) // 执行查询并返回一行结果

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetOne 根据 ID 返回一个用户。
func (u *User) GetOne(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where id = $1` // 根据 ID 查询用户

	var user User
	row := db.QueryRowContext(ctx, query, id) // 执行查询并返回一行结果

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update 根据接收者 u 中的信息更新数据库中的一个用户。
func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update users set
		email = $1,
		first_name = $2,
		last_name = $3,
		user_active = $4,
		updated_at = $5
		where id = $6
	` // 更新用户信息的 SQL 语句

	_, err := db.ExecContext(ctx, stmt,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Active,
		time.Now(),
		u.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// Delete 删除数据库中的一个用户，根据 User.ID。
func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1` // 根据 ID 删除用户

	_, err := db.ExecContext(ctx, stmt, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteByID 根据 ID 删除数据库中的一个用户。
func (u *User) DeleteByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1` // 根据 ID 删除用户

	_, err := db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// Insert 插入一个新用户到数据库，并返回新插入行的 ID。
func (u *User) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12) // 对密码进行加密
	if err != nil {
		return 0, err
	}

	var newID int
	stmt := `insert into users (email, first_name, last_name, password, user_active, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7) returning id` // 插入用户的 SQL 语句

	err = db.QueryRowContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		user.Active,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// ResetPassword 是用于更改用户密码的方法。
func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12) // 对密码进行加密
	if err != nil {
		return err
	}

	stmt := `update users set password = $1 where id = $2` // 更新用户密码的 SQL 语句
	_, err = db.ExecContext(ctx, stmt, hashedPassword, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// PasswordMatches 使用 Go 的 bcrypt 包比较用户提供的密码和数据库中存储的密码哈希值。如果密码匹配，返回 true；否则，返回 false。
func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText)) // 比较密码和哈希值
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// 密码不匹配
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

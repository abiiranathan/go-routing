// Code generated by "apigen"; DO NOT EDIT.

package svc

import (
	"gorm.io/gorm"
	"hello/models"
	"math"
)

/*
Options to pass to GORM.
Preferred order of option: Select,Group, Where, Or, Joins, Preload, Order, Limit, Offset, etc
For example:

		Select("column1, column2").
	    Where("column3 = ?", value).
	    Joins("JOIN t1 ON t1.id = t.id").
	    Order("column ASC").
	    Limit(10).
	    Offset(0).
	    Find(&results)
*/
type Option func(db *gorm.DB) *gorm.DB

/*
Where add conditions
See the docs for details on the various formats that where clauses can take. By default, where clauses chain with AND.

	Find the first user with name john
	Where("name = ?", "john")
*/
func Where(query string, args ...any) Option {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where(query, args...)
		return db
	}
}

/*
Preload preload associations with given conditions

	Preload("Orders", "state NOT IN (?)", "cancelled")
*/
func Preload(query string, args ...any) Option {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Preload(query, args...)
		return db
	}
}

/*
Select specify fields that you want when querying, creating, updating

Use Select when you only want a subset of the fields.
By default, GORM will select all fields.

	Select("name", "age")
*/
func Select(columns ...string) Option {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Select(columns)
		return db
	}
}

// Omit specify fields that you want to ignore when creating, updating and querying
func Omit(columns ...string) Option {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Omit(columns...)
		return db
	}
}

// Order("name DESC")
func Order(order string) Option {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Order(order)
		return db
	}
}

// Select("name, sum(age) as total"), Group("name")
func Group(group string) Option {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Group(group)
		return db
	}
}

// Select("name, sum(age) as total"), Group("name"), Having("name = ?", "john")
func Having(query interface{}, args ...interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Having(query, args...)
		return db
	}
}

// Limit(3) - Retrieve three users
func Limit(limit int) Option {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Limit(limit)
		return db
	}
}

// Distinct(args ...any) - Specify Distinct columns
func Distinct(args ...any) Option {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Distinct(args...)
		return db
	}
}

// Offset(2) - Select the third user
func Offset(offset int) Option {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Offset(offset)
		return db
	}
}

/*
Joins specify Joins conditions

	Joins("Account")
	Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "johndoe@example.org")
*/
func Joins(query string, args ...interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Joins(query, args...)
		return db
	}
}

func applyOptions(db *gorm.DB, options ...Option) *gorm.DB {
	for _, option := range options {
		db = option(db)
	}
	return db
}

// PaginatedResults defines options for paginated queries.
type PaginatedResults[T any] struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int64 `json:"total_pages"`
	Count      int64 `json:"count"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
	Results    []T   `json:"results"`
}

type postService interface {

	// Create new post
	Create(post *models.Post, options ...Option) error

	// Create multiple posts
	CreateMany(posts *[]models.Post, options ...Option) error

	// Update post with all the fields. Uses gorm.DB.Save()
	Update(postId int, post *models.Post, options ...Option) (models.Post, error)

	// Update a single column with specified conditions
	UpdateColumn(columnName string, value any, where string, target ...any) error

	// PartialUpdate for post. Only updates fields with non-zero values using gorm.DB.Updates(). Returns the updated post
	PartialUpdate(id int, post models.Post, options ...Option) (models.Post, error)

	// Permanently Delete post from the database by primary key
	Delete(id int) error

	// Permanently Delete post from the database matching conditions.
	DeleteWhere(value string, conds ...any) error

	// Get a single post by id (primary key)
	Get(id int, options ...Option) (models.Post, error)

	// Retrieves all posts from the database
	GetAll(options ...Option) (results []models.Post, err error)

	// Find the first record matching condition specified by query &amp; args
	FindOne(options ...Option) (models.Post, error)

	// Find all records matching condition specified by query &amp; args
	FindMany(options ...Option) (results []models.Post, err error)

	// GetPaginated retrieves a paginated list of posts
	GetPaginated(page int, pageSize int, options ...Option) (*PaginatedResults[models.Post], error)
}

// Implementation for postService interface
type postRepo struct {
	DB *gorm.DB
}

// Returns a post service that accesses the gorm.DB
// instance through dependancy injection
func newPostService(db *gorm.DB) postService {
	return &postRepo{DB: db}
}

// Create new post
func (repo *postRepo) CreateMany(posts *[]models.Post, options ...Option) error {
	if err := repo.DB.Omit().Create(posts).Error; err != nil {
		return err
	}

	// Refetch to load associations if any
	for i, record := range *posts {
		record, err := repo.Get(record.ID, options...)
		if err != nil {
			return err
		}
		(*posts)[i] = record
	}
	return nil
}

// Create new post
func (repo *postRepo) Create(post *models.Post, options ...Option) error {
	if err := repo.DB.Omit().Create(post).Error; err != nil {
		return err
	}

	// Refetch to load associations if any
	record, err := repo.Get(post.ID, options...)
	if err != nil {
		return err
	}
	*post = record
	return nil
}

// Update post with all the fields. Uses gorm.DB.Save()
func (repo *postRepo) Update(id int, post *models.Post, options ...Option) (models.Post, error) {
	// Make sure the ID is set on object to use Save(), otherwise you get unique constraint error.
	post.ID = id
	if err := repo.DB.Omit().Save(post).Error; err != nil {
		return models.Post{}, err
	}
	return repo.Get(id, options...)
}

// Update a single column. Gorm hooks will be fired because it uses Update() method.
func (repo *postRepo) UpdateColumn(columnName string, value any, where string, target ...any) error {
	return repo.DB.Model(&models.Post{}).Where(where, target...).Update(columnName, value).Error
}

// PartialUpdate for post. Only updates fields with no zero values. Returns the updated post
func (repo *postRepo) PartialUpdate(id int, post models.Post, options ...Option) (models.Post, error) {
	if err := repo.DB.Omit().Where("id=?", id).Model(&models.Post{}).Updates(post).Error; err != nil {
		return models.Post{}, err
	}

	var updatedpost models.Post
	db := repo.DB
	db = applyOptions(db, options...)
	if err := db.First(&updatedpost, id).Error; err != nil {
		return models.Post{}, err
	}
	return updatedpost, nil
}

// Permanently Delete post from the database by id
func (repo *postRepo) Delete(id int) error {
	if err := repo.DB.Unscoped().Delete(&models.Post{}, id).Error; err != nil {
		return err
	}
	return nil
}

// Permanently Delete post from the database matching conditions
func (repo *postRepo) DeleteWhere(value string, conds ...any) error {
	if err := repo.DB.Unscoped().Delete(&models.Post{}, value, conds).Error; err != nil {
		return err
	}
	return nil
}

// Make sure struct has an ID field as PK.
// Else skip this.
// Get a single post by id primary key
// Warning: Do not pass Where() option in options when using id, you will get unexpected results.
// (unless that's what you want!)
func (repo *postRepo) Get(id int, options ...Option) (models.Post, error) {
	var post models.Post
	db := repo.DB
	db = applyOptions(db, options...)
	if err := db.First(&post, id).Error; err != nil {
		return models.Post{}, err
	}
	return post, nil
}

// Retries all posts
func (repo *postRepo) GetAll(options ...Option) (results []models.Post, err error) {
	db := repo.DB
	db = applyOptions(db, options...)
	err = db.Find(&results).Error
	return
}

// GetPaginated retrieves a paginated list of users
func (repo *postRepo) GetPaginated(page int, pageSize int, options ...Option) (
	*PaginatedResults[models.Post], error) {

	var results []models.Post

	// Page must be >= 1
	if page < 1 {
		page = 1
	}

	// Calculate offset and limit
	offset := (page - 1) * pageSize

	// Query the database
	db := repo.DB
	db = applyOptions(db, options...)

	// Retrieve total count of records after applying options
	var totalCount int64
	if err := db.Model(&models.Post{}).Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if err := db.Offset(offset).Limit(pageSize).Find(&results).Error; err != nil {
		return nil, err
	}

	paginatedResults := &PaginatedResults[models.Post]{
		Page:       page,
		PageSize:   pageSize,
		HasNext:    int64(page*pageSize) < totalCount,
		HasPrev:    page > 1,
		Results:    results,
		Count:      totalCount,
		TotalPages: int64(math.Ceil(float64(totalCount) / float64(pageSize))),
	}

	return paginatedResults, nil
}

func (repo *postRepo) FindOne(options ...Option) (models.Post, error) {
	var post models.Post
	db := repo.DB
	db = applyOptions(db, options...)
	if err := db.First(&post).Error; err != nil {
		return models.Post{}, err
	}
	return post, nil
}

func (repo *postRepo) FindMany(options ...Option) (results []models.Post, err error) {
	db := repo.DB
	db = applyOptions(db, options...)
	err = db.Find(&results).Error
	return
}

// Service embeds all generated services
type Service struct {
	PostService postService
}

// Returns a new Service that embeds all the generated services.
func NewService(db *gorm.DB) *Service {
	return &Service{
		PostService: newPostService(db),
	}
}
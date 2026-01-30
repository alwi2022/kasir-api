package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll() ([]models.Product, error) {
	query := `
SELECT p.id, p.name, p.price, p.stock, p.category_id, c.name
FROM product p
LEFT JOIN category c ON c.id = p.category_id
ORDER BY p.id
`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		var cID sql.NullInt64
		var cName sql.NullString

		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &cID, &cName); err != nil {
			return nil, err
		}

		if cID.Valid {
			v := cID.Int64
			p.CategoryID = &v
		}
		if cName.Valid {
			v := cName.String
			p.CategoryName = &v
		}

		products = append(products, p)
	}

	return products, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO product (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	return err
}

// GetByID - ambil produk by ID
func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := `
SELECT p.id, p.name, p.price, p.stock, p.category_id, c.name
FROM product p
LEFT JOIN category c ON c.id = p.category_id
WHERE p.id = $1

`

	var p models.Product
	var cID sql.NullInt64
	var cName sql.NullString

	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &cID, &cName)
	if err == sql.ErrNoRows {
		return nil, errors.New("produk tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	if cID.Valid {
		v := cID.Int64
		p.CategoryID = &v
	}
	if cName.Valid {
		v := cName.String
		p.CategoryName = &v
	}

	return &p, nil
}

func (repo *ProductRepository) Update(product *models.Product) error {
	query := "UPDATE product SET name=$1, price=$2, stock=$3, category_id=$4 WHERE id=$5"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}
	return nil
}


func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM product  WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return err
}

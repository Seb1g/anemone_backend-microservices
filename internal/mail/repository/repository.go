package repository

import (
	"anemone_backend-microservices/internal/mail/model"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateAddress(address string, userID int64) (*model.TempAddress, error) {
	addr := &model.TempAddress{
		Address: address,
		UserID:  userID,
	}

	q := `
		INSERT INTO temp_addresses (address, user_id)
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	if err := r.db.QueryRow(q, addr.Address, addr.UserID).
		Scan(&addr.ID, &addr.CreatedAt); err != nil {
		return nil, err
	}

	return addr, nil
}

func (r *Repository) FindAddressByString(address string) (*model.TempAddress, error) {
	var addr model.TempAddress
	q := `SELECT id, address, user_id, created_at FROM temp_addresses WHERE address = $1`
	if err := r.db.Get(&addr, q, address); err != nil {
		return nil, err
	}
	return &addr, nil
}

func (r *Repository) SaveEmail(email *model.Email) error {
	q := `
		INSERT INTO emails (address_id, sender, recipients, subject, body, raw_data)
		VALUES ($1,$2,$3,$4,$5,$6)
	`
	_, err := r.db.Exec(
		q,
		email.AddressID,
		email.Sender,
		pq.Array(email.Recipients),
		email.Subject,
		email.Body,
		email.RawData,
	)
	return err
}

func (r *Repository) GetEmailsForAddress(addressID int) ([]model.Email, error) {
	var emails []model.Email
	q := `
		SELECT id, sender, recipients, subject, body, raw_data, received_at
		FROM emails
		WHERE address_id = $1
		ORDER BY received_at DESC
	`
	if err := r.db.Select(&emails, q, addressID); err != nil {
		return nil, err
	}
	return emails, nil
}

func (r *Repository) GetAddressesForUser(userID int64) ([]model.TempAddress, error) {
	var addresses []model.TempAddress
	q := `
		SELECT id, address, created_at
		FROM temp_addresses
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	if err := r.db.Select(&addresses, q, userID); err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *Repository) DeleteAddress(addressID int, userID int64) error {
	_, err := r.db.Exec(
		`DELETE FROM temp_addresses WHERE id=$1 AND user_id=$2`,
		addressID,
		userID,
	)
	return err
}

func (r *Repository) CheckAddressOwner(addressID int, userID int64) (bool, error) {
	var exists bool
	err := r.db.Get(
		&exists,
		`SELECT EXISTS(SELECT 1 FROM temp_addresses WHERE id=$1 AND user_id=$2)`,
		addressID,
		userID,
	)
	return exists, err
}

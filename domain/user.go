// Package domain includes domain definition files
package domain

import "time"

type User struct {
	ID        string    `bson:"_id" json:"id" db:"id"`
	Username  string    `bson:"username" json:"username" db:"username"`
	Password  string    `bson:"password" json:"password" db:"password"`
	Email     string    `bson:"email" json:"email" db:"email"`
	Bio       string    `bson:"bio" json:"bio" db:"bio"`
	Credit    int       `bson:"credit" json:"credit" db:"credit"`
	CreatedAt time.Time `bson:"created_at" json:"created_at" db:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at" db:"updated_at"`
}

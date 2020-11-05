// Package domain includes domain definition files
package domain

import "time"

type User struct {
	ID        string    `bson:"_id" json:"id"`
	Username  string    `bson:"username" json:"username"`
	Password  string    `bson:"password" json:"password"`
	Email     string    `bson:"email" json:"email"`
	Bio       string    `bson:"bio" json:"bio"`
	Credit    int       `bson:"credit" json:"credit"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

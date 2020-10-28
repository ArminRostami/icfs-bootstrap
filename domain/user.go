// Package domain includes domain definition files
package domain

type User struct {
	Username string `bson:"username" json:"username" binding:"required"`
	Password string `bson:"password" json:"password" binding:"required"`
	Credit   int    `bson:"credit"`
	ID       string `bson:"_id"`
}

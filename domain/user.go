// Package domain includes domain definition files
package domain

type User struct {
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
	Credit   int    `bson:"credit" json:"credit"`
	ID       string `bson:"_id"`
}

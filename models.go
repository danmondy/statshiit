package main

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

//User
type User struct{
	Id int64
	Email string
	Hashword string
	Rank int
	Since time.Time
}

func NewUser(email string, password string, rank int)User{
	u := User{Email: email, Rank: rank}
	u.SetHashword([]byte(password))
	u.Since = time.Now()
	return u
}

func (u *User) SetHashword(password []byte)error{
	hash, err := bcrypt.GenerateFromPassword(password, 0)//0 -> uses default cost (10 at time of writing, 4 min / 31 max)
	if err != nil {
		return err
	}else{ 
		u.Hashword = string(hash)
		return nil
	}
}

func (u *User) CompareHash(p []byte)error{
	return bcrypt.CompareHashAndPassword([]byte(u.Hashword), p)
}

//Character
type Character struct{
	Id int64
	UserId int64
	Name string
	StatsJson string
	Image string
	DateModified time.Time
}
func NewCharacter(userId int64, name string, stats string, image string)*Character{
	c := &Character{}
	c.Name = name
	c.UserId = userId
	c.StatsJson = stats
	c.Image = image
	return c
}

package user

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/jinzhu/gorm"
	"golite/internal/db"
	"golite/model/user"
)

var (
	U *userService
)

func init() {
	U = new()
}

type userService struct {
	o *gorm.DB
}

func new() *userService {
	return &userService{o: db.O}
}

func (us *userService) Login(u *user.User) error {
	username := u.Username
	ou, err := us.GetUserByUsername(username)
	if err != nil {
		return err
	}
	if ou == nil {
		return fmt.Errorf("user %s does not exist", username)
	}
	h := md5.New()
	h.Write([]byte(u.Password))
	u.Password = hex.EncodeToString(h.Sum(nil))
	if ou.Password != u.Password {
		return fmt.Errorf("user %s password is incorrect", username)
	}
	*u = *ou
	return nil
}

func (us *userService) Save(u *user.User) error {
	username := u.Username
	ou, err := us.GetUserByUsername(username)
	if err != nil {
		return err
	}
	if ou == nil { // 第一次登录
		h := md5.New()
		h.Write([]byte(u.Password))
		u.Password = hex.EncodeToString(h.Sum(nil))
		return us.o.Save(u).Error
	}
	u.Id = ou.Id
	err = us.o.Save(u).Error
	return err
}

func (us *userService) GetUserByUsername(username string) (*user.User, error) {
	u := &user.User{}
	err := us.o.Where("`username`=?", username).First(u).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (us *userService) GetUserByToken(token string) (*user.User, error) {
	u := &user.User{}
	err := us.o.Where("`token`=?", token).First(u).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (us *userService) GetUserById(id int64) (*user.User, error) {
	u := &user.User{}
	err := us.o.Where("`id`=?", id).First(u).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (us *userService) Delete(id int64) error {
	return us.o.Unscoped().Where("`id`=?", id).Delete(user.User{}).Error
}

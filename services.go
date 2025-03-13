package main

import (
	"errors"
	"log"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Register(username string, password string) (User, error) {
	var existUserCnt int64
	if err := DB.Model(&User{}).Where("username = ?", username).Count(&existUserCnt).Error; err != nil {
		log.Printf("Database error: %v", err)
		return User{}, err
	}
	if existUserCnt != 0 {
		return User{}, errors.New("username already exist")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return User{}, err
	}

	newUser := User{Username: username, Password: string(hashedPassword)}
	if err := DB.Create(&newUser).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		return User{}, err
	}
	return newUser, nil
}

func Login(username string, password string) (User, error) {
	var user User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		log.Printf("Invalid username or password %v", err)
		return User{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("Invalid username or password %v", err)
		return User{}, err
	}
	return user, nil
}

func ListUrlMapping(userID uint) ([]URLMapping, error) {
	var urlList []URLMapping
	if err := DB.Where("user_id = ?", userID).Find(&urlList).Error; err != nil {
		log.Printf("Database error: %v", err)
		return []URLMapping{}, err
	}
	return urlList, nil
}

func GetUrlMapping(urlMappingID string) (URLMapping, error) {
	var urlMapping URLMapping
	if err := DB.Where("ID = ?", urlMappingID).Take(&urlMapping).Error; err != nil {
		log.Printf("Error quering URLMapping: %v", err)
		return URLMapping{}, err
	}
	return urlMapping, nil
}

func GetUrlMappingActionRecord(urlMappingID string) (URLMappingActionRecord, error) {
	var urlMappingActionRecord URLMappingActionRecord
	if err := DB.Where("url_mapping_id = ?", urlMappingID).Take(&urlMappingActionRecord).Error; err != nil {
		log.Printf("Error quering URLMappingActionRecord: %v", err)
		return URLMappingActionRecord{}, err
	}
	return urlMappingActionRecord, nil
}

func CreateUrlMapping(userID uint, originURL string) (URLMapping, error) {
	var randomID string
	for randomID == "" {
		randomID = randomString(6)
		var existUrlMappingCnt int64
		if err := DB.Model(&URLMapping{}).Where("ID = ?", randomID).Count(&existUrlMappingCnt).Error; err != nil {
			log.Printf("Database error: %v", err)
			return URLMapping{}, err
		}
		if existUrlMappingCnt != 0 {
			randomID = ""
		}
	}

	newURLMapping := URLMapping{ID: randomID, OriginURL: originURL, UserID: userID}
	if err := DB.Create(&newURLMapping).Error; err != nil {
		log.Printf("Error creating url mapping: %v", err)
		return URLMapping{}, err
	}
	newURLMappingActionRecord := URLMappingActionRecord{URLMapping: newURLMapping, ClickCount: 0}
	if err := DB.Create(&newURLMappingActionRecord).Error; err != nil {
		log.Printf("Error creating url mapping action record: %v", err)
		return newURLMapping, err
	}
	return newURLMapping, nil
}

func Redirect(urlMappingID string) (URLMapping, error) {
	urlMappingChan := make(chan URLMapping, 1)
	urlMappingActionRecordChan := make(chan URLMappingActionRecord, 1)
	errChan := make(chan error, 2)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		var urlMapping URLMapping
		if err := DB.Where("ID =?", urlMappingID).Take(&urlMapping).Error; err != nil {
			errChan <- err
			return
		}
		urlMappingChan <- urlMapping
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var urlMappingActionRecord URLMappingActionRecord
		if err := DB.Where("url_mapping_id = ?", urlMappingID).Take(&urlMappingActionRecord).Error; err != nil {
			errChan <- err
			return
		}
		urlMappingActionRecord.ClickCount += 1
		urlMappingActionRecord.LastAccess = time.Now()
		if err := DB.Save(&urlMappingActionRecord).Error; err != nil {
			errChan <- err
			return
		}
		urlMappingActionRecordChan <- urlMappingActionRecord
	}()

	wg.Wait()
	close(urlMappingChan)
	close(urlMappingActionRecordChan)
	close(errChan)

	for err := range errChan {
		log.Printf("Error: %v", err)
		return URLMapping{}, err
	}

	return <-urlMappingChan, nil
}
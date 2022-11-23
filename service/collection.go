package service

import (
	"fmt"
	"nfthelper/database"
	"nfthelper/logger"
	"nfthelper/model"
	"nfthelper/model/dbmodel"
)

type CollectionService struct {
}

func (c *CollectionService) GetByID(id uint) (collection dbmodel.Collection) {
	msg := fmt.Sprintf("get collection by id=%d", id)

	result := database.DB.Where("id=?", id).First(&collection)
	if result.Error != nil {
		logger.Error("%s, error is %+v", msg, result.Error)
		return
	}
	return
}

func (c *CollectionService) GetByAddr(address string) (collection dbmodel.Collection) {
	msg := fmt.Sprintf("get collection by address=%s", address)

	result := database.DB.Where("address=?", address).First(&collection)
	if result.Error != nil {
		logger.Error("%s, error is %+v", msg, result.Error)
		return
	}
	return
}

func (c *CollectionService) Search(keyword string) (collections []model.Collection) {
	msg := fmt.Sprintf("get collection by keyword=%s", keyword)

	result := database.DB.Where("name like ?", "%"+keyword+"%").Order("price desc").Limit(8).Find(&collections)
	if result.Error != nil {
		logger.Error("%s, error is %+v", msg, result.Error)
		return
	}
	return
}

func (c *CollectionService) HasAlreadyWatched(uid uint, collectionID uint) bool {
	msg := fmt.Sprintf("has user already subscribed the collection uid=%d, collectionID=%d", uid, collectionID)
	logger.Info(msg)

	var count int64
	result := database.DB.Table("user_collection").Where("user_id=? and collection_id=?", uid, collectionID).Count(&count)
	if result.Error != nil {
		logger.Error("%s, error is %+v", msg, result.Error)
	}
	logger.Info("user collection count is %d", count)
	return count > 0
}

func (c *CollectionService) ListByUserID(uid uint) (collections []model.Collection) {
	msg := fmt.Sprintf("list user collections uid=%d", uid)
	logger.Info(msg)

	result := database.DB.Table("user_collection").Select("collection.id, collection.name, collection.address, collection.price, collection.os_url").Joins("left join collection on user_collection.collection_id = collection.id where user_id=?", uid).Scan(&collections)
	if result.Error != nil {
		logger.Error("%s, error is %+v", msg, result.Error)
	}
	return
}

func (c *CollectionService) AddUserCollection(uid uint, collection dbmodel.Collection) {
	msg := fmt.Sprintf("add user collection by collection=%+v", collection)

	userCollection := dbmodel.UserCollection{
		UserID:       uid,
		CollectionID: collection.ID,
	}
	result := database.DB.Save(&userCollection)
	if result.Error != nil {
		logger.Error("%s, error is %+v", msg, result.Error)
	}
}

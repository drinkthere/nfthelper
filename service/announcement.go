package service

import (
	"fmt"
	"nfthelper/database"
	"nfthelper/logger"
	"nfthelper/model/dbmodel"
)

type AnnouncementService struct {
}

func (c *AnnouncementService) ListByCollectionID(collectionID uint) (announcements []dbmodel.Announcement) {
	msg := fmt.Sprintf("list announcemnets by collectionId=%d", collectionID)
	logger.Info(msg)

	result := database.DB.Where("collection_id=?", collectionID).Order("id desc").Limit(3).Find(&announcements)
	if result.Error != nil {
		logger.Error("%s, error is %+v", msg, result.Error)
	}
	return
}

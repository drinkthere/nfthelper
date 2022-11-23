package service

import "nfthelper/model"

type AnnouncementService struct {
}

func (c *AnnouncementService) GetByCollectionIDAndUserID(id uint, uid uint) []model.Announcement {
	return []model.Announcement{
		{
			ID:  1,
			URL: "https://opensea.io/collection/azuki",
		},
		{
			ID:  2,
			URL: "https://opensea.io/collection/cryptopunks",
		},
		{
			ID:  3,
			URL: "https://opensea.io/collection/clonex",
		},
	}
}

package service

import "nfthelper/model"

type CollectionService struct {
}

func (c *CollectionService) GetCollectionByID(id int64) model.Collection {
	return model.Collection{
		ID:      1,
		Name:    "Azuki",
		Address: "0x6C869A43A9D362eF870d75daE56A01887578421d",
		Price:   6.1,
	}
}

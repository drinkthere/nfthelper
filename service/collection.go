package service

import "nfthelper/model"

type CollectionService struct {
}

func (c *CollectionService) GetByID(id int64) model.Collection {
	return model.Collection{
		ID:      1,
		Name:    "Azuki",
		Address: "0x6C869A43A9D362eF870d75daE56A01887578421d",
		Price:   6.1,
		OsURL:   "https://opensea.io/collection/azuki",
	}
}

func (c *CollectionService) GetByAddr(address string) model.Collection {
	return model.Collection{
		ID:      1,
		Name:    "Homa Gang - Valentine (Homa Gang - Valentine)",
		Address: "0x6C869A43A9D362eF870d75daE56A01887578421d",
		Price:   6.1,
		OsURL:   "https://opensea.io/collection/azuki",
	}
}

func (c *CollectionService) Search(keyword string) []model.Collection {
	return []model.Collection{
		{ID: 2, Name: "Azuki"},
		{ID: 3, Name: "AzukiApeSocialClub"},
		{ID: 4, Name: "OkayAzukis"},
	}
}

func (c *CollectionService) ListByUserID(uid int64) []model.Collection {
	return []model.Collection{
		{ID: 2, Name: "Azuki", OsURL: "https://opensea.io/collection/azuki"},
		{ID: 3, Name: "AzukiApeSocialClub", OsURL: "https://opensea.io/collection/azuki"},
		{ID: 4, Name: "OkayAzukis", OsURL: "https://opensea.io/collection/azuki"},
		{ID: 5, Name: "Azuki", OsURL: "https://opensea.io/collection/azuki"},
	}
}

package model

func GetAllGolies() ([]Goly, error) {
	var golies []Goly

	tx := db.Find(&golies)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return golies, nil
}

func GetGoly(id uint64) (Goly, error) {
	var goly Goly

	tx := db.Where("id = ?", id).First(&goly)
	if tx.Error != nil {
		return Goly{}, tx.Error
	}

	return goly, nil
}

func CreateGoly(goly Goly) (Goly, error) {
	tx := db.Create(&goly)
	if tx.Error != nil {
		return Goly{}, tx.Error
	}

	return goly, nil
}

func UpdateGoly(goly Goly) (Goly, error) {
	tx := db.Save(&goly)
	if tx.Error != nil {
		return Goly{}, tx.Error
	}

	return goly, nil
}

func DeleteGoly(id uint64) error {
	tx := db.Unscoped().Delete(&Goly{}, id)
	return tx.Error
}

func FindByGolyUrl(url string) (Goly, error) {
	var goly Goly
	tx := db.Where("goly = ?", url).First(&goly)
	return goly, tx.Error
}

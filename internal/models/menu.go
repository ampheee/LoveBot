package models

var AdminStartMenu = struct {
	AdminEnter string
}{
	AdminEnter: "Войти в меню",
}

var UserStartMenu = struct {
	InsertSomeThoughts string
	GetComplimentNow   string
}{
	InsertSomeThoughts: "Поделиться мыслями",
	GetComplimentNow:   "Получить факт",
}

var AdminMenu = struct {
	InsertNewPhoto      string
	InsertNewCompliment string
	GetAllCompliments   string
	GetAllPhotos        string
	Back                string
}{
	InsertNewPhoto:      "Добавить новое фото",
	InsertNewCompliment: "Добавить новый комплимент",
	GetAllCompliments:   "Все комплименты",
	GetAllPhotos:        "Все фото",
	//Back:              "Назад",
}

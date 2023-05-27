package models

var AdminStartMenu = struct {
	AdminEnter string
}{
	AdminEnter: "Войти в меню",
}

var UserStartMenu = struct {
	InsertSomeThoughts string
	GetComplimentNow   string
	Back               string
}{
	InsertSomeThoughts: "Поделиться мыслями 😳",
	GetComplimentNow:   "Получить факт 😋",
	Back:               "Вернуться назад ❤️‍🩹",
}

var AdminMenu = struct {
	InsertNewPhoto      string
	InsertNewCompliment string
	GetAllCompliments   string
	GetAllPhotos        string
	GetComplimentNow    string
	Back                string
}{
	InsertNewPhoto:      "Добавить новое фото",
	InsertNewCompliment: "Добавить новый комплимент",
	GetAllCompliments:   "Все комплименты",
	GetAllPhotos:        "Все фото",
	//Back:              "Назад",
}

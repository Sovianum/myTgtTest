package model

type StatsSlice struct {
	Items []Item `json:"items"`
}

type Item struct {
	Date Calendar `json:"date"`
	Rows []Row    `json:"rows"`
}

type Row struct {
	Id    uint   `json:"id"`
	Age   uint   `json:"age"`
	Sex   string `json:"sex"`
	Count uint   `json:"count"`
}

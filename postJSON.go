package letter

type PostData struct {
	Name         string `json:"name"`
	RelationShip Person `json:"relationship"`
}

type Person struct {
	Name string `json:"name"`
}

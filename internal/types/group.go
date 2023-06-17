package types

type Group struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	OwnerID uint   `json:"ownerId"`
}
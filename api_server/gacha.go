package main

import (
	"math/rand"
	"sort"
	"strconv"
	"time"
)

type Character struct {
	CharacterID string
	Name        string
}

type Results struct {
	Results []Character `json:"results"`
}

type OwnerShip struct {
	UserCharacterID string `json:"userCharacterID"`
	CharacterID     string `json:"charcterID"`
	Name            string `json:"name"`
}
type Characters struct {
	Characters []OwnerShip
}
type CharacterProbability struct {
	ID          int
	Name        string
	Probability int
}

type Posession struct {
	UserID      int
	CharacterID int
}
type Time struct {
	Times int
}

func GetProbabilityList() ([]CharacterProbability, error) {
	var list []CharacterProbability

	//characterprobabilityテーブルから確率のリストを取得
	q := `SELECT id, name, probability FROM characters`
	rows, err := db.Query(q)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var character CharacterProbability
		if err := rows.Scan(&character.ID, &character.Name, &character.Probability); err != nil {
			return nil, err
		}
		list = append(list, character)
	}

	return list, nil
}

// DrawGacha ガチャを一回引く　キャラクターのIDを返す
func DrawGacha(list []CharacterProbability) int {
	rand.Seed(time.Now().UnixNano())

	boundaries := make([]int, len(list)+1)
	for i := 1; i < len(boundaries); i++ {
		boundaries[i] = boundaries[i-1] + list[i-1].Probability
	}

	x := rand.Intn(boundaries[len(boundaries)-1] + 1)
	idx := sort.SearchInts(boundaries, x) - 1

	return idx
}

func (poses *Posession) RegistrateOwnership() error {
	if _, err := db.Exec(`INSERT INTO characters_possession (userid, characterid) VALUES ( ?, ? )`, poses.UserID, poses.CharacterID); err != nil {
		return err
	}
	return nil
}

func (user *User) GetCharacters() (*Characters, error) {
	q := `SELECT characters_possession.usercharacterid, characters.id, characters.name
			FROM characters INNER JOIN characters_possession
			ON characters.id=characters_possession.characterid
			WHERE characters_possession.userid = ?`

	rows, err := db.Query(q, user.ID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var characters Characters
	var UserCharacterID int
	var CharacterID int
	for rows.Next() {
		owener := OwnerShip{}
		if err := rows.Scan(&UserCharacterID, &CharacterID, &owener.Name); err != nil {
			return nil, err
		}
		owener.UserCharacterID = strconv.Itoa(UserCharacterID)
		owener.CharacterID = strconv.Itoa(CharacterID)
		characters.Characters = append(characters.Characters, owener)
	}

	return &characters, nil
}

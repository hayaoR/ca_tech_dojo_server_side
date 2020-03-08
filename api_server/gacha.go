package main

import (
	"math/rand"
	"sort"
	"time"
)

type Character struct {
	CharacterID string
	Name        string
}

type Results struct {
	Results []Character `json:"results"`
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
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
	if _, err := db.Exec(`INSERT INTO characters_possession VALUES ( ?, ? )`, poses.UserID, poses.CharacterID); err != nil {
		return err
	}
	return nil
}

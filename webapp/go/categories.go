package main

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	inMemoryCategories = map[int]Category{
		1:  {ID: 1, ParentID: 0, CategoryName: "ソファー", ParentCategoryName: ""},
		2:  {ID: 2, ParentID: 1, CategoryName: "一人掛けソファー", ParentCategoryName: "ソファー"},
		3:  {ID: 3, ParentID: 1, CategoryName: "二人掛けソファー", ParentCategoryName: "ソファー"},
		4:  {ID: 4, ParentID: 1, CategoryName: "コーナーソファー", ParentCategoryName: "ソファー"},
		5:  {ID: 5, ParentID: 1, CategoryName: "二段ソファー", ParentCategoryName: "ソファー"},
		6:  {ID: 6, ParentID: 1, CategoryName: "ソファーベッド", ParentCategoryName: "ソファー"},
		10: {ID: 10, ParentID: 0, CategoryName: "家庭用チェア", ParentCategoryName: ""},
		11: {ID: 11, ParentID: 10, CategoryName: "スツール", ParentCategoryName: "家庭用チェア"},
		12: {ID: 12, ParentID: 10, CategoryName: "クッションスツール", ParentCategoryName: "家庭用チェア"},
		13: {ID: 13, ParentID: 10, CategoryName: "ダイニングチェア", ParentCategoryName: "家庭用チェア"},
		14: {ID: 14, ParentID: 10, CategoryName: "リビングチェア", ParentCategoryName: "家庭用チェア"},
		15: {ID: 15, ParentID: 10, CategoryName: "カウンターチェア", ParentCategoryName: "家庭用チェア"},
		20: {ID: 20, ParentID: 0, CategoryName: "キッズチェア", ParentCategoryName: ""},
		21: {ID: 21, ParentID: 20, CategoryName: "学習チェア", ParentCategoryName: "キッズチェア"},
		22: {ID: 22, ParentID: 20, CategoryName: "ベビーソファ", ParentCategoryName: "キッズチェア"},
		23: {ID: 23, ParentID: 20, CategoryName: "キッズハイチェア", ParentCategoryName: "キッズチェア"},
		24: {ID: 24, ParentID: 20, CategoryName: "テーブルチェア", ParentCategoryName: "キッズチェア"},
		30: {ID: 30, ParentID: 0, CategoryName: "オフィスチェア", ParentCategoryName: ""},
		31: {ID: 31, ParentID: 30, CategoryName: "デスクチェア", ParentCategoryName: "オフィスチェア"},
		32: {ID: 32, ParentID: 30, CategoryName: "ビジネスチェア", ParentCategoryName: "オフィスチェア"},
		33: {ID: 33, ParentID: 30, CategoryName: "回転チェア", ParentCategoryName: "オフィスチェア"},
		34: {ID: 34, ParentID: 30, CategoryName: "リクライニングチェア", ParentCategoryName: "オフィスチェア"},
		35: {ID: 35, ParentID: 30, CategoryName: "投擲用椅子", ParentCategoryName: "オフィスチェア"},
		40: {ID: 40, ParentID: 0, CategoryName: "折りたたみ椅子", ParentCategoryName: ""},
		41: {ID: 41, ParentID: 40, CategoryName: "パイプ椅子", ParentCategoryName: "折りたたみ椅子"},
		42: {ID: 42, ParentID: 40, CategoryName: "木製折りたたみ椅子", ParentCategoryName: "折りたたみ椅子"},
		43: {ID: 43, ParentID: 40, CategoryName: "キッチンチェア", ParentCategoryName: "折りたたみ椅子"},
		44: {ID: 44, ParentID: 40, CategoryName: "アウトドアチェア", ParentCategoryName: "折りたたみ椅子"},
		45: {ID: 45, ParentID: 40, CategoryName: "作業椅子", ParentCategoryName: "折りたたみ椅子"},
		50: {ID: 50, ParentID: 0, CategoryName: "ベンチ", ParentCategoryName: ""},
		51: {ID: 51, ParentID: 50, CategoryName: "一人掛けベンチ", ParentCategoryName: "ベンチ"},
		52: {ID: 52, ParentID: 50, CategoryName: "二人掛けベンチ", ParentCategoryName: "ベンチ"},
		53: {ID: 53, ParentID: 50, CategoryName: "アウトドア用ベンチ", ParentCategoryName: "ベンチ"},
		54: {ID: 54, ParentID: 50, CategoryName: "収納付きベンチ", ParentCategoryName: "ベンチ"},
		55: {ID: 55, ParentID: 50, CategoryName: "背もたれ付きベンチ", ParentCategoryName: "ベンチ"},
		56: {ID: 56, ParentID: 50, CategoryName: "ベンチマーク", ParentCategoryName: "ベンチ"},
		60: {ID: 60, ParentID: 0, CategoryName: "座椅子", ParentCategoryName: ""},
		61: {ID: 61, ParentID: 60, CategoryName: "和風座椅子", ParentCategoryName: "座椅子"},
		62: {ID: 62, ParentID: 60, CategoryName: "高座椅子", ParentCategoryName: "座椅子"},
		63: {ID: 63, ParentID: 60, CategoryName: "ゲーミング座椅子", ParentCategoryName: "座椅子"},
		64: {ID: 64, ParentID: 60, CategoryName: "ロッキングチェア", ParentCategoryName: "座椅子"},
		65: {ID: 65, ParentID: 60, CategoryName: "座布団", ParentCategoryName: "座椅子"},
		66: {ID: 66, ParentID: 60, CategoryName: "空気椅子", ParentCategoryName: "座椅子"},
	}
)

func getCategoryByID(q sqlx.Queryer, categoryID int) (category Category, err error) {
	if c, ok := inMemoryCategories[categoryID]; ok {
		return c, nil
	}
	return Category{}, errors.New("category is not found")
}

func getParentCategory(q sqlx.Queryer, base *Category) error {
	if c, ok := inMemoryCategories[base.ParentID]; ok {
		base.ParentCategoryName = c.CategoryName
		return nil
	}
	return errors.New("category is not found")
}

package dao

import (
	"errors"
	"strconv"
)

type Flag struct {
	dbRecord
	Id_   int
	Key   string
	Value interface{}
}

type FlagType string

const (
	FlagTypeString FlagType = "string"
	FlagTypeInt    FlagType = "int"
	FlagTypeBool   FlagType = "boolean"
)

func (main *Dao) GetFlag(key string, typeTo FlagType) (*Flag, error) {
	prefix := "main.Dao.GetFlag"
	query := "SELECT _id, flag_key, flag_value, created_at, updated_at, status FROM news_api_flags WHERE flag_key = ?"

	rows, err := main.db.Query(query, key)
	if err != nil {
		return nil, errors.New(prefix + " (query): " + err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		flag := &Flag{}
		var value string
		rows.Scan(&flag.Id_, &flag.Key, &value, &flag.CreatedAt, &flag.UpdatedAt, &flag.Status)

		switch typeTo {
		case FlagTypeString:
			flag.Value = value

		case FlagTypeInt:
			flag.Value, err = strconv.Atoi(value)
			if err != nil {
				return nil, errors.New(prefix + "Invalid type int: " + err.Error())
			}

		case FlagTypeBool:
			flag.Value, err = strconv.ParseBool(value)
			if err != nil {
				return nil, errors.New(prefix + "Invalid type bool: " + err.Error())
			}
		}

		return flag, nil
	}

	return nil, nil
}

func (main *Dao) SetFlag(key, value string, typeTo FlagType) error {
	prefix := "main.Dao.SetFlag"
	flag, err := main.GetFlag(key, typeTo)
	if err != nil {
		return errors.New(prefix + " (get flag): " + err.Error())
	}

	if flag != nil {
		return main.updateFlag(key, value)
	} else {
		return main.createFlag(key, value)
	}
}

func (main *Dao) updateFlag(key, value string) error {
	prefix := "main.Dao.createFlag"
	query := "UPDATE news_api_flags SET flag_value = ? WHERE flag_key = ?"
	stmt, err := main.db.Prepare(query)
	if err != nil {
		return errors.New(prefix + " (prepare): " + err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(value, key)
	if err != nil {
		return errors.New(prefix + " (exec): " + err.Error())
	}

	return nil
}

func (main *Dao) createFlag(key, value string) error {
	prefix := "main.Dao.updateFlag"
	query := "INSERT INTO news_api_flags (flag_key, flag_value) VALUES (?, ?)"
	stmt, err := main.db.Prepare(query)
	if err != nil {
		return errors.New(prefix + " (Prepare): " + err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(key, value)
	if err != nil {
		return errors.New(prefix + " (Prepare): " + err.Error())
	}

	return nil
}

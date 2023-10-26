package drafts

import (
	"encoding/json"
	"strconv"

	"github.com/Gealber/limengo/repositories/models"
	badger "github.com/dgraph-io/badger/v4"
)

type repo struct {
	db *badger.DB
}

func New(db *badger.DB) *repo {
	return &repo{db: db}
}

func (r *repo) List(id int) ([]models.Draft, error) {
	var drafts []models.Draft
	err := r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(strconv.Itoa(id)))
		if err != nil {
			return models.NotFoundErr
		}

		valCopy, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		return json.Unmarshal(valCopy, &drafts)
	})

	return drafts, err
}

func (r *repo) Get(id int, context string) (*models.Draft, error) {
	var draft models.Draft
	err := r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(strconv.Itoa(id)))
		if err != nil {
			return err
		}

		var drafts []models.Draft
		valCopy, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		err = json.Unmarshal(valCopy, &drafts)
		if err != nil {
			return err
		}

		for _, d := range drafts {
			if d.Context == context {
				draft = d
				return nil
			}
		}

		return models.NotFoundErr
	})

	return &draft, err
}

func (r *repo) Create(draft models.Draft) error {
	return r.db.Update(func(txn *badger.Txn) error {
		key := strconv.Itoa(draft.ID)

		item, err := txn.Get([]byte(key))
		if err != nil {
			// means that the key is not in the db, so we create a brand new resource
			drafts := []models.Draft{draft}
			data, err := json.Marshal(drafts)
			if err != nil {
				return err
			}

			return txn.Set([]byte(key), data)
		}

		valCopy, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		var drafts []models.Draft
		err = json.Unmarshal(valCopy, &drafts)
		if err != nil {
			return err
		}

		for _, d := range drafts {
			if d.Context == draft.Context {
				return models.DuplicateValueErr
			}
		}

		drafts = append(drafts, draft)

		data, err := json.Marshal(drafts)
		if err != nil {
			return err
		}

		return txn.Set([]byte(key), data)
	})
}

func (r *repo) Delete(id int, context string) error {
	return r.db.Update(func(txn *badger.Txn) error {
		key := strconv.Itoa(id)

		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		valCopy, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		var drafts []models.Draft
		err = json.Unmarshal(valCopy, &drafts)
		if err != nil {
			return err
		}

		pos := -1
		for i, d := range drafts {
			if d.Context == context {
				pos = i
				break
			}
		}

		if pos == -1 {
			return models.DuplicateValueErr
		}

		drafts = append(drafts[:pos], drafts[pos+1:]...)

		data, err := json.Marshal(drafts)
		if err != nil {
			return err
		}

		return txn.Set([]byte(key), data)
	})
}

func (r *repo) Update(draft models.Draft) (*models.Draft, error) {
	var result models.Draft
	err := r.db.Update(func(txn *badger.Txn) error {
		key := strconv.Itoa(draft.ID)

		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		valCopy, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		var drafts []models.Draft
		err = json.Unmarshal(valCopy, &drafts)
		if err != nil {
			return err
		}

		pos := -1
		for i, d := range drafts {
			if d.Context == draft.Context {
				pos = i
				break
			}
		}

		if pos == -1 {
			return models.NotFoundErr
		}

		// overwritting draft
		drafts[pos].Type = draft.Type
		drafts[pos].Data = draft.Data

		result = drafts[pos]

		data, err := json.Marshal(drafts)
		if err != nil {
			return err
		}

		return txn.Set([]byte(key), data)
	})

	return &result, err
}

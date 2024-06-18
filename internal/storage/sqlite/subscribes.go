package sqlite

import "fmt"

func (s *Storage) Subscribe(subID, uID int64) (int64, error) {
	const fn = "storage.sqlite.Subscribe"

	stmt, err := s.db.Prepare("INSERT INTO subscribes (user_id, sub_id) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(uID, subID)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	return id, nil
}

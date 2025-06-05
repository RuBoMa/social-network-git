package seed

import "database/sql"

func RunSeeds(db *sql.DB) {
	err := SeedUsers(db)
	if err == nil {
		err = SeedPosts(db)
		if err == nil {
			err = SeedGroups(db)
			if err == nil {
				_ = SeedMessages(db)
			}
		}
	}
}

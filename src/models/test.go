package models

// TestUsers runs some mock actions.
func TestUsers(users IUsers) error {
	// adding users
	userList := []User{
		User{
			DomainID: 1,
			Email:    "email1@example.com",
			Hash:     "password123",
		},
		User{
			DomainID: 1,
			Email:    "email2@example.com",
			Hash:     "password123",
		},
		User{
			DomainID: 2,
			Email:    "email1@example.com",
			Hash:     "password123",
		},
	}
	for _, u := range userList {
		_, err := users.Add(u.Email, u.DomainID, u.Hash)
		if err != nil {
			return err
		}
	}

	// testing user funcs
	user := User{
		DomainID: 2,
		Email:    "email1@example.com",
	}

	id, err := users.Authenticate(user.Email, user.DomainID, "password123")
	if err != nil {
		return err
	}

	user.ID = id

	err = users.ChangePassword(user.ID, "123password")
	if err != nil {
		return err
	}

	_, err = users.Authenticate(user.Email, user.DomainID, "123password")
	if err != nil {
		return err
	}

	return nil
}


// TestDomains runs some mock actions
func TestDomains(domains IDomains) error {
	// adding domains
	domainList := []Domain{
		Domain{
			Name: "name1",
			Data: "{}",
		},
		Domain{
			Name: "name2",
			Data: "{}",
		},
	}
	for _, d := range domainList {
		_, err := domains.Add(d.Name, d.Data)
		if err != nil {
			return err
		}
	}

	// testing domain funcs
	domains.UpdateData(1, `{"updated": true}`)

	return nil
}

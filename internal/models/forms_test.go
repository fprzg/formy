package models

/*
func TestFormsInsert(t *testing.T) {
	_, m := setupTestDB(t)

	err := m.Users.Insert("Alice", "alice@example.com", "securepass")
	assert.NoError(t, err)

	userID, err := m.Users.Authenticate("alice@example.com", "securepass")
	assert.NoError(t, err)
	assert.NotZero(t, userID)

	const fields = `
	[
		{ "field_name": "name", "field_type": "string", "constraints": ["required"] },
		{ "field_name": "email", "field_type": "string", "constraints": ["unique", "required"] },
		{ "field_name": "phone_number", "field_type": "string", "constraints": ["required"] },
		{ "field_name": "message", "field_type": "string", "constraints": [] }
	]
	`
	err = m.Forms.InsertForm(userID, "SimpleForm", "Simple contact form.", fields)
	assert.NoError(t, err)

	forms, err := m.Forms.GetFormsByUser(userID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(forms))
}
*/

/*
func TestFormsGet(t *testing.T) {
	ctx := setupTestDB(t)
	m := GetModels(ctx.AppDB)

	err := m.Users.Insert("Alice", "alice@example.com", "securepass")
	assert.NoError(t, err)

	userID, err := m.Users.Authenticate("alice@example.com", "securepass")
	assert.NoError(t, err)
	assert.NotZero(t, err)

	err = m.Forms.Insert(userID, "SimpleForm", "This is a simple form")
	assert.NoError(t, err)

}

*/

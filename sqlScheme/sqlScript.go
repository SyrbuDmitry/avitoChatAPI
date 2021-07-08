package sqlScheme

//Схема базы данных

const (
	userTable = `CREATE TABLE IF NOT EXISTS users (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					name VARCHAR(100) NOT NULL,
					createdAt Date
				 );`

	chatTable = `CREATE TABLE IF NOT EXISTS chats (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					name VARCHAR(100) NOT NULL,
					createdAt Date
				);`

	messageTable = `CREATE TABLE IF NOT EXISTS messages (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						text VARCHAR(220), 
						chat_id INTEGER REFERENCES chats(id) NOT NULL,
						user_id INTEGER REFERENCES users(id) NOT NULL,
						createdAt Date
                    );`

	chatsJoinUsersTable = `CREATE TABLE IF NOT EXISTS chatsJoinUsers (
							chat_id INTEGER REFERENCES chats(id) NOT NULL,
							user_id INTEGER REFERENCES users(id) NOT NULL
                           );`

	Scheme = userTable + chatTable + messageTable + chatsJoinUsersTable
)

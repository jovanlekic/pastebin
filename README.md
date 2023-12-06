# Pastebin

Pastebin is an online platform that allows users to store and share plain text, code snippets, and other types of content. It provides a convenient way to share snippets of code, configuration files, or any text-based information with others, often used by programmers, developers, and individuals looking to exchange information in a simple, temporary, and publicly accessible manner.

# DB part
- Install Postgres locally
- Get MongoDB connection string and paste it into ".env" file as MONGO_URI
- **(MONGO)Enable current IP - Enable current IP address**
- Use this to create new db with collection
`dbObj := db.NewMongoDB(client, "sample_mflix", "movies")`
- Use this to create new db with collection
`postgresClient, err := db.ConnectToPostgresDb("dbName", "user", password")`

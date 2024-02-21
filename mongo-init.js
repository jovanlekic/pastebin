db = db.getSiblingDB('paste');
db.createUser({
    user: 'admin',
    pwd: 'secret',
    roles: [{ role: 'readWrite', db: 'paste' }]
});
db.createCollection('messages');

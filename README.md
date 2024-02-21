# Pastebin Backend endpoint

Pastebin is an online platform that allows users to store and share plain text (for now, idea is to make additional updates to support code sharing, also to recognize what type of content is stored, etc.). It provides a convenient way to share text-based information with others, often used by individuals looking to exchange information in a simple, temporary, and publicly accessible manner.

## Table of Contents
- Starting point
- API
- DB
- KGS
- Final words

### Starting point
- Whole app works in 3 containers, 1 server and 1 for PostgresDB, and 1 for MongoDB
- Run `docker-compose up -d`
- App is available at port 8080

### API 
| Path | Type | Explaination |
| ------------ | ------------- | ------------- |
| /api/register | POST  | User registration |
| /api/login | POST  | User login |
| /api/check | GET | Check Authorization |
| /api/checkandparse | GET  | Check Authorization |
| /api/createPaste | POST  | Create Paste |
| /api//getPaste/{pasteKey} | GET  | Get Paste by key |
| /api/deletePaste | POST  | Delete Paste |
| /api/getUserInfo | GET  | Get user metadata |
| /api/getUserPastes | GET  | Get user pastes |

### DB 
- Handle all necessary CRUD operations needed for this API actions.

### KGS
- Detached entity made to work only as key generator service, it populates its table with all combinations of keys and gives free key each time.

### Final words
- It's important to mention that whole app is made to serve request sequentually, and ofcourse its could be speed up with starting new goroutine each time new request comes, or choosing more complex architecture solution with multiplicating servers, adding caches, load balancers, etc..
- App will gradually be developed into more robust one, currently its made monolithic but in the future it's planed to decouple certain parts into microservices.
- It's planned to include code sharing, with code detection.
- We are almost done with frontend in React.
- So, till next update.. (Feb 2024)

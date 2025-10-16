```mermaid
erDiagram
    client {
        string[uuid7] id PK
        string email
        float balance
    }
    message {
        string id
        string sender_client_id FK
        string receiver_phone "e.g. +989123456789"
        text content
        float price
        boolean is_express
        string status "pending, sent, failed"
		timestamp created_at
		timestamp updated_at
    }
	client 1 to 1+ message : "Has many"
```
# Схема БД HoneyGame

## Обзор

Проект использует PostgreSQL с отдельной схемой `honey`.

## Таблицы

```mermaid
erDiagram
    users ||--o{ user_honey : "имеет"
    users ||--o{ challenges : "участвует"
    users ||--o{ challenges : "победил"
    users ||--o{ challenges : "проиграл"

    users {
        int id PK
        bigint tg_chat_id
        varchar user_name
    }

    user_honey {
        int user_id FK
        bigint honey
    }

    challenges {
        int id PK
        bigint amount
        int winner_user_id FK
        int loser_user_id FK
        timestamp created_at
        timestamp completed_at
    }
```

## Детали таблиц

### `honey.users`

| Поле | Тип | Описание |
|------|-----|----------|
| id | SERIAL | Первичный ключ |
| tg_chat_id | BIGINT | ID чата в Telegram (уникальный) |
| user_name | VARCHAR(255) | Имя пользователя |

### `honey.user_honey`

| Поле | Тип | Описание |
|------|-----|----------|
| user_id | INTEGER | Внешний ключ на users (PK) |
| honey | BIGINT | Баланс меда |

### `honey.challenges`

| Поле | Тип | Описание |
|------|-----|----------|
| id | SERIAL | Первичный ключ |
| amount | BIGINT | Количество меда в ставке |
| winner_user_id | INTEGER | ID победителя |
| loser_user_id | INTEGER | ID проигравшего |
| created_at | TIMESTAMP | Время создания вызова |
| completed_at | TIMESTAMP | Время завершения (NULL — не завершён) |

## Индексы

- `idx_challenges_winner` — по `winner_user_id`
- `idx_challenges_loser` — по `loser_user_id`
- `idx_challenges_completed_at` — по `completed_at`

## Миграции

- **up:** `migrations/up.sql`
- **down:** `migrations/down.sql`

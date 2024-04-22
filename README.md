# atlas-logout-gift
Mushroom Game - Logout Gift Service

## Overview
A RESTful resource which provides logout gifts to characters. This is backed by a Postgres database, but is seeded by a JSON file.
## Environment
- JAEGER_HOST - Jaeger [host]:[port]
- LOG_LEVEL - Logging level - Panic / Fatal / Error / Warn / Info / Debug / Trace
- DB_USER - Postgres user name
- DB_PASSWORD - Postgres user password
- DB_HOST - Postgres Database host
- DB_PORT - Postgres Database port
- DB_NAME - Postgres Database name

## API
    GET - /api/lgs/worlds/{worldId}/gifts
    POST - /api/lgs/worlds/{worldId}/characters/{characterId}/choices
    DELETE - /api/lgs/worlds/{worldId}/characters/{characterId}/choices
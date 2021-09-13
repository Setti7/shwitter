Shwitter

## Database migrations

We are using [Journey](https://github.com/db-journey/journey) to manage the database migrations.

Common commands:

```bash
make jobs             # Prints the base Journey command.
make create-keyspace  # Create the keyspace for Cassandra.
make migrate-up       # Apply all migrations.
make create-migration # Creates an empty migration.
```

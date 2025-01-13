## SQLSlap 👋

It's a CLI tool to generate fake data for database and simulate real world load test.

### Features

- [x] Fill database with fake data
- [x] Random Insert/Delete/Update operation
- [ ] Slow Query Simulation
- [ ] Deadlock Simulation

### Setup

Build and install on system

```bash
go build -o sqlslap .
mv ./sqlslap /usr/bin
```

SQLSlap requires a single `config.yaml` to provide it the metadata and config.

```yaml
connection_info:
  host: <db_host_ip>
  port: 3306
  user: root
  password: toor
databases:
  - name: bloggers_world
    schema: >
      CREATE TABLE users (
          id CHAR(36) NOT NULL, -- UUID
          username VARCHAR(255) NOT NULL,
          name VARCHAR(255) NOT NULL,
          email VARCHAR(255) NOT NULL,
          created_at DATETIME NOT NULL,
          PRIMARY KEY (id)
      ) ENGINE=InnoDB;

      CREATE TABLE posts (
          id CHAR(36) NOT NULL, -- UUID
          title VARCHAR(255) NOT NULL,
          content TEXT NOT NULL,
          created_at DATETIME NOT NULL,
          user_id CHAR(36) NOT NULL, -- Foreign key to users.id
          PRIMARY KEY (id),
          KEY user_id (user_id),
          CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
      ) ENGINE=InnoDB;

    tables:
      users:
        engine: InnoDB
        columns:
          - name: id
            generator: uuid
            primary_key: true
          - name: username
            generator: username
          - name: name
            generator: name
          - name: email
            generator: email
          - name: created_at
            generator: datetime
      posts:
        engine: InnoDB
        columns:
          - name: id
            generator: uuid
          - name: title
            generator: title
          - name: content
            generator: sentence
          - name: created_at
            generator: datetime
          - name: user_id
            generator: fk(users.id)
```

### Usage

- Cleanup and Re-create database tables
  ```bash
  sqlslap --db bloggers_world init-db
  ```
- Insert rows in database with generated fake datas

  ```bash
  sqlslap --db bloggers_world random insert --table users --clients 100
  ```

  **Output Format** -

  ```
  Job: Insert Job | Runners : 10
  ┌────────┬──────────┬─────────┬──────────────┐
  │ Runner │  Total   │ Failed  │ Op/s         │
  ├────────┼──────────┼─────────┼──────────────┤
  │ 0      │ 460      │ 0       │ 28.75        │
  │ 1      │ 433      │ 0       │ 27.06        │
  │ 2      │ 428      │ 0       │ 26.75        │
  │ 3      │ 430      │ 0       │ 26.87        │
  │ 4      │ 460      │ 0       │ 28.75        │
  │ 5      │ 460      │ 0       │ 28.75        │
  │ 6      │ 459      │ 0       │ 28.68        │
  │ 7      │ 434      │ 0       │ 27.12        │
  │ 8      │ 433      │ 0       │ 27.06        │
  │ 9      │ 460      │ 0       │ 28.75        │
  ├────────┴──────────┴─────────┴──────────────┤
  │ Total  │ 4457     │ 0       │ 278.6        │
  └────────────────────────────────────────────┘
  ```

- Delete random rows
  ```bash
  sqlslap --db bloggers_world random delete --table users --clients 100
  ```
- Update random rows
  ```bash
  sqlslap --db bloggers_world random update --table users --clients 100
  ```

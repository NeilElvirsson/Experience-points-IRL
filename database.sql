CREATE TABLE IF NOT EXISTS user (
    id TEXT NOT NULL,
    user_name TEXT NOT NULL UNIQUE,
    level INTEGER,
    password TEXT NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS task (
    id TEXT NOT NULL,
    task_name TEXT NOT NULL UNIQUE,
    xp_value INTEGER,
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS log (
    id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    task_id TEXT NOT NULL,
    timestamp INTEGER,
    PRIMARY KEY(id),
    FOREIGN KEY(user_id) REFERENCES user(id),
    FOREIGN KEY(task_id) REFERENCES task(id) 
);
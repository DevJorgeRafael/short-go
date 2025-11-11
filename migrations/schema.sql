-- 👤 AUTH CONTEXT

CREATE TABLE users (
    id          TEXT PRIMARY KEY,           -- UUID
    email       TEXT UNIQUE NOT NULL,
    password    TEXT NOT NULL,              -- bcrypt hash
    name        TEXT NOT NULL,
    is_active   BOOLEAN DEFAULT TRUE,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
    id              TEXT PRIMARY KEY,
    user_id         TEXT NOT NULL,
    refresh_token   TEXT NOT NULL,
    expires_at      TIMESTAMP NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_refresh_token ON sessions(refresh_token);

-- ✅ TASKS CONTEXT

-- Tabla de catálogo: Estados de tareas
CREATE TABLE task_statuses (
    id          INTEGER PRIMARY KEY,
    code        TEXT UNIQUE NOT NULL,       -- PENDING, IN_PROGRESS, COMPLETED
    name        TEXT NOT NULL,              -- "Pendiente", "En Progreso", "Completada"
    description TEXT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Datos iniciales de estados
INSERT INTO task_statuses (id, code, name, description) VALUES
(1, 'PENDING', 'Pendiente', 'Tarea pendiente por iniciar'),
(2, 'IN_PROGRESS', 'En Progreso', 'Tarea en proceso'),
(3, 'COMPLETED', 'Completada', 'Tarea finalizada');

-- Tabla de catálogo: Prioridades
CREATE TABLE task_priorities (
    id          INTEGER PRIMARY KEY,
    code        TEXT UNIQUE NOT NULL,       -- LOW, MEDIUM, HIGH
    name        TEXT NOT NULL,              -- "Baja", "Media", "Alta"
    level       INTEGER NOT NULL,           -- 1=LOW, 2=MEDIUM, 3=HIGH (para ORDER BY)
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Datos iniciales de prioridades
INSERT INTO task_priorities (id, code, name, level) VALUES
(1, 'LOW', 'Baja', 1),
(2, 'MEDIUM', 'Media', 2),
(3, 'HIGH', 'Alta', 3);

-- Tabla principal: Tareas
CREATE TABLE tasks (
    id              TEXT PRIMARY KEY,           -- UUID
    user_id         TEXT NOT NULL,              -- FK → users
    title           TEXT NOT NULL,
    description     TEXT,
    status_id       INTEGER NOT NULL DEFAULT 1, -- FK → task_statuses (default: PENDING)
    priority_id     INTEGER NOT NULL DEFAULT 2, -- FK → task_priorities (default: MEDIUM)
    starts_at       TIMESTAMP,
    due_date        TIMESTAMP,
    completed_at    TIMESTAMP,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (status_id) REFERENCES task_statuses(id),
    FOREIGN KEY (priority_id) REFERENCES task_priorities(id)
);

-- Índices optimizados
CREATE INDEX idx_tasks_user_id ON tasks(user_id);
CREATE INDEX idx_tasks_status_id ON tasks(status_id);
CREATE INDEX idx_tasks_priority_id ON tasks(priority_id);
CREATE INDEX idx_tasks_due_date ON tasks(due_date);
CREATE INDEX idx_tasks_user_status ON tasks(user_id, status_id);
CREATE INDEX idx_tasks_user_priority ON tasks(user_id, priority_id);

-- Vista opcional para queries más simples (JOIN automático)
CREATE VIEW v_tasks_detailed AS
SELECT 
    t.id,
    t.user_id,
    t.title,
    t.description,
    ts.code as status_code,
    ts.name as status_name,
    tp.code as priority_code,
    tp.name as priority_name,
    tp.level as priority_level,
    t.starts_at,
    t.due_date,
    t.completed_at,
    t.created_at,
    t.updated_at,
    CASE 
        WHEN t.due_date < datetime('now') AND t.status_id != 3 
        THEN 1 
        ELSE 0 
    END as is_overdue
FROM tasks t
INNER JOIN task_statuses ts ON t.status_id = ts.id
INNER JOIN task_priorities tp ON t.priority_id = tp.id;
-- Mini-MES 初始 Schema
-- 可重複執行（IF NOT EXISTS）

CREATE TABLE IF NOT EXISTS recipe (
    id               SERIAL PRIMARY KEY,
    name             VARCHAR(100) NOT NULL UNIQUE,
    equipment_type   VARCHAR(20)  NOT NULL CHECK (equipment_type IN ('CVD','Etch','CMP','Diffusion')),
    target_temp      FLOAT        NOT NULL,
    target_pressure  FLOAT        NOT NULL,
    duration_min     INT          NOT NULL,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS lot (
    id           SERIAL PRIMARY KEY,
    lot_number   VARCHAR(50)  NOT NULL UNIQUE,
    product      VARCHAR(100) NOT NULL,
    recipe_id    INT          NOT NULL REFERENCES recipe(id),
    priority     INT          NOT NULL DEFAULT 3 CHECK (priority BETWEEN 1 AND 5),
    status       VARCHAR(20)  NOT NULL DEFAULT 'QUEUED'
                     CHECK (status IN ('QUEUED','RUNNING','COMPLETED','ON_HOLD')),
    wafer_count  INT          NOT NULL DEFAULT 25,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS equipment (
    id             SERIAL PRIMARY KEY,
    name           VARCHAR(50)  NOT NULL UNIQUE,
    type           VARCHAR(20)  NOT NULL CHECK (type IN ('CVD','Etch','CMP','Diffusion')),
    status         VARCHAR(20)  NOT NULL DEFAULT 'IDLE'
                       CHECK (status IN ('IDLE','RUNNING','DOWN','PM')),
    current_lot_id INT          REFERENCES lot(id) ON DELETE SET NULL,
    utilization    FLOAT        NOT NULL DEFAULT 0,
    temperature    FLOAT        NOT NULL DEFAULT 25,
    pressure       FLOAT        NOT NULL DEFAULT 0,
    ucl_temp       FLOAT        NOT NULL DEFAULT 100,
    lcl_temp       FLOAT        NOT NULL DEFAULT 0,
    ucl_pressure   FLOAT        NOT NULL DEFAULT 1000,
    lcl_pressure   FLOAT        NOT NULL DEFAULT 0,
    is_alarm       BOOLEAN      NOT NULL DEFAULT FALSE,
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS wafer (
    id          SERIAL PRIMARY KEY,
    lot_id      INT         NOT NULL REFERENCES lot(id) ON DELETE CASCADE,
    sequence    INT         NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'QUEUED'
                    CHECK (status IN ('QUEUED','PROCESSING','COMPLETED','DEFECT')),
    defect_flag BOOLEAN     NOT NULL DEFAULT FALSE,
    UNIQUE (lot_id, sequence)
);

CREATE TABLE IF NOT EXISTS spc_record (
    id           SERIAL PRIMARY KEY,
    equipment_id INT         NOT NULL REFERENCES equipment(id) ON DELETE CASCADE,
    parameter    VARCHAR(20) NOT NULL CHECK (parameter IN ('temperature','pressure')),
    value        FLOAT       NOT NULL,
    ucl          FLOAT       NOT NULL,
    lcl          FLOAT       NOT NULL,
    is_alarm     BOOLEAN     NOT NULL DEFAULT FALSE,
    timestamp    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS alarm_event (
    id             SERIAL PRIMARY KEY,
    equipment_id   INT         NOT NULL REFERENCES equipment(id) ON DELETE CASCADE,
    parameter      VARCHAR(20) NOT NULL,
    value          FLOAT       NOT NULL,
    ucl            FLOAT       NOT NULL,
    lcl            FLOAT       NOT NULL,
    severity       VARCHAR(10) NOT NULL CHECK (severity IN ('WARNING','CRITICAL')),
    acknowledged   BOOLEAN     NOT NULL DEFAULT FALSE,
    timestamp      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 常用查詢 index
CREATE INDEX IF NOT EXISTS idx_equipment_status      ON equipment (status);
CREATE INDEX IF NOT EXISTS idx_equipment_type        ON equipment (type);
CREATE INDEX IF NOT EXISTS idx_lot_status            ON lot (status);
CREATE INDEX IF NOT EXISTS idx_spc_record_equipment  ON spc_record (equipment_id, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_alarm_event_equipment ON alarm_event (equipment_id);
CREATE INDEX IF NOT EXISTS idx_alarm_acknowledged    ON alarm_event (acknowledged);

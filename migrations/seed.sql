-- Seed 資料：Recipe 8 筆、Equipment 22 台、Lot 10 筆
-- 可重複執行（ON CONFLICT DO NOTHING）

-- ===== Recipe =====
INSERT INTO recipe (name, equipment_type, target_temp, target_pressure, duration_min) VALUES
    ('CVD-SiO2-STD',    'CVD',       680,  350,  90),
    ('CVD-SiN-STD',     'CVD',       700,  370,  120),
    ('Etch-Poly-DRY',   'Etch',       45,    8,   30),
    ('Etch-Oxide-DRY',  'Etch',       38,    6,   25),
    ('CMP-STI-STD',     'CMP',        60,  180,   45),
    ('CMP-Metal-STD',   'CMP',        65,  200,   60),
    ('Diff-N-Well',     'Diffusion', 950,  760,  180),
    ('Diff-P-Well',     'Diffusion', 960,  755,  180)
ON CONFLICT (name) DO NOTHING;

-- ===== Lot =====
INSERT INTO lot (lot_number, product, recipe_id, priority, status, wafer_count) VALUES
    ('LOT-2024-001', 'NAND-128G',  1, 1, 'RUNNING',   25),
    ('LOT-2024-002', 'DRAM-8G',    2, 2, 'RUNNING',   25),
    ('LOT-2024-003', 'NAND-128G',  1, 1, 'RUNNING',   25),
    ('LOT-2024-010', 'Logic-28nm', 3, 2, 'RUNNING',   25),
    ('LOT-2024-011', 'DRAM-8G',    4, 3, 'RUNNING',   25),
    ('LOT-2024-020', 'NAND-256G',  5, 1, 'RUNNING',   25),
    ('LOT-2024-030', 'Logic-14nm', 7, 1, 'RUNNING',   25),
    ('LOT-2024-040', 'NAND-128G',  1, 3, 'QUEUED',    25),
    ('LOT-2024-041', 'DRAM-16G',   2, 2, 'QUEUED',    25),
    ('LOT-2024-050', 'Logic-28nm', 3, 2, 'COMPLETED', 25)
ON CONFLICT (lot_number) DO NOTHING;

-- ===== Equipment =====
INSERT INTO equipment (name, type, status, current_lot_id, utilization, temperature, pressure, ucl_temp, lcl_temp, ucl_pressure, lcl_pressure, is_alarm) VALUES
    -- CVD 群組（6台）
    ('CVD-01', 'CVD', 'RUNNING', (SELECT id FROM lot WHERE lot_number='LOT-2024-001'), 85, 680, 350, 720, 640, 400, 300, false),
    ('CVD-02', 'CVD', 'RUNNING', (SELECT id FROM lot WHERE lot_number='LOT-2024-002'), 92, 695, 370, 720, 640, 400, 300, false),
    ('CVD-03', 'CVD', 'RUNNING', (SELECT id FROM lot WHERE lot_number='LOT-2024-003'), 78, 735, 360, 720, 640, 400, 300, true),
    ('CVD-04', 'CVD', 'IDLE',    NULL,                                                 45,  25,  50, 720, 640, 400, 300, false),
    ('CVD-05', 'CVD', 'PM',      NULL,                                                  0,  30,  10, 720, 640, 400, 300, false),
    ('CVD-06', 'CVD', 'DOWN',    NULL,                                                  0,  28,   8, 720, 640, 400, 300, false),
    -- Etch 群組（6台）
    ('Etch-01', 'Etch', 'RUNNING', (SELECT id FROM lot WHERE lot_number='LOT-2024-010'), 88, 45, 8,  80, 20, 15, 3, false),
    ('Etch-02', 'Etch', 'RUNNING', (SELECT id FROM lot WHERE lot_number='LOT-2024-011'), 76, 38, 6,  80, 20, 15, 3, false),
    ('Etch-03', 'Etch', 'RUNNING', NULL,                                                  91, 52, 18, 80, 20, 15, 3, true),
    ('Etch-04', 'Etch', 'IDLE',    NULL,                                                  50, 22,  2, 80, 20, 15, 3, false),
    ('Etch-05', 'Etch', 'IDLE',    NULL,                                                  42, 25,  3, 80, 20, 15, 3, false),
    ('Etch-06', 'Etch', 'DOWN',    NULL,                                                   0, 20,  1, 80, 20, 15, 3, false),
    -- CMP 群組（5台）
    ('CMP-01', 'CMP', 'RUNNING', (SELECT id FROM lot WHERE lot_number='LOT-2024-020'), 82, 60, 180, 90, 40, 220, 140, false),
    ('CMP-02', 'CMP', 'RUNNING', NULL,                                                  95, 65, 200, 90, 40, 220, 140, false),
    ('CMP-03', 'CMP', 'IDLE',    NULL,                                                  55, 30,  50, 90, 40, 220, 140, false),
    ('CMP-04', 'CMP', 'PM',      NULL,                                                   0, 25,  10, 90, 40, 220, 140, false),
    ('CMP-05', 'CMP', 'IDLE',    NULL,                                                  38, 28,  45, 90, 40, 220, 140, false),
    -- Diffusion 群組（5台）
    ('Diff-01', 'Diffusion', 'RUNNING', (SELECT id FROM lot WHERE lot_number='LOT-2024-030'), 90, 950, 760, 1000, 900, 800, 720, false),
    ('Diff-02', 'Diffusion', 'RUNNING', NULL,                                                  87, 960, 755, 1000, 900, 800, 720, false),
    ('Diff-03', 'Diffusion', 'IDLE',    NULL,                                                  60,  25, 760, 1000, 900, 800, 720, false),
    ('Diff-04', 'Diffusion', 'DOWN',    NULL,                                                   0,  25, 760, 1000, 900, 800, 720, false),
    ('Diff-05', 'Diffusion', 'RUNNING', NULL,                                                  83, 945, 758, 1000, 900, 800, 720, false)
ON CONFLICT (name) DO NOTHING;

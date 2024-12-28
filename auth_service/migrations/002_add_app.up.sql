INSERT INTO apps (id, name, secret) VALUES (1, 'file_service', 'top-secret') ON CONFLICT DO NOTHING;  
INSERT INTO apps (id, name, secret) VALUES (2, 'museum_admin', 'top1-secret') ON CONFLICT DO NOTHING;

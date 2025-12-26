INSERT INTO public.tag (id, version, name, created_date, last_modified_date)
VALUES
  (100, 1, 'backend', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (101, 1, 'urgent', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

SELECT setval('tag_id_seq', (SELECT COALESCE(MAX(id), 1) FROM public.tag));

INSERT INTO public.item (id, version, status, description, assignee_id, created_date, last_modified_date)
VALUES
  (100, 1, 'TODO', 'First task', 100, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (101, 1, 'DONE', 'Second task', 101, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

SELECT setval('item_id_seq', (SELECT COALESCE(MAX(id), 1) FROM public.item));

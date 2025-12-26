INSERT INTO public.person (id, version, first_name, last_name, created_date, last_modified_date)
VALUES
  (100, 1, 'John', 'Doe', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (101, 1, 'Jane', 'Doe', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

SELECT setval('person_id_seq', (SELECT COALESCE(MAX(id), 1) FROM public.person));

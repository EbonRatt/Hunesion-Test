INSERT INTO public.item_tag (id, item_id, tag_id)
VALUES
  (100, 100, 100),
  (101, 100, 101),
  (102, 101, 100)
ON CONFLICT (id) DO NOTHING;

SELECT setval('item_tag_id_seq', (SELECT COALESCE(MAX(id), 1) FROM public.item_tag));

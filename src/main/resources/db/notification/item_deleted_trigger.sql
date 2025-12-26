CREATE OR REPLACE FUNCTION public.notify_item_deleted()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  PERFORM pg_notify('item_deleted', OLD.id::text);
  RETURN OLD;
END;
$$;

DROP TRIGGER IF EXISTS trg_item_deleted ON public.item;
CREATE TRIGGER trg_item_deleted
AFTER DELETE ON public.item
FOR EACH ROW
EXECUTE FUNCTION public.notify_item_deleted();

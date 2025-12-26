CREATE OR REPLACE FUNCTION public.notify_item_saved()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  PERFORM pg_notify('item_saved', NEW.id::text);
  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS trg_item_saved ON public.item;
CREATE TRIGGER trg_item_saved
AFTER INSERT OR UPDATE ON public.item
FOR EACH ROW
EXECUTE FUNCTION public.notify_item_saved();

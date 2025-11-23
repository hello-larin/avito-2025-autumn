-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_merged_at_on_merge()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = 'MERGED' AND OLD.merged_at IS NULL THEN
        NEW.merged_at := NOW();
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_set_merged_at
    BEFORE UPDATE OF status ON pull_requests
    FOR EACH ROW
    EXECUTE FUNCTION set_merged_at_on_merge();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_set_merged_at ON pull_requests;
DROP FUNCTION IF EXISTS set_merged_at_on_merge();
-- +goose StatementEnd

-- Restores the unqualified body from 000038.
--
-- Note this reintroduces the pg_upgrade failure: a major version upgrade cannot
-- complete while the body is unqualified.
CREATE OR REPLACE FUNCTION immutable_unaccent(text) RETURNS text
  LANGUAGE sql IMMUTABLE STRICT PARALLEL SAFE AS
$$ SELECT unaccent('unaccent', $1) $$;

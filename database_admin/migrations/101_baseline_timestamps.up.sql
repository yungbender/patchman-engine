ALTER TABLE baseline ADD COLUMN creator TEXT CHECK (NOT empty(creator));
ALTER TABLE baseline ADD COLUMN published TIMESTAMP WITH TIME ZONE;
ALTER TABLE baseline ADD COLUMN last_edited TIMESTAMP WITH TIME ZONE;
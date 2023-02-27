create or replace function enable_uuid_extension()
    returns void
    language plpgsql
as
    $function$
declare
v_owner text;
begin
begin
            v_owner := _dba.get_fa_role();
exception
            when undefined_function or invalid_schema_name then
                v_owner := null;
end;

        if v_owner is not null
        then
            perform _dba.create_extension('uuid-ossp', 'public');
else
            create extension if not exists "uuid-ossp" with schema public;
end if;
end;
$function$;

select enable_uuid_extension();

drop function if exists enable_uuid_extension();
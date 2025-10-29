-- Schema for storing blockchain wallet metadata alongside Supabase auth users.
create table if not exists public.wallets (
    supabase_user_id uuid primary key,
    address text not null unique,
    public_key text not null,
    private_key_encrypted text not null,
    created_at bigint not null,
    block_index integer not null default -1
);

comment on table public.wallets is 'Maps Supabase auth users to their on-chain wallet metadata.';
comment on column public.wallets.supabase_user_id is 'Auth user identifier from Supabase.';
comment on column public.wallets.address is 'Deterministic wallet address assigned to the user.';
comment on column public.wallets.public_key is 'Ed25519 public key (Base64 encoded).';
comment on column public.wallets.private_key_encrypted is 'AES-GCM encrypted private key payload.';
comment on column public.wallets.created_at is 'Unix timestamp (seconds) when the wallet was created.';
comment on column public.wallets.block_index is 'Block index where the wallet transaction was finalized (-1 until committed).';

-- Ensure only privileged service role can modify wallet records.
alter table public.wallets enable row level security;

-- Example policy allowing service role to manage records.
-- Drop and recreate policy to keep script idempotent when re-run.
drop policy if exists service_role_manage_wallets on public.wallets;

create policy service_role_manage_wallets on public.wallets
    for all
    using (auth.role() = 'service_role')
    with check (auth.role() = 'service_role');

-- Grant read access to authenticated users if desired (optional).
-- grant select on public.wallets to authenticated;

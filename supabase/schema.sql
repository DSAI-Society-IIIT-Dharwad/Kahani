-- Kahani collaborative storytelling schema
-- Run this in Supabase SQL editor

create extension if not exists "pgcrypto";

-- Enum to track the state of each story project
create type story_status as enum ('draft', 'recruiting', 'in_progress', 'completed', 'archived');

-- Templates authors can use to bootstrap new stories
create table if not exists public.story_templates (
  id uuid primary key default gen_random_uuid(),
  title text not null,
  description text,
  genre text,
  difficulty text,
  estimated_length text,
  prompt jsonb,
  created_by uuid references auth.users (id) on delete set null,
  created_at timestamptz not null default now()
);

-- Collaborative story projects
create table if not exists public.story_projects (
  id uuid primary key default gen_random_uuid(),
  title text not null,
  summary text,
  status story_status not null default 'draft',
  host_id uuid references auth.users (id) on delete set null,
  template_id uuid references public.story_templates (id) on delete set null,
  slots_total integer not null default 4,
  slots_taken integer not null default 1,
  tags text[] default '{}',
  metadata jsonb default '{}',
  is_chain_backed boolean not null default false,
  chain_reference text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

-- Canonical lines for each project
create table if not exists public.story_lines (
  id uuid primary key default gen_random_uuid(),
  story_id uuid not null references public.story_projects (id) on delete cascade,
  author_id uuid references auth.users (id) on delete set null,
  author_handle text,
  content text not null,
  position integer,
  source text default 'user',
  is_final boolean not null default false,
  created_at timestamptz not null default now()
);

create table if not exists public.story_nfts (
  id uuid primary key default gen_random_uuid(),
  token_id text not null unique,
  story_chain_id text,
  story_project_id uuid references public.story_projects (id) on delete set null,
  story_title text,
  minted_by uuid references auth.users (id) on delete set null,
  minted_by_handle text,
  minted_at timestamptz not null default now(),
  metadata jsonb default '{}'
);

create index if not exists story_nfts_minted_at_idx on public.story_nfts (minted_at desc);

-- Participants in each story project
create table if not exists public.story_participants (
  id uuid primary key default gen_random_uuid(),
  story_id uuid not null references public.story_projects (id) on delete cascade,
  user_id uuid not null references auth.users (id) on delete cascade,
  role text default 'writer',
  joined_at timestamptz not null default now(),
  unique (story_id, user_id)
);

-- Maintain updated_at on story_projects
create or replace function public.set_story_updated_at()
returns trigger as $$
begin
  new.updated_at = now();
  return new;
end;
$$ language plpgsql;

create trigger story_projects_set_updated_at
before update on public.story_projects
for each row
execute procedure public.set_story_updated_at();

-- Simple view to summarise participant counts (optional)
create or replace view public.story_project_stats as
select
  sp.id,
  sp.title,
  sp.status,
  sp.is_chain_backed,
  sp.chain_reference,
  sp.created_at,
  sp.updated_at,
  count(p.id) as participant_count
from public.story_projects sp
left join public.story_participants p on p.story_id = sp.id
group by sp.id;



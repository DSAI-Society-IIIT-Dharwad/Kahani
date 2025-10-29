# Kahani Website Guide

## Application Pages

- **Home (`/`)**: Fixed "Kahani" masthead with session-aware login/logout, backed by an auto-advancing story carousel. Visitors can open template spotlights, inspect richer copy in the modal, or jump into the experience without signing in.
- **Authentication (`/auth`)**: Supabase-powered sign-in and sign-up sheet presented over an animated emerald gradient. The page introduces Kahani's collaborative storytelling pitch before exposing email/password entry, password toggles, and status toasts.
- **Dashboard (`/dashboard`)**: Authenticated control center that surfaces collaboration metrics, in-progress rooms, recruiting projects, and seeded templates. Stories pull live data from Supabase, highlight on-chain state, show seat availability, and provide deep links into Story Studio, Browse, and template previews.
- **Story Template Explorer (`/story`)**: Rich single-story view that blends Supabase lines, Kahani API suggestions, and blockchain health checks. Writers can join a room, request AI guidance, submit or edit lines, and monitor canonicalization, chain sync, and download status from one screen.
- **Story Studio (`/story-studio`)**: Modernized studio tailored for Supabase projects. It fetches curated story metadata, maps template-specific themes, lets contributors join with a pen name, streams lines, handles AI suggestions, canonicalizes drafts, and exposes minting state with resilient error messaging.
- **Browse Minted NFTs (`/browse`)**: Public-facing gallery cataloging tokens stored in the `story_nfts` table. When the chain API is configured, the page enriches each card with on-chain metadata, image references, and external links, while noting local mint history and Supabase relationships.

## Image Reference

- `AI_Help.png`: Snapshot of the Story Studio AI assistance drawer showing how writers request, review, and insert generated lines alongside context cards.
- `Browser.png`: Minted NFT gallery from `/browse`, highlighting token badges, chain metadata callouts, and quick links back to Story Studio.
- `Dashboard_1.png`: Hero section of the dashboard with "Command Center" welcome, collaboration metrics, and quick action buttons.
- `Dashboard_2.png`: Dashboard detail view focusing on the Active Adventures grid, status badges, template tags, and entry points to writers' rooms.
- `Game_Story_Studio.png`: Full Story Studio canvas illustrating the live transcript, collaborator sidebar, theme-aware action buttons, and join/score widgets.
- `Game_Story_Template.png`: Template-driven story experience that showcases AI suggestions, chain status blocks, and the sentence submission workflow.
- `Landing_Page.png`: Landing carousel with gradient-backed story cards, navigation dots, and the oversized "Kahani" logotype.
- `Login_SignUp.png`: Authentication page featuring the Kahani logotype, onboarding copy, and stacked login/register form over soft gradients.
- `Story_Summarizer.png`: Modal summarizer that appears after canonicalization, presenting condensed prose with download and copy controls.
- `Story_Template_Based.png`: Template selection interface emphasizing grid cards, genre badges, and preview buttons for launching new collaborations.

"use client"

import supabase from "@/lib/supabaseClient"

export interface MintedNftRecord {
  tokenId: string
  storyChainId?: string | null
  storyProjectId?: string | null
  storyTitle?: string | null
  mintedBy?: string | null
  mintedByHandle?: string | null
  metadata?: Record<string, unknown> | null
  mintedAt?: string | null
}

export const recordMintedNft = async ({
  tokenId,
  storyChainId,
  storyProjectId,
  storyTitle,
  mintedBy,
  mintedByHandle,
  metadata,
  mintedAt,
}: MintedNftRecord) => {
  const payload: Record<string, unknown> = {
    token_id: tokenId,
    story_chain_id: storyChainId ?? null,
    story_project_id: storyProjectId ?? null,
    story_title: storyTitle ?? null,
    minted_by: mintedBy ?? null,
    minted_by_handle: mintedByHandle ?? null,
    metadata: metadata ?? null,
  }

  if (mintedAt) {
    payload.minted_at = mintedAt
  }

  const { error } = await supabase.from("story_nfts").upsert(payload, { onConflict: "token_id" })

  if (error) {
    throw error
  }
}

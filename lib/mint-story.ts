"use client"

import { mintChainStory, type ChainMintPayload, type ChainMintResponse } from "@/lib/blockchain-api"
import { recordMintedNft, type MintedNftRecord } from "@/lib/nft-registry"

export interface MintStoryAndRecordOptions extends Omit<MintedNftRecord, "tokenId"> {
  tokenId?: string
}

export const mintStoryAndRecord = async (
  token: string,
  storyId: string,
  payload: ChainMintPayload,
  options: MintStoryAndRecordOptions = {}
): Promise<ChainMintResponse> => {
  const response = await mintChainStory(token, storyId, payload)

  const derivedTokenId = (() => {
    if (options.tokenId) return options.tokenId
    const nftPayload = response.nft as Record<string, unknown> | undefined
    if (nftPayload && typeof nftPayload.token_id === "string") {
      return nftPayload.token_id
    }
    if (nftPayload && typeof nftPayload.id === "string") {
      return nftPayload.id
    }
    const transactionPayload = response.transaction as Record<string, unknown> | undefined
    if (transactionPayload && typeof transactionPayload.token_id === "string") {
      return transactionPayload.token_id
    }
    return null
  })()

  if (derivedTokenId) {
    await recordMintedNft({
      tokenId: derivedTokenId,
      storyChainId: options.storyChainId,
      storyProjectId: options.storyProjectId,
      storyTitle: options.storyTitle,
      mintedBy: options.mintedBy,
      mintedByHandle: options.mintedByHandle,
      metadata: options.metadata,
      mintedAt: options.mintedAt,
    })
  }

  return response
}

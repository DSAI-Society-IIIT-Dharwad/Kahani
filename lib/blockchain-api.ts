const RAW_CHAIN_BASE_URL = process.env.NEXT_PUBLIC_CHAIN_API_BASE_URL

const CHAIN_API_BASE = RAW_CHAIN_BASE_URL?.replace(/\/$/, "") ?? ""

const jsonHeaders: HeadersInit = {
  "Content-Type": "application/json",
  Accept: "application/json",
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  if (!CHAIN_API_BASE) {
    throw new Error("Chain API base URL is not configured.")
  }

  const url = `${CHAIN_API_BASE}${path.startsWith("/") ? path : `/${path}`}`
  let response: Response

  try {
    response = await fetch(url, {
      ...init,
      headers: {
        ...jsonHeaders,
        ...init.headers,
      },
      cache: "no-store",
    })
  } catch (rawError) {
    const reason = rawError instanceof Error ? rawError.message : "Unknown network error"
    throw new Error(`Chain API request failed: ${reason}`)
  }

  if (!response.ok) {
    const message = await response.text()
    throw new Error(message || response.statusText)
  }

  if (response.status === 204) {
    return undefined as T
  }

  return (await response.json()) as T
}

export interface ChainHealth {
  status: string
  uptime_seconds?: number
  latest_block_height?: number
  latest_tx_height?: number
  consensus_ready?: boolean
  metrics?: Record<string, unknown>
}

export interface ChainStoryContribution {
  contribution_id: string
  author: string
  text: string
  timestamp: string
  transaction_hash?: string
}

export interface ChainStoryResponse {
  story_id: string
  title: string
  summary: string | null
  contributions: ChainStoryContribution[]
  authors: string[]
  nfts: Array<Record<string, unknown>>
}

export interface ChainWalletResponse {
  wallet_id: string
  address: string
  created_at: string
  last_seen_at?: string
  balances?: Record<string, unknown>
}

export interface ChainContributePayload {
  story_id: string
  author: string
  text: string
}

export interface ChainMintPayload {
  story_id: string
  token_metadata: Record<string, unknown>
}

export interface ChainTransaction<T = Record<string, unknown>> {
  transaction: T
}

export interface ChainMintResponse {
  nft: Record<string, unknown>
  transaction: Record<string, unknown>
}

export const isChainApiConfigured = () => Boolean(CHAIN_API_BASE)

export const fetchChainHealth = () => request<ChainHealth>("/api/health")

export const fetchChainStory = (storyId: string) => request<ChainStoryResponse>(`/api/story/${storyId}`)

export const fetchChainWallet = (supabaseUserId: string) => request<ChainWalletResponse>(`/api/wallet/${supabaseUserId}`)

export const fetchChainNft = (tokenId: string) => request<Record<string, unknown>>(`/api/nft/${tokenId}`)

export const fetchChainNftAuthors = (tokenId: string) => request<Record<string, unknown>>(`/api/nft/${tokenId}/authors`)

export const submitChainContribution = (token: string, payload: ChainContributePayload) =>
  request<ChainTransaction>(`/api/story/contribute`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(payload),
  })

export const mintChainStory = (token: string, storyId: string, payload: ChainMintPayload) =>
  request<ChainMintResponse>(`/api/story/${storyId}/mint`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(payload),
  })

export const createChainEventsSocket = () => {
  if (!CHAIN_API_BASE) {
    throw new Error("Chain API base URL is not configured.")
  }

  const url = new URL(`${CHAIN_API_BASE}/api/events`)
  const wsUrl = url.toString().replace(/^http/, "ws")
  return new WebSocket(wsUrl)
}

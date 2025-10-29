const RAW_BASE_URL = process.env.NEXT_PUBLIC_KAHANI_API_BASE_URL

const API_PROXY_BASE = "/api/kahani"

export class KahaniApiError extends Error {
    constructor(message: string, public readonly status?: number) {
        super(message)
        this.name = "KahaniApiError"
    }
}

const request = async <T>(path: string, init: RequestInit = {}): Promise<T> => {
    const response = await fetch(`${API_PROXY_BASE}${path}`, {
        headers: {
            "Content-Type": "application/json",
            Accept: "application/json",
            ...(init.headers || {}),
        },
        cache: "no-store",
        ...init,
    })

    if (!response.ok) {
        const message = await response.text()
        throw new KahaniApiError(message || response.statusText, response.status)
    }

    if (response.status === 204) {
        return undefined as T
    }

    return (await response.json()) as T
}

export interface SuggestStoryRequest {
    user_prompt: string
    user_id?: string
}

export interface SuggestStoryResponse {
    id: number
    suggestion: string
    context_used: Array<Record<string, unknown>>
    context_count: number
    verified: boolean
    embedding_id: string
}

export interface EditStoryRequest {
    llm_proposed: string
    final_text: string
    user_id?: string
}

export interface EditStoryResponse {
    id: number
    suggestion: string
    context_used: Array<Record<string, unknown>>
    context_count: number
    verified: boolean
    embedding_id: string
}

export interface StoryLinePayload {
    id: number
    final_text?: string
    suggestion?: string
    user_id?: string
    verified?: boolean
    created_at?: string
}

export interface ExtractLoreRequest {
    line_ids: number[]
}

export interface LoreExtractionResponse {
    characters: Array<Record<string, unknown>>
    locations: Array<Record<string, unknown>>
    events: Array<Record<string, unknown>>
    items: Array<Record<string, unknown>>
    total_entries: number
}

export type LoreCollectionResponse = Record<string, unknown>

export interface CanonicalizeStoryRequest {
    line_ids: number[]
    title: string
}

export interface CanonicalStoryResponse {
    id: number
    title: string
    full_text: string
    original_lines_count: number
    created_at: string
}

export const suggestStoryLine = (payload: SuggestStoryRequest) =>
    request<SuggestStoryResponse>("/api/story/suggest", {
        method: "POST",
        body: JSON.stringify({ user_id: "default_user", ...payload }),
    })

export const editStoryLine = (payload: EditStoryRequest) =>
    request<EditStoryResponse>("/api/story/edit", {
        method: "POST",
        body: JSON.stringify({ user_id: "default_user", ...payload }),
    })

export const fetchStoryLines = (verifiedOnly = false) =>
    request<StoryLinePayload[]>(`/api/story/lines${verifiedOnly ? "?verified_only=true" : ""}`)

export const verifyStoryLine = (lineId: number, signature = "user_signed") =>
    request<string>(`/api/story/verify/${lineId}`, {
        method: "POST",
        body: JSON.stringify({ line_id: lineId, signature }),
    })

export const retrieveContext = (query: string, top_k = 5, content_type?: string) =>
    request<Record<string, unknown>>("/api/context/retrieve", {
        method: "POST",
        body: JSON.stringify({ query, top_k, content_type }),
    })

export const extractLore = (payload: ExtractLoreRequest) =>
    request<LoreExtractionResponse>("/api/lore/extract", {
        method: "POST",
        body: JSON.stringify(payload),
    })

export const fetchAllLore = () => request<LoreCollectionResponse>("/api/lore/all")

export const canonicalizeStory = (payload: CanonicalizeStoryRequest) =>
    request<CanonicalStoryResponse>("/api/story/canonicalize", {
        method: "POST",
        body: JSON.stringify(payload),
    })

export const fetchCanonicalStory = (storyId: number) =>
    request<CanonicalStoryResponse | string>(`/api/story/canonical/${storyId}`)

export const isKahaniApiConfigured = () => Boolean(RAW_BASE_URL)

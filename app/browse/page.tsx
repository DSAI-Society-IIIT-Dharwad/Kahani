import { cookies } from "next/headers"
import Link from "next/link"
import { createServerComponentClient } from "@supabase/auth-helpers-nextjs"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { fetchChainNft, isChainApiConfigured } from "@/lib/blockchain-api"

export const dynamic = "force-dynamic"

type StoryNftRow = {
  id: string
  token_id: string
  story_chain_id: string | null
  story_project_id: string | null
  story_title: string | null
  minted_by_handle: string | null
  minted_at: string | null
  metadata: Record<string, unknown> | null
}

interface ChainNftDescriptor {
  tokenId: string
  payload: Record<string, unknown> | null
  error?: string
}

const formatDate = (value: string | null) => {
  if (!value) return "Unknown"
  const date = new Date(value)
  if (Number.isNaN(date.valueOf())) return value
  return date.toLocaleString()
}

const resolveTitle = (row: StoryNftRow, chainPayload: Record<string, unknown> | null) => {
  if (row.story_title) return row.story_title
  if (chainPayload && typeof chainPayload === "object" && typeof chainPayload["title"] === "string") {
    return chainPayload["title" as keyof typeof chainPayload] as string
  }
  if (row.metadata && typeof row.metadata === "object" && typeof row.metadata["title"] === "string") {
    return row.metadata["title" as keyof typeof row.metadata] as string
  }
  return `Story token ${row.token_id}`
}

const resolveImage = (chainPayload: Record<string, unknown> | null, metadata: Record<string, unknown> | null) => {
  const tryGet = (source: Record<string, unknown> | null) => {
    if (!source || typeof source !== "object") return null
    const imageValue = source["image"]
    if (typeof imageValue === "string" && imageValue.length > 0) {
      return imageValue
    }
    const imageCid = source["image_cid"]
    if (typeof imageCid === "string" && imageCid.length > 0) {
      return imageCid.startsWith("ipfs://") ? imageCid : `ipfs://${imageCid}`
    }
    return null
  }

  return tryGet(chainPayload) ?? tryGet(metadata)
}

const resolveDescription = (chainPayload: Record<string, unknown> | null, metadata: Record<string, unknown> | null) => {
  const tryGet = (source: Record<string, unknown> | null) => {
    if (!source || typeof source !== "object") return null
    const descriptionValue = source["description"]
    if (typeof descriptionValue === "string" && descriptionValue.length > 0) {
      return descriptionValue
    }
    return null
  }

  return tryGet(chainPayload) ?? tryGet(metadata) ?? "No description available yet."
}

export default async function BrowsePage() {
  const chainConfigured = isChainApiConfigured()
  const cookieStore = await cookies()
  const supabase = createServerComponentClient({
    cookies: () => cookieStore as unknown as ReturnType<typeof cookies>,
  })

  const { data: mintedRows, error } = await supabase
    .from("story_nfts")
    .select("id, token_id, story_chain_id, story_project_id, story_title, minted_by_handle, minted_at, metadata")
    .order("minted_at", { ascending: false })
    .limit(36)

  if (error) {
    console.error("Unable to load minted NFTs", error)
  }

  const rows: StoryNftRow[] = mintedRows ?? []

  const chainDescriptors: Record<string, ChainNftDescriptor> = {}

  if (chainConfigured && rows.length) {
    const chainResults = await Promise.all(
      rows.map(async (row) => {
        try {
          const payload = await fetchChainNft(row.token_id)
          return { tokenId: row.token_id, payload } satisfies ChainNftDescriptor
        } catch (fetchError) {
          const errorMessage = fetchError instanceof Error ? fetchError.message : "Unable to load chain NFT"
          return { tokenId: row.token_id, payload: null, error: errorMessage } satisfies ChainNftDescriptor
        }
      })
    )

    for (const descriptor of chainResults) {
      chainDescriptors[descriptor.tokenId] = descriptor
    }
  }

  return (
    <main className="min-h-screen bg-gradient-to-br from-purple-50 via-white to-emerald-50">
      <div className="mx-auto flex max-w-6xl flex-col gap-12 px-6 py-24">
        <section className="rounded-3xl bg-white/80 p-10 shadow-xl ring-1 ring-purple-100/70 backdrop-blur-lg">
          <div className="flex flex-col gap-4">
            <div>
              <Badge variant="secondary" className="rounded-full">On-chain gallery</Badge>
              <h1 className="mt-3 text-4xl font-semibold text-slate-900">Browse Minted Kahani Stories</h1>
              <p className="mt-3 max-w-2xl text-base text-slate-600">
                Every NFT listed here was recorded on the Kahani blockchain service. When you mint a story, we capture the
                token id so you and other writers can revisit the on-chain edition.
              </p>
            </div>
            <div className="flex flex-wrap gap-3 text-sm text-slate-500">
              <span>{rows.length} token{rows.length === 1 ? "" : "s"} tracked locally</span>
              <span>
                {chainConfigured
                  ? "Blockchain API configured"
                  : "Set NEXT_PUBLIC_CHAIN_API_BASE_URL to resolve chain metadata"}
              </span>
            </div>
            <div className="flex flex-wrap gap-3">
              <Button asChild variant="outline" className="rounded-full border-emerald-300 text-emerald-700">
                <Link href="/story-studio">Create a new story</Link>
              </Button>
              <Button asChild variant="ghost" className="rounded-full text-purple-700">
                <Link href="/dashboard">Back to dashboard</Link>
              </Button>
            </div>
          </div>
        </section>

        <section className="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
          {rows.length === 0 ? (
            <Card className="md:col-span-2 xl:col-span-3 text-center">
              <CardHeader>
                <CardTitle>No minted NFTs recorded</CardTitle>
                <CardDescription>
                  Mint a story from the studio to populate this gallery. We automatically keep track of each minted token.
                </CardDescription>
              </CardHeader>
            </Card>
          ) : (
            rows.map((row) => {
              const descriptor = chainDescriptors[row.token_id]
              const chainPayload = descriptor?.payload ?? null
              const nftTitle = resolveTitle(row, chainPayload)
              const imageUrl = resolveImage(chainPayload, row.metadata)
              const description = resolveDescription(chainPayload, row.metadata)

              return (
                <Card key={row.id} className="border-purple-100/70 bg-white/80">
                  <CardHeader>
                    <div className="flex items-center justify-between">
                      <Badge variant="outline">Token #{row.token_id}</Badge>
                      <span className="text-xs text-slate-500">Minted {formatDate(row.minted_at)}</span>
                    </div>
                    <CardTitle className="text-lg text-slate-900">{nftTitle}</CardTitle>
                    <CardDescription>{description}</CardDescription>
                  </CardHeader>
                  <CardContent className="flex flex-col gap-3 text-sm text-slate-600">
                    {row.minted_by_handle && (
                      <span>Minted by {row.minted_by_handle}</span>
                    )}
                    {row.story_chain_id && (
                      <span>Story chain id: {row.story_chain_id}</span>
                    )}
                    {row.story_project_id && (
                      <span>Linked Supabase project: {row.story_project_id}</span>
                    )}
                    {imageUrl && (
                      <div className="overflow-hidden rounded-xl border border-slate-200 bg-slate-100">
                        <div className="flex items-center justify-center bg-slate-200/60 px-3 py-2 text-xs text-slate-600">
                          Image reference
                        </div>
                        <div className="px-4 py-3 text-xs break-all text-slate-500">{imageUrl}</div>
                      </div>
                    )}
                    {descriptor?.error && (
                      <p className="text-xs text-red-500">{descriptor.error}</p>
                    )}
                  </CardContent>
                  <CardFooter className="flex items-center justify-between">
                    {row.story_chain_id ? (
                      <Button asChild variant="outline" className="rounded-full">
                        <Link href={`/story-studio?id=${row.story_chain_id}`}>
                          View story
                        </Link>
                      </Button>
                    ) : (
                      <div className="text-xs text-slate-400">No linked studio project</div>
                    )}
                    {chainPayload && typeof chainPayload === "object" && chainPayload["external_url"] ? (
                      <Button asChild className="rounded-full">
                        <Link href={String(chainPayload["external_url"])} target="_blank" rel="noopener noreferrer">
                          View metadata
                        </Link>
                      </Button>
                    ) : null}
                  </CardFooter>
                </Card>
              )
            })
          )}
        </section>
      </div>
    </main>
  )
}

import { NextRequest, NextResponse } from "next/server"

const RAW_BASE_URL = process.env.NEXT_PUBLIC_KAHANI_API_BASE_URL

const ensureBaseUrl = () => {
  if (!RAW_BASE_URL) {
    throw new Error("Missing NEXT_PUBLIC_KAHANI_API_BASE_URL. Add it to your .env.local file.")
  }

  return RAW_BASE_URL.replace(/\/$/, "")
}

const buildTargetUrl = (segments: string[], search: string) => {
  const base = ensureBaseUrl()
  const path = segments.join("/")
  const query = search ? search : ""
  return `${base}/${path}${query}`
}

const forwardRequest = async (request: NextRequest, targetSegments: string[]) => {
  if (request.method === "OPTIONS") {
    return new NextResponse(null, { status: 204 })
  }

  const targetUrl = buildTargetUrl(targetSegments, request.nextUrl.search)

  const headers = new Headers(request.headers)
  headers.set("accept", "application/json")
  headers.set("ngrok-skip-browser-warning", "true")
  headers.delete("host")
  headers.delete("connection")
  headers.delete("content-length")
  headers.delete("accept-encoding")

  const init: RequestInit = {
    method: request.method,
    headers,
    cache: "no-store",
  }

  if (!["GET", "HEAD"].includes(request.method)) {
    const bodyText = await request.text()
    init.body = bodyText
  }

  const response = await fetch(targetUrl, init)

  const responseHeaders = new Headers(response.headers)
  responseHeaders.delete("content-length")
  responseHeaders.delete("transfer-encoding")

  return new NextResponse(response.body, {
    status: response.status,
    statusText: response.statusText,
    headers: responseHeaders,
  })
}

const createHandler = () => async (
  request: NextRequest,
  context: { params: Promise<{ kahaniPath?: string[] }> },
) => {
  const { kahaniPath = [] } = await context.params
  const segments = kahaniPath
  try {
    return await forwardRequest(request, segments)
  } catch (error) {
    const message = error instanceof Error ? error.message : "Unexpected proxy error"
    return NextResponse.json({ error: message }, { status: 500 })
  }
}

export const GET = createHandler()
export const POST = createHandler()
export const PUT = createHandler()
export const PATCH = createHandler()
export const DELETE = createHandler()
export const OPTIONS = createHandler()

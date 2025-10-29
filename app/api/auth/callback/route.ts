import { NextResponse } from "next/server"
import { cookies } from "next/headers"
import { createRouteHandlerClient } from "@supabase/auth-helpers-nextjs"

export const dynamic = "force-dynamic"

export async function POST(request: Request) {
  const cookieStore = await cookies()
  const supabase = createRouteHandlerClient({
    cookies: () => cookieStore as unknown as ReturnType<typeof cookies>,
  })
  const { event, session } = await request.json()

  if ((event === "SIGNED_IN" || event === "TOKEN_REFRESHED") && session) {
    await supabase.auth.setSession(session)
  }

  if (event === "SIGNED_OUT") {
    await supabase.auth.signOut()
  }

  return NextResponse.json({ success: true })
}

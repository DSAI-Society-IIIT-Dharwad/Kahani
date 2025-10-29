"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import { usePathname, useRouter } from "next/navigation"
import type { Session } from "@supabase/supabase-js"

import { Button } from "@/components/ui/button"
import supabase from "@/lib/supabaseClient"

export default function TopBar() {
  const pathname = usePathname()
  const router = useRouter()
  const [session, setSession] = useState<Session | null>(null)

  useEffect(() => {
    let mounted = true

    const syncSession = async () => {
      const { data } = await supabase.auth.getSession()
      if (mounted) {
        setSession(data.session ?? null)
      }
    }

    syncSession()

    const {
      data: { subscription },
    } = supabase.auth.onAuthStateChange((_event, newSession) => {
      if (mounted) {
        setSession(newSession)
      }
    })

    return () => {
      mounted = false
      subscription.unsubscribe()
    }
  }, [])

  if (pathname === "/") {
    return null
  }

  const handleLogout = async () => {
    await supabase.auth.signOut()
    setSession(null)
    router.push("/")
  }

  return (
    <div className="sticky top-0 z-50 w-full bg-white/80 backdrop-blur-md border-b border-emerald-100 shadow-sm">
      <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-3">
        <div className="flex items-center gap-6">
          <Link href="/" className="text-lg font-semibold text-emerald-700 hover:text-emerald-900 transition-colors">
            Kahani Dashboard
          </Link>
          <nav className="hidden items-center gap-4 text-sm text-slate-600 md:flex">
            <Link href="/dashboard" className={pathname === "/dashboard" ? "text-emerald-700 font-medium" : "hover:text-emerald-600"}>
              Dashboard
            </Link>
            <Link href="/story-studio" className={pathname.startsWith("/story-studio") ? "text-emerald-700 font-medium" : "hover:text-emerald-600"}>
              Story Studio
            </Link>
            <Link href="/browse" className={pathname === "/browse" ? "text-emerald-700 font-medium" : "hover:text-emerald-600"}>
              Browse NFTs
            </Link>
          </nav>
        </div>
        {session ? (
          <div className="flex items-center gap-3">
            <span className="hidden text-sm text-slate-600 sm:block">
              {session.user.email ?? "Logged in"}
            </span>
            <Button variant="secondary" className="rounded-full" onClick={handleLogout}>
              Log out
            </Button>
          </div>
        ) : (
          <Button asChild className="rounded-full">
            <Link href="/auth">Login</Link>
          </Button>
        )}
      </div>
    </div>
  )
}

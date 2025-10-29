"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import { useRouter } from "next/navigation"
import type { Session } from "@supabase/supabase-js"
import { Button } from "@/components/ui/button"
import supabase from "@/lib/supabaseClient"

export default function Header() {
  const router = useRouter()
  const [session, setSession] = useState<Session | null>(null)

  useEffect(() => {
    let isMounted = true

    supabase.auth.getSession().then(({ data }) => {
      if (isMounted) {
        setSession(data.session ?? null)
      }
    })

    const {
      data: { subscription },
    } = supabase.auth.onAuthStateChange((_event, newSession) => {
      if (isMounted) {
        setSession(newSession)
      }
    })

    return () => {
      isMounted = false
      subscription.unsubscribe()
    }
  }, [])

  const handleSignOut = async () => {
    await supabase.auth.signOut()
    setSession(null)
    router.push("/")
  }

  return (
    <>
      <header className="fixed top-0 left-0 right-0 z-50 px-6 py-0 flex items-start justify-center bg-transparent">
        <div className="flex justify-center gap-4 w-full mt-0 ">
          <div className="text-center">
            <div
              className="text-[150px] font-bold mb-[50px] drop-shadow-lg"
              style={{
                fontFamily: 'var(--font-titan-one), sans-serif',
                WebkitTextStroke: '5px black',
                color: 'transparent'
              }}
            >
              Kahani
            </div>
          </div>
        </div>
        {session ? (
          <div className="absolute top-10 right-6 flex items-center gap-3">
            <span className="rounded-full bg-black/50 text-white/90 px-4 py-2 text-sm font-medium backdrop-blur-md border border-white/30 shadow-lg">
              Logged in as {session.user.email}
            </span>
            <Button
              onClick={handleSignOut}
              className="rounded-full bg-white/80 hover:bg-white text-slate-900 font-semibold text-lg px-6 py-3 backdrop-blur-lg border border-white/80 shadow-xl transition-all hover:scale-105"
            >
              Log out
            </Button>
          </div>
        ) : (
          <Button
            asChild
            className="absolute top-10 right-6 rounded-full bg-white/70 hover:bg-white/90 text-slate-900 font-semibold text-lg px-10 py-3 backdrop-blur-lg border border-white/80 shadow-xl transition-all hover:scale-105"
          >
            <Link href="/auth">Login</Link>
          </Button>
        )}
      </header>
    </>
  )
}

"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import LoginModal from "@/components/login-modal"
import Link from "next/dist/client/link"

export default function Header() {
  const [showLogin, setShowLogin] = useState(false)

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
        <Button
          onClick={() => setShowLogin(true)}
          className="absolute right-6 rounded-full bg-white/20 hover:bg-white/30 text-white font-medium px-8 py-2 backdrop-blur-sm border border-white/30 transition-all hover:scale-105"
        >
          <Link href="/auth">Login</Link>
        </Button>
      </header>
    </>
  )
}

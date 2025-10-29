import type React from "react"
import type { Metadata } from "next"
import { Poppins } from "next/font/google"
import { Shadows_Into_Light } from "next/font/google"
import { Titan_One } from "next/font/google"
import { Analytics } from "@vercel/analytics/next"
import TopBar from "@/components/top-bar"
import "./globals.css"

const poppins = Poppins({
  subsets: ["latin"],
  weight: ["400", "500", "600", "700"],
  variable: "--font-poppins",
})

const shadowsIntoLight = Shadows_Into_Light({
  subsets: ["latin"],
  weight: ["400"],
  variable: "--font-shadows-into-light",
})

const titanOne = Titan_One({
  subsets: ["latin"],
  weight: ["400"],
  variable: "--font-titan-one",
})

export const metadata: Metadata = {
  title: "Kahani - Stories that Inspire",
  description: "Discover magical stories that bring joy and wonder",
  generator: "v0.app",
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en">
      <body className={`${poppins.variable} ${shadowsIntoLight.variable} ${titanOne.variable} font-sans antialiased`}>
        <TopBar />
        {children}
        <Analytics />
      </body>
    </html>
  )
}

"use client"

import type React from "react"

import { cn } from "@/lib/utils"
import Image from "next/image"
import { useState } from "react"

interface Story {
  id: number
  title: string
  description: string
  color: string
  bgColor: string
  image: string
}

interface StoryCardProps {
  story: Story
  position: number
  isActive: boolean
  onClick: () => void
}

export default function StoryCard({ story, position, isActive, onClick }: StoryCardProps) {
  const [rotateX, setRotateX] = useState(0)
  const [rotateY, setRotateY] = useState(0)

  // Calculate transform based on position for 5 visible cards with continuous rotation
  const getTransform = () => {
    const baseTranslate = position * 240 // Adjusted spacing for 5 cards
    const scale = isActive ? 1.05 : 0.85 - Math.abs(position) * 0.08 // Progressive scaling
    const zIndex = 50 - Math.abs(position) * 10
    const opacity = 1 - Math.abs(position) * 0.2 // More opacity reduction for far cards

    // Add rotation for cards on the sides
    const baseRotateY = position * 8 // More rotation for depth

    return {
      transform: `translateX(${baseTranslate}px) scale(${Math.max(scale, 0.65)}) rotateX(${rotateX}deg) rotateY(${rotateY + baseRotateY}deg) translateZ(${isActive ? 50 : 0}px)`,
      zIndex,
      opacity: Math.max(opacity, 0.5),
    }
  }

  const handleMouseMove = (e: React.MouseEvent<HTMLDivElement>) => {
    if (!isActive) return

    const card = e.currentTarget
    const rect = card.getBoundingClientRect()
    const x = e.clientX - rect.left
    const y = e.clientY - rect.top

    const centerX = rect.width / 2
    const centerY = rect.height / 2

    const rotateXValue = ((y - centerY) / centerY) * -10
    const rotateYValue = ((x - centerX) / centerX) * 10

    setRotateX(rotateXValue)
    setRotateY(rotateYValue)
  }

  const handleMouseLeave = () => {
    setRotateX(0)
    setRotateY(0)
  }

  const style = getTransform()

  return (
    <div
      className={cn(
        "absolute w-80 h-[480px] rounded-3xl overflow-hidden transition-all duration-700 cursor-pointer animate-spin-360",
        "hover:shadow-2xl",
        isActive && "ring-2 ring-white/60 shadow-2xl",
        !isActive && "hover:scale-105"
      )}
      style={{ ...style, perspective: "1000px", transformStyle: "preserve-3d" }}
      onClick={onClick}
      onMouseMove={handleMouseMove}
      onMouseLeave={handleMouseLeave}
    >
      {/* Enhanced animated background border - full opacity for active card */}
      <div className={`absolute inset-0 rounded-3xl bg-gradient-to-r ${story.color} p-[2px] ${isActive ? 'opacity-100' : 'opacity-60'}`}>
        <div className={`w-full h-full ${isActive ? 'bg-black/5' : 'bg-black/2'} backdrop-blur-xs rounded-[1.4rem] overflow-hidden relative`}>
          <div className="relative w-full h-full">
            {/* Enhanced background layers - full opacity for active card */}
            <div className={`absolute inset-0 bg-gradient-to-br ${story.color} ${isActive ? 'opacity-30' : 'opacity-15'}`} />
            <div className={`absolute inset-0 bg-gradient-to-t ${isActive ? 'from-black/60 via-black/10' : 'from-black/30 via-black/5'} to-transparent`} />
            <div className="absolute inset-0 bg-gradient-radial from-white/5 via-transparent to-black/10" />

            {/* Animated mesh gradient overlay - full opacity for active card */}
            <div className={`absolute inset-0 ${isActive ? 'opacity-30' : 'opacity-15'} mix-blend-overlay bg-gradient-conic from-transparent via-white/10 to-transparent animate-spin-slow`} />

            {/* Image with better integration - full opacity for active card */}
            <Image
              src={story.image || "/placeholder.svg"}
              alt={story.title}
              fill
              className={`object-cover mix-blend-soft-light ${isActive ? 'opacity-60' : 'opacity-30'}`}
            />

            {/* Static glowing particles effect - reduced opacity */}
            <div className="absolute inset-0 opacity-10">
              <div className="absolute top-1/4 left-1/4 w-2 h-2 bg-white rounded-full animate-float" />
              <div className="absolute top-3/4 right-1/3 w-1 h-1 bg-white rounded-full animate-float delay-1000" />
              <div className="absolute top-1/2 right-1/4 w-1.5 h-1.5 bg-white rounded-full animate-float delay-2000" />
            </div>

            {/* Text content ON the card, not behind it */}
            <div className="absolute inset-0 flex flex-col justify-between p-6 z-10">
              {/* Top section with chapter number */}
              <div className="flex justify-between items-start">
                <div className={`px-3 py-1 rounded-full bg-gradient-to-r ${story.color} text-white text-xs font-bold uppercase tracking-wider shadow-lg backdrop-blur-sm`}>
                  Chapter {story.id}
                </div>
                {isActive && (
                  <div className="flex gap-1">
                    <div className="w-2 h-2 bg-white rounded-full animate-bounce delay-0" />
                    <div className="w-2 h-2 bg-white rounded-full animate-bounce delay-100" />
                    <div className="w-2 h-2 bg-white rounded-full animate-bounce delay-200" />
                  </div>
                )}
              </div>

              {/* Bottom section with title - reduced background opacity for main background visibility */}
              <div className="space-y-3">
                <h3 className="text-2xl font-bold text-white drop-shadow-2xl leading-tight bg-black/25 backdrop-blur-sm rounded-2xl px-4 py-3 border border-white/15">
                  {story.title}
                </h3>
                {isActive && (
                  <p className="text-white/90 text-sm leading-relaxed bg-black/30 backdrop-blur-sm rounded-xl px-4 py-3 border border-white/20 shadow-lg">
                    {story.description}
                  </p>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

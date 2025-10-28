"use client"

import { useState, useEffect } from "react"
import StoryCard from "@/components/story-card"
import AnimatedBackground from "@/components/animated-background"
import StoryDetailModal from "@/components/story-detail-modal"

const stories = [
  {
    id: 1,
    title: "The Magical Forest",
    description: "A tale of wonder and discovery in an enchanted woodland",
    color: "from-emerald-400 via-emerald-500 to-teal-600",
    bgColor: "bg-gradient-to-br from-emerald-50 via-emerald-100 to-teal-200",
    image: "/magical-forest-illustration.jpg",
  },
  {
    id: 2,
    title: "Ocean Dreams",
    description: "Dive deep into the mysteries beneath the waves",
    color: "from-blue-400 via-cyan-500 to-indigo-600",
    bgColor: "bg-gradient-to-br from-blue-50 via-cyan-100 to-indigo-200",
    image: "/ocean-underwater-illustration.jpg",
  },
  {
    id: 3,
    title: "Mountain Quest",
    description: "An adventure to the peaks where eagles soar",
    color: "from-orange-400 via-red-500 to-pink-600",
    bgColor: "bg-gradient-to-br from-orange-50 via-red-100 to-pink-200",
    image: "/mountain-peak-illustration.jpg",
  },
  {
    id: 4,
    title: "Desert Mirage",
    description: "Secrets hidden in the golden sands of time",
    color: "from-yellow-400 via-orange-500 to-red-600",
    bgColor: "bg-gradient-to-br from-yellow-50 via-orange-100 to-red-200",
    image: "/desert-oasis-illustration.jpg",
  },
  {
    id: 5,
    title: "Starlight Journey",
    description: "Travel through galaxies and cosmic wonders",
    color: "from-indigo-400 via-purple-500 to-pink-600",
    bgColor: "bg-gradient-to-br from-indigo-50 via-purple-100 to-pink-200",
    image: "/space-stars-galaxy-illustration.jpg",
  },
  {
    id: 6,
    title: "Garden of Joy",
    description: "Where flowers bloom and happiness grows",
    color: "from-pink-400 via-rose-500 to-purple-600",
    bgColor: "bg-gradient-to-br from-pink-50 via-rose-100 to-purple-200",
    image: "/flower-garden-illustration.jpg",
  },
  {
    id: 7,
    title: "Crystal Caves",
    description: "Hidden treasures in underground crystal formations",
    color: "from-violet-400 via-purple-500 to-indigo-600",
    bgColor: "bg-gradient-to-br from-violet-50 via-purple-100 to-indigo-200",
    image: "/crystal-cave-illustration.jpg",
  },
  {
    id: 8,
    title: "Sky Castle",
    description: "Floating kingdoms high above the clouds",
    color: "from-sky-400 via-blue-500 to-cyan-600",
    bgColor: "bg-gradient-to-br from-sky-50 via-blue-100 to-cyan-200",
    image: "/sky-castle-illustration.jpg",
  },
]

export default function StoryCarousel() {
  const [activeIndex, setActiveIndex] = useState(0)
  const [isAutoPlaying, setIsAutoPlaying] = useState(true)
  const [selectedStory, setSelectedStory] = useState<typeof stories[0] | null>(null)
  const [showStoryModal, setShowStoryModal] = useState(false)

  // Auto-rotate cards from left to right
  useEffect(() => {
    if (!isAutoPlaying) return

    const interval = setInterval(() => {
      setActiveIndex((prev) => (prev + 1) % stories.length)
    }, 3000) // Change every 3 seconds

    return () => clearInterval(interval)
  }, [isAutoPlaying])

  const handlePrevious = () => {
    setIsAutoPlaying(false)
    setActiveIndex((prev) => (prev === 0 ? stories.length - 1 : prev - 1))
    setTimeout(() => setIsAutoPlaying(true), 5000) // Resume after 5 seconds
  }

  const handleNext = () => {
    setIsAutoPlaying(false)
    setActiveIndex((prev) => (prev + 1) % stories.length)
    setTimeout(() => setIsAutoPlaying(true), 5000) // Resume after 5 seconds
  }

  const getCardPosition = (index: number) => {
    return index - activeIndex
  }

  return (
    <div
      className={`min-h-screen flex items-center justify-center transition-all duration-1000 ${stories[activeIndex].bgColor} pt-60 relative overflow-hidden`}
    >
      <AnimatedBackground activeIndex={activeIndex} colorGradient={stories[activeIndex].color} />

      <div className="relative w-full max-w-6xl h-[600px] flex items-center justify-center px-4">
        {/* Navigation Buttons */}
        <button
          onClick={handlePrevious}
          className="absolute left-4 z-40 w-12 h-12 rounded-full bg-white/20 backdrop-blur-md shadow-lg flex items-center justify-center hover:bg-white/30 transition-all hover:scale-110 border border-white/30"
          aria-label="Previous story"
        >
          <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
          </svg>
        </button>

        <button
          onClick={handleNext}
          className="absolute right-4 z-40 w-12 h-12 rounded-full bg-white/20 backdrop-blur-md shadow-lg flex items-center justify-center hover:bg-white/30 transition-all hover:scale-110 border border-white/30"
          aria-label="Next story"
        >
          <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
          </svg>
        </button>

        {/* Story Cards Container - Expanded to show more cards */}
        <div className="relative w-full h-full flex items-center justify-center overflow-hidden">
          <div className="relative w-[1200px] h-full flex items-center justify-center">
            {stories.map((story, index) => {
              const position = getCardPosition(index)
              const isVisible = Math.abs(position) <= 2 // Show more cards at once
              const isActive = position === 0

              if (!isVisible) return null

              return (
                <StoryCard
                  key={story.id}
                  story={story}
                  position={position}
                  isActive={isActive}
                  onClick={() => {
                    if (isActive) {
                      // If card is active, show story details
                      setSelectedStory(story)
                      setShowStoryModal(true)
                    } else {
                      // If card is not active, make it active
                      setIsAutoPlaying(false)
                      setActiveIndex(index)
                      setTimeout(() => setIsAutoPlaying(true), 5000)
                    }
                  }}
                />
              )
            })}
          </div>
        </div>

        {/* Progress Indicators */}
        <div className="absolute bottom-20 left-1/2 -translate-x-1/2 flex gap-2 z-30">
          {stories.map((_, index) => (
            <button
              key={index}
              onClick={() => {
                setIsAutoPlaying(false)
                setActiveIndex(index)
                setTimeout(() => setIsAutoPlaying(true), 5000)
              }}
              className={`w-3 h-3 rounded-full transition-all duration-300 ${index === activeIndex
                ? 'bg-white scale-125'
                : 'bg-white/50 hover:bg-white/70'
                }`}
              aria-label={`Go to story ${index + 1}`}
            />
          ))}
        </div>
      </div>

      <StoryDetailModal
        story={selectedStory}
        open={showStoryModal}
        onOpenChange={setShowStoryModal}
      />
    </div>
  )
}

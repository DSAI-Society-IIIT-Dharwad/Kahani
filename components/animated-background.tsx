"use client"

import { useEffect, useState } from "react"

interface Circle {
  id: number
  size: number
  x: number
  y: number
  duration: number
  delay: number
}

interface AnimatedBackgroundProps {
  activeIndex: number
  colorGradient: string
}

export default function AnimatedBackground({ activeIndex, colorGradient }: AnimatedBackgroundProps) {
  const [circles, setCircles] = useState<Circle[]>([])

  useEffect(() => {
    // Generate random circles
    const newCircles = Array.from({ length: 8 }, (_, i) => ({
      id: i,
      size: Math.random() * 300 + 100,
      x: Math.random() * 100,
      y: Math.random() * 100,
      duration: Math.random() * 3 + 2,
      delay: Math.random() * 2,
    }))
    setCircles(newCircles)
  }, [])

  return (
    <div className="absolute inset-0 overflow-hidden pointer-events-none">
      {circles.map((circle) => (
        <div
          key={circle.id}
          className={`absolute rounded-full bg-gradient-to-br ${colorGradient} opacity-20 blur-3xl animate-pulse-scale`}
          style={{
            width: `${circle.size}px`,
            height: `${circle.size}px`,
            left: `${circle.x}%`,
            top: `${circle.y}%`,
            animationDuration: `${circle.duration}s`,
            animationDelay: `${circle.delay}s`,
            transform: "translate(-50%, -50%)",
            transition: "all 0.7s ease-in-out",
          }}
        />
      ))}
    </div>
  )
}

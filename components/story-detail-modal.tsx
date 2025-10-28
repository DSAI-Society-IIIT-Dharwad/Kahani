"use client"

import React from "react"
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import Image from "next/image"
import { useRouter } from "next/navigation"

interface Story {
    id: number
    title: string
    description: string
    color: string
    bgColor: string
    image: string
}

interface StoryDetailModalProps {
    story: Story | null
    open: boolean
    onOpenChange: (open: boolean) => void
}

const storyDetails = {
    1: {
        fullDescription: "Deep in the heart of an ancient woodland lies a magical forest where time moves differently. Here, talking animals share wisdom with wandering travelers, and every tree holds a secret from centuries past. The protagonist, a young adventurer named Luna, discovers she has the ability to communicate with the forest spirits and must help them save their home from a creeping darkness that threatens to consume all magic.",
        characters: ["Luna - The brave young adventurer", "Elderoak - The wise tree spirit", "Whisper - A mischievous fairy companion"],
        themes: ["Friendship", "Environmental protection", "Self-discovery", "Magic and wonder"],
        ageGroup: "Ages 8-12",
        readingTime: "15-20 minutes"
    },
    2: {
        fullDescription: "Beneath the endless blue of the ocean lies a world more vibrant and mysterious than any land-dweller could imagine. Captain Marina, a skilled underwater explorer, embarks on a quest to find the legendary Pearl of Dreams, said to grant the deepest wishes of those pure of heart. Along her journey, she encounters bioluminescent creatures, ancient underwater cities, and learns that the greatest treasures aren't always made of gold.",
        characters: ["Captain Marina - Fearless ocean explorer", "Coral - A wise sea turtle", "Neptune - Guardian of the deep seas"],
        themes: ["Adventure", "Ocean conservation", "Courage", "Hidden worlds"],
        ageGroup: "Ages 10-14",
        readingTime: "20-25 minutes"
    },
    3: {
        fullDescription: "High above the clouds, where the air is thin and the view stretches to infinity, lies the challenge of a lifetime. Young mountaineer Alex has always dreamed of reaching the summit of Mount Eternus, the tallest peak in the mystical Skyreach Mountains. But this isn't just any climbing adventure - the mountain is alive with ancient magic, and only those who prove their worth through acts of kindness and bravery can reach the top.",
        characters: ["Alex - Determined young climber", "Storm Eagle - Majestic mountain guardian", "Yuki - A wise mountain guide"],
        themes: ["Perseverance", "Respect for nature", "Personal growth", "Achievement"],
        ageGroup: "Ages 9-13",
        readingTime: "18-22 minutes"
    },
    4: {
        fullDescription: "In the vast expanse of the Whispering Desert, mirages aren't just tricks of light - they're doorways to hidden realms. Desert wanderer Zara possesses the rare gift of seeing through illusions and discovering the truth behind the shimmering heat waves. When her village's oasis begins to dry up, she must venture deep into the desert's heart to find the legendary Spring of Life, facing sandstorm spirits and solving ancient riddles along the way.",
        characters: ["Zara - Desert guide with special sight", "Mirage - A shapeshifting desert spirit", "Cactus King - Guardian of desert wisdom"],
        themes: ["Perseverance", "Desert ecology", "Ancient wisdom", "Resourcefulness"],
        ageGroup: "Ages 8-12",
        readingTime: "16-20 minutes"
    },
    5: {
        fullDescription: "Among the countless stars and swirling galaxies, a young astronaut named Nova discovers that space isn't empty at all - it's filled with cosmic music, dancing asteroids, and friendly alien civilizations. When a mysterious signal from a distant planet reaches Earth, Nova volunteers for humanity's first intergalactic friendship mission. The journey teaches her that despite our differences, kindness and curiosity are universal languages.",
        characters: ["Nova - Young space explorer", "Stardust - An alien friend made of cosmic matter", "Captain Cosmos - Veteran space navigator"],
        themes: ["Space exploration", "Friendship across differences", "Scientific wonder", "Unity"],
        ageGroup: "Ages 10-14",
        readingTime: "20-25 minutes"
    },
    6: {
        fullDescription: "In a magical garden where emotions bloom as flowers, young gardener Joy discovers that her grandmother's secret garden has the power to heal hearts and spread happiness. When sadness begins to wilt the flowers and darkness creeps across the land, Joy must learn to tend not just the plants, but also the feelings of everyone around her. Each flower she nurtures teaches her a lesson about empathy, love, and the importance of emotional wellness.",
        characters: ["Joy - Empathetic young gardener", "Bloom - A talking sunflower", "Grandmother Iris - Keeper of garden secrets"],
        themes: ["Emotional intelligence", "Mental health", "Family bonds", "Nature's healing power"],
        ageGroup: "Ages 7-11",
        readingTime: "14-18 minutes"
    },
    7: {
        fullDescription: "Deep beneath the earth in the Crystal Caves, where gemstones sing ancient melodies and underground rivers flow with liquid starlight, young geologist Ruby makes an incredible discovery. The caves are home to crystal creatures who have been guarding Earth's memories for millions of years. When mining operations threaten their home, Ruby must find a way to protect these magical beings while teaching humans about the incredible world that exists beneath their feet.",
        characters: ["Ruby - Passionate young geologist", "Quartz - Eldest crystal guardian", "Echo - A playful cave sprite"],
        themes: ["Environmental protection", "Hidden wonders", "Scientific discovery", "Coexistence"],
        ageGroup: "Ages 9-13",
        readingTime: "17-21 minutes"
    },
    8: {
        fullDescription: "High above the clouds, where castles float on wisps of air and sky ships sail between floating islands, lives Princess Celeste in the magnificent Sky Castle. When storm clouds begin to gather and threaten to ground all the floating kingdoms, Celeste must embark on a quest through the cloud layers to restore the ancient Wind Crystals. Her journey teaches her that true leadership comes from serving others and that even the smallest acts of kindness can change the world.",
        characters: ["Princess Celeste - Brave sky kingdom heir", "Nimbus - A loyal cloud dragon", "Captain Gale - Sky ship navigator"],
        themes: ["Leadership", "Weather and nature", "Responsibility", "Courage in adversity"],
        ageGroup: "Ages 8-12",
        readingTime: "19-23 minutes"
    }
}

export default function StoryDetailModal({ story, open, onOpenChange }: StoryDetailModalProps) {
    const router = useRouter()

    if (!story) return null

    const details = storyDetails[story.id as keyof typeof storyDetails]

    const handleStartReading = () => {
        onOpenChange(false) // Close the modal
        router.push(`/story?id=${story.id}`) // Navigate to collaborative story page
    }

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle className="text-2xl font-bold text-center">{story.title}</DialogTitle>
                </DialogHeader>

                <div className="space-y-6">
                    {/* Story Image and Basic Info */}
                    <div className="flex flex-col md:flex-row gap-6">
                        <div className="relative w-full md:w-1/3 h-64 rounded-lg overflow-hidden">
                            <Image
                                src={story.image || "/placeholder.svg"}
                                alt={story.title}
                                fill
                                className="object-cover"
                            />
                            <div className={`absolute inset-0 bg-gradient-to-t ${story.color} opacity-30`} />
                        </div>

                        <div className="flex-1 space-y-4">
                            <div className="grid grid-cols-2 gap-4 text-sm">
                                <div>
                                    <span className="font-semibold text-gray-600">Age Group:</span>
                                    <p className="text-gray-800">{details?.ageGroup}</p>
                                </div>
                                <div>
                                    <span className="font-semibold text-gray-600">Reading Time:</span>
                                    <p className="text-gray-800">{details?.readingTime}</p>
                                </div>
                            </div>

                            <div>
                                <h3 className="font-semibold text-gray-600 mb-2">Themes:</h3>
                                <div className="flex flex-wrap gap-2">
                                    {details?.themes.map((theme, index) => (
                                        <span
                                            key={index}
                                            className={`px-3 py-1 rounded-full text-xs font-medium bg-gradient-to-r ${story.color} text-white`}
                                        >
                                            {theme}
                                        </span>
                                    ))}
                                </div>
                            </div>
                        </div>
                    </div>

                    {/* Full Description */}
                    <div>
                        <h3 className="text-lg font-semibold mb-3 text-gray-800">Story Overview</h3>
                        <p className="text-gray-700 leading-relaxed">{details?.fullDescription}</p>
                    </div>

                    {/* Characters */}
                    <div>
                        <h3 className="text-lg font-semibold mb-3 text-gray-800">Main Characters</h3>
                        <ul className="space-y-2">
                            {details?.characters.map((character, index) => (
                                <li key={index} className="flex items-start gap-2 text-gray-700">
                                    <span className="text-blue-500 mt-1">•</span>
                                    {character}
                                </li>
                            ))}
                        </ul>
                    </div>

                    {/* Action Buttons */}
                    <div className="flex gap-3 pt-4 border-t">
                        <Button
                            onClick={handleStartReading}
                            className={`flex-1 bg-gradient-to-r ${story.color} text-white hover:opacity-90`}
                        >
                            Start Reading
                        </Button>
                        <Button variant="outline" className="flex-1">
                            Add to Library
                        </Button>
                        <Button variant="outline">
                            Share Story
                        </Button>
                    </div>
                </div>
            </DialogContent>
        </Dialog>
    )
}
"use client"

import { useState, useEffect, useRef } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Textarea } from "@/components/ui/textarea"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { ArrowLeft, Clock, BookOpen, Send, Zap } from "lucide-react"

interface StorySentence {
    id: number
    text: string
    author: string
    timestamp: string
    color: string
    likes?: number
    isLiked?: boolean
}

interface ActivePlayer {
    id: string
    name: string
    color: string
    isTyping: boolean
    lastActive: string
    score: number
    avatar: string
}



interface Story {
    id: number
    title: string
    description: string
    color: string
    bgColor: string
    image: string
}

// Mock data for collaborative sentences
const mockSentences: Record<number, StorySentence[]> = {
    1: [
        {
            id: 1,
            text: "Luna stepped into the magical forest, her heart pounding with excitement as ancient trees whispered secrets in languages she had never heard before.",
            author: "Emma_Writer",
            timestamp: "2 hours ago",
            color: "from-emerald-400 to-emerald-600"
        },
        {
            id: 2,
            text: "The path ahead shimmered with golden light, and she noticed that her footsteps left tiny flowers blooming in the moss behind her.",
            author: "StoryMaster23",
            timestamp: "1 hour ago",
            color: "from-green-400 to-green-600"
        },
        {
            id: 3,
            text: "Suddenly, a wise old owl perched on a branch above called out, 'Welcome, Luna! The forest has been waiting for someone with your pure heart to help us.'",
            author: "MagicTeller",
            timestamp: "45 minutes ago",
            color: "from-teal-400 to-teal-600"
        }
    ],
    2: [
        {
            id: 1,
            text: "Captain Marina dove deeper into the crystal-clear waters, her breathing apparatus humming softly as schools of luminescent fish danced around her.",
            author: "OceanExplorer",
            timestamp: "3 hours ago",
            color: "from-blue-400 to-blue-600"
        },
        {
            id: 2,
            text: "The ancient coral formations below seemed to pulse with their own inner light, revealing intricate patterns that told stories of civilizations long forgotten.",
            author: "DeepSeaDreamer",
            timestamp: "2 hours ago",
            color: "from-cyan-400 to-cyan-600"
        }
    ],
    // Add more mock data for other stories...
    3: [
        {
            id: 1,
            text: "Alex gripped the rocky ledge, feeling the mountain's ancient power thrumming through the stone as eagles soared majestically overhead.",
            author: "PeakSeeker",
            timestamp: "4 hours ago",
            color: "from-orange-400 to-orange-600"
        }
    ]
}

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
    // Add other stories as needed...
]

export default function CollaborativeStoryGame() {
    const router = useRouter()
    const searchParams = useSearchParams()
    const storyId = parseInt(searchParams.get('id') || '1')

    const [story, setStory] = useState<Story | null>(null)
    const [sentences, setSentences] = useState<StorySentence[]>([])
    const [newSentence, setNewSentence] = useState("")
    const [isSubmitting, setIsSubmitting] = useState(false)
    const [playerName, setPlayerName] = useState("")
    const [hasJoined, setHasJoined] = useState(false)
    const [currentPlayer, setCurrentPlayer] = useState<ActivePlayer | null>(null)

    const messagesEndRef = useRef<HTMLDivElement>(null)

    // Initialize game
    useEffect(() => {
        const foundStory = stories.find(s => s.id === storyId)
        if (foundStory) {
            setStory(foundStory)
        }

        const storySentences = mockSentences[storyId] || []
        setSentences(storySentences)
    }, [storyId])

    // Auto-scroll to bottom
    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
    }, [sentences])

    const handleJoinGame = () => {
        if (!playerName.trim()) return

        const newPlayer: ActivePlayer = {
            id: 'current-user',
            name: playerName,
            color: story?.color || "from-gray-400 to-gray-600",
            isTyping: false,
            lastActive: "Just now",
            score: 0,
            avatar: playerName.slice(0, 2).toUpperCase()
        }

        setCurrentPlayer(newPlayer)
        setHasJoined(true)
    }

    const handleSubmitSentence = async () => {
        if (!newSentence.trim() || isSubmitting || !currentPlayer) return

        setIsSubmitting(true)

        // Simulate API call delay
        await new Promise(resolve => setTimeout(resolve, 800))

        const newSentenceObj: StorySentence = {
            id: sentences.length + 1,
            text: newSentence.trim(),
            author: currentPlayer.name,
            timestamp: "Just now",
            color: currentPlayer.color,
            likes: 0,
            isLiked: false
        }

        setSentences([...sentences, newSentenceObj])
        setNewSentence("")
        setIsSubmitting(false)

        // Update player score
        setCurrentPlayer(prev => prev ? { ...prev, score: prev.score + 10 } : null)
    }

    const handleLikeSentence = (sentenceId: number) => {
        setSentences(prev => prev.map(sentence =>
            sentence.id === sentenceId
                ? {
                    ...sentence,
                    likes: (sentence.likes || 0) + (sentence.isLiked ? -1 : 1),
                    isLiked: !sentence.isLiked
                }
                : sentence
        ))
    }



    if (!story) {
        return <div>Loading...</div>
    }

    // Join Game Screen
    if (!hasJoined) {
        return (
            <div className={`min-h-screen ${story.bgColor} flex items-center justify-center`}>
                <Card className="w-full max-w-md mx-4 bg-white/90 backdrop-blur-sm border-white/30">
                    <CardHeader className="text-center">
                        <CardTitle className="text-2xl font-bold bg-gradient-to-r from-gray-800 to-gray-600 bg-clip-text text-transparent">
                            Join Story Game
                        </CardTitle>
                        <p className="text-gray-600">{story.title}</p>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div className="text-center space-y-2">
                            <div className="flex items-center justify-center gap-2 text-sm text-gray-600">
                                <BookOpen className="w-4 h-4" />
                                <span>Ready to create a story together?</span>
                            </div>
                        </div>

                        <div className="space-y-3">
                            <input
                                type="text"
                                placeholder="Enter your player name..."
                                value={playerName}
                                onChange={(e) => setPlayerName(e.target.value)}
                                className="w-full px-4 py-3 rounded-lg bg-white/70 border border-white/50 focus:outline-none focus:ring-2 focus:ring-blue-500"
                                onKeyPress={(e) => e.key === 'Enter' && handleJoinGame()}
                            />
                            <Button
                                onClick={handleJoinGame}
                                disabled={!playerName.trim()}
                                className={`w-full bg-gradient-to-r ${story.color} text-white hover:opacity-90 py-3`}
                            >
                                <Zap className="w-4 h-4 mr-2" />
                                Join Game
                            </Button>
                        </div>

                        <div className="text-xs text-gray-500 text-center">
                            <p>📖 Add sentences to build the story together!</p>
                            <p>✨ Be creative and keep the story flowing</p>
                            <p>🏆 Earn points for your contributions</p>
                        </div>
                    </CardContent>
                </Card>
            </div>
        )
    }

    return (
        <div className={`min-h-screen ${story.bgColor} transition-all duration-1000`}>
            {/* Game Header */}
            <div className="sticky top-0 z-50 bg-white/90 backdrop-blur-md border-b border-white/20">
                <div className="max-w-4xl mx-auto px-4 py-3">
                    <div className="flex items-center justify-between">
                        <div className="flex items-center gap-3">
                            <Button
                                onClick={() => router.back()}
                                variant="ghost"
                                size="sm"
                                className="flex items-center gap-2"
                            >
                                <ArrowLeft className="w-4 h-4" />
                                Back to Stories
                            </Button>
                            <div>
                                <h1 className="text-lg font-bold text-gray-800">{story.title}</h1>
                                <p className="text-xs text-gray-600">Collaborative Story</p>
                            </div>
                        </div>

                        {/* Player Info */}
                        {currentPlayer && (
                            <div className="flex items-center gap-2">
                                <Avatar className="w-8 h-8">
                                    <AvatarFallback className={`bg-gradient-to-r ${currentPlayer.color} text-white text-xs font-semibold`}>
                                        {currentPlayer.avatar}
                                    </AvatarFallback>
                                </Avatar>
                                <div className="text-sm">
                                    <p className="font-medium text-gray-800">{currentPlayer.name}</p>
                                    <p className="text-xs text-gray-600">{currentPlayer.score} points</p>
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            </div>

            <div className="max-w-4xl mx-auto p-4">
                {/* Story Content */}
                <div className="flex-1 space-y-4">
                    {/* Story Introduction */}
                    <Card className="bg-white/70 backdrop-blur-sm border-white/30">
                        <CardHeader>
                            <CardTitle className="text-xl text-center bg-gradient-to-r from-gray-800 to-gray-600 bg-clip-text text-transparent">
                                {story.title}
                            </CardTitle>
                            <p className="text-center text-gray-600 italic">{story.description}</p>
                        </CardHeader>
                    </Card>

                    {/* Story Sentences */}
                    <div className="space-y-3 max-h-96 overflow-y-auto">
                        {sentences.map((sentence, index) => (
                            <Card key={sentence.id} className="bg-white/60 backdrop-blur-sm border-white/30 hover:bg-white/70 transition-all">
                                <CardContent className="p-4">
                                    <div className="flex items-start gap-3">
                                        <Avatar className="w-8 h-8 flex-shrink-0">
                                            <AvatarFallback className={`bg-gradient-to-r ${sentence.color} text-white text-xs font-semibold`}>
                                                {sentence.author.slice(0, 2).toUpperCase()}
                                            </AvatarFallback>
                                        </Avatar>
                                        <div className="flex-1">
                                            <div className="flex items-center gap-2 mb-1">
                                                <span className="font-medium text-gray-800 text-sm">{sentence.author}</span>
                                                <Badge variant="outline" className="text-xs">
                                                    #{index + 1}
                                                </Badge>
                                                <span className="text-xs text-gray-500">{sentence.timestamp}</span>
                                            </div>
                                            <p className="text-gray-700 leading-relaxed">
                                                {sentence.text}
                                            </p>
                                            <div className="flex items-center gap-2 mt-2">
                                                <Button
                                                    variant="ghost"
                                                    size="sm"
                                                    onClick={() => handleLikeSentence(sentence.id)}
                                                    className={`text-xs ${sentence.isLiked ? 'text-red-500' : 'text-gray-500'}`}
                                                >
                                                    ❤️ {sentence.likes || 0}
                                                </Button>
                                            </div>
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>
                        ))}
                        <div ref={messagesEndRef} />
                    </div>

                    {/* Input Area */}
                    <Card className="bg-white/80 backdrop-blur-sm border-white/30 sticky bottom-4">
                        <CardContent className="p-4">
                            <div className="space-y-3">
                                <div className="flex items-center justify-between">
                                    <h3 className="font-semibold text-gray-800">Continue the Story</h3>
                                    <div className="text-sm text-gray-600">
                                        Add your creative sentence!
                                    </div>
                                </div>
                                <Textarea
                                    placeholder="Write the next sentence in the story..."
                                    value={newSentence}
                                    onChange={(e) => setNewSentence(e.target.value)}
                                    className="min-h-[80px] resize-none bg-white/70 border-white/50"
                                    maxLength={300}
                                    disabled={isSubmitting}
                                />
                                <div className="flex justify-between items-center">
                                    <span className="text-xs text-gray-500">
                                        {newSentence.length}/300 characters
                                    </span>
                                    <Button
                                        onClick={handleSubmitSentence}
                                        disabled={!newSentence.trim() || isSubmitting}
                                        className={`bg-gradient-to-r ${story.color} text-white hover:opacity-90`}
                                    >
                                        {isSubmitting ? (
                                            <>
                                                <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin mr-2" />
                                                Adding...
                                            </>
                                        ) : (
                                            <>
                                                <Send className="w-4 h-4 mr-2" />
                                                Add to Story (+10 points)
                                            </>
                                        )}
                                    </Button>
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                </div>
            </div>
        </div>
    )
}
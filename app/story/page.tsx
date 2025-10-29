"use client"

import { useState, useEffect, useRef } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Textarea } from "@/components/ui/textarea"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import {
    KahaniApiError,
    canonicalizeStory,
    editStoryLine,
    fetchStoryLines,
    suggestStoryLine,
    type CanonicalStoryResponse,
    type StoryLinePayload,
    isKahaniApiConfigured,
} from "@/lib/kahani-api"
import { ArrowLeft, BookOpen, Download, Send, Users, Zap } from "lucide-react"

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
    accentColor: string
    backgroundTone: string
    buttonClass: string
    textAccent: string
    image: string
    authorName: string
    authorTitle: string
    authorImage: string
}

// Mock data for collaborative sentences
const mockSentences: Record<number, StorySentence[]> = {
    1: [
        {
            id: 1,
            text: "Luna stepped into the magical forest, her heart pounding with excitement as ancient trees whispered secrets in languages she had never heard before.",
            author: "Emma_Writer",
            timestamp: "2 hours ago",
            color: "bg-emerald-500"
        },
        {
            id: 2,
            text: "The path ahead shimmered with golden light, and she noticed that her footsteps left tiny flowers blooming in the moss behind her.",
            author: "StoryMaster23",
            timestamp: "1 hour ago",
            color: "bg-green-500"
        },
        {
            id: 3,
            text: "Suddenly, a wise old owl perched on a branch above called out, 'Welcome, Luna! The forest has been waiting for someone with your pure heart to help us.'",
            author: "MagicTeller",
            timestamp: "45 minutes ago",
            color: "bg-teal-500"
        }
    ],
    2: [
        {
            id: 1,
            text: "Captain Marina dove deeper into the crystal-clear waters, her breathing apparatus humming softly as schools of luminescent fish danced around her.",
            author: "OceanExplorer",
            timestamp: "3 hours ago",
            color: "bg-sky-600"
        },
        {
            id: 2,
            text: "The ancient coral formations below seemed to pulse with their own inner light, revealing intricate patterns that told stories of civilizations long forgotten.",
            author: "DeepSeaDreamer",
            timestamp: "2 hours ago",
            color: "bg-cyan-500"
        }
    ],
    // Add more mock data for other stories...
    3: [
        {
            id: 1,
            text: "Alex gripped the rocky ledge, feeling the mountain's ancient power thrumming through the stone as eagles soared majestically overhead.",
            author: "PeakSeeker",
            timestamp: "4 hours ago",
            color: "bg-orange-500"
        }
    ]
}

const staticOnlinePlayers = [
    { name: "StoryWeaver92", role: "Plot Architect" },
    { name: "LunaDreams", role: "Character Artist" },
    { name: "InkSlinger", role: "Scene Builder" },
    { name: "MythMaker", role: "Twist Specialist" }
]

const accentPalette = [
    "bg-emerald-600",
    "bg-sky-600",
    "bg-amber-600",
    "bg-rose-600",
    "bg-indigo-600",
    "bg-teal-600",
    "bg-purple-600",
    "bg-slate-600"
]

const getRandomAccent = () => accentPalette[Math.floor(Math.random() * accentPalette.length)]

const stories = [
    {
        id: 1,
        title: "The Magical Forest",
        description: "A tale of wonder and discovery in an enchanted woodland",
        accentColor: "bg-emerald-600",
        backgroundTone: "bg-emerald-50",
        buttonClass: "bg-emerald-600 hover:bg-emerald-700 focus-visible:ring-emerald-500",
        textAccent: "text-emerald-700",
        image: "/magical-forest-illustration.jpg",
        authorName: "Maya Evergreen",
        authorTitle: "Fantasy Storyteller",
        authorImage: "https://images.unsplash.com/photo-1544723795-3fb6469f5b39?auto=format&fit=crop&w=400&q=80",
    },
    {
        id: 2,
        title: "Ocean Dreams",
        description: "Dive deep into the mysteries beneath the waves",
        accentColor: "bg-sky-600",
        backgroundTone: "bg-sky-50",
        buttonClass: "bg-sky-600 hover:bg-sky-700 focus-visible:ring-sky-500",
        textAccent: "text-sky-700",
        image: "/ocean-underwater-illustration.jpg",
        authorName: "Kai Mariner",
        authorTitle: "Explorer of the Deep",
        authorImage: "https://images.unsplash.com/photo-1524504388940-b1c1722653e1?auto=format&fit=crop&w=400&q=80",
    },
    {
        id: 3,
        title: "Mountain Quest",
        description: "An adventure to the peaks where eagles soar",
        accentColor: "bg-amber-600",
        backgroundTone: "bg-amber-50",
        buttonClass: "bg-amber-600 hover:bg-amber-700 focus-visible:ring-amber-500",
        textAccent: "text-amber-700",
        image: "/mountain-peak-illustration.jpg",
        authorName: "Elias Summit",
        authorTitle: "Keeper of Legends",
        authorImage: "https://images.unsplash.com/photo-1521572267360-ee0c2909d518?auto=format&fit=crop&w=400&q=80",
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
    const [showPlayers, setShowPlayers] = useState(false)
    const [isSyncing, setIsSyncing] = useState(false)
    const [isSuggesting, setIsSuggesting] = useState(false)
    const [apiError, setApiError] = useState<string | null>(null)
    const [suggestedLine, setSuggestedLine] = useState<string | null>(null)
    const [suggestionContext, setSuggestionContext] = useState<Array<Record<string, unknown>>>([])
    const [lastProposedLine, setLastProposedLine] = useState<string | null>(null)
    const [storyLines, setStoryLines] = useState<StoryLinePayload[]>([])
    const [isFinishing, setIsFinishing] = useState(false)
    const [finishError, setFinishError] = useState<string | null>(null)
    const [canonicalStory, setCanonicalStory] = useState<CanonicalStoryResponse | null>(null)
    const [isDownloading, setIsDownloading] = useState(false)
    const [downloadError, setDownloadError] = useState<string | null>(null)
    const apiConfigured = isKahaniApiConfigured()

    const messagesEndRef = useRef<HTMLDivElement>(null)

    const transformStoryLines = (lines: StoryLinePayload[]) => {
        return lines
            .map((line, index) => {
                const text = line.final_text || line.suggestion
                if (!text) return null
                return {
                    id: line.id ?? index + 1,
                    text,
                    author: line.user_id || "Kahani AI",
                    timestamp: line.created_at ? new Date(line.created_at).toLocaleString() : "Just now",
                    color: getRandomAccent(),
                    likes: 0,
                    isLiked: false
                } satisfies StorySentence
            })
            .filter(Boolean) as StorySentence[]
    }

    // Initialize game
    useEffect(() => {
        const foundStory = stories.find(s => s.id === storyId)
        if (foundStory) {
            setStory(foundStory)
        }

        const fallbackSentences = mockSentences[storyId] || []
        setSentences(fallbackSentences)
        setApiError(null)

        if (!apiConfigured) {
            return
        }

        let cancelled = false

        const syncStoryLines = async () => {
            setIsSyncing(true)
            try {
                const lines = await fetchStoryLines()
                if (!cancelled && Array.isArray(lines) && lines.length > 0) {
                    setStoryLines(lines)
                    const mapped = transformStoryLines(lines)
                    if (mapped.length) {
                        setSentences(mapped)
                    }
                }
            } catch (error) {
                if (!cancelled) {
                    const message = error instanceof KahaniApiError ? error.message : "Unable to load story lines from Kahani service."
                    setApiError(message)
                }
            } finally {
                if (!cancelled) {
                    setIsSyncing(false)
                }
            }
        }

        syncStoryLines()

        return () => {
            cancelled = true
        }
    }, [storyId, apiConfigured])

    // Auto-scroll to bottom
    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
    }, [sentences])

    const handleJoinGame = () => {
        if (!playerName.trim()) return

        const newPlayer: ActivePlayer = {
            id: 'current-user',
            name: playerName,
            color: story?.accentColor || "bg-slate-500",
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
        setApiError(null)

        const trimmedText = newSentence.trim()
        const sentenceColor = currentPlayer.color || getRandomAccent()
        const timestampLabel = new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })

        const newSentenceObj: StorySentence = {
            id: sentences.length + 1,
            text: trimmedText,
            author: currentPlayer.name,
            timestamp: timestampLabel,
            color: sentenceColor,
            likes: 0,
            isLiked: false
        }

        setSentences(prev => [...prev, newSentenceObj])
        setNewSentence("")
        setSuggestedLine(null)
        setSuggestionContext([])
        setLastProposedLine(null)
        setCanonicalStory(null)
        setDownloadError(null)

        // Update player score
        setCurrentPlayer(prev => prev ? { ...prev, score: prev.score + 10 } : null)

        try {
            if (apiConfigured) {
                const payload = {
                    llm_proposed: lastProposedLine || suggestedLine || trimmedText,
                    final_text: trimmedText,
                    user_id: currentPlayer.id
                }
                await editStoryLine(payload)
                const updatedLines = await fetchStoryLines()
                setStoryLines(updatedLines)
                const mapped = transformStoryLines(updatedLines)
                if (mapped.length) {
                    setSentences(mapped)
                }
            }
        } catch (error) {
            const message = error instanceof KahaniApiError ? error.message : "We saved the line locally but syncing with Kahani failed."
            setApiError(message)
        } finally {
            setIsSubmitting(false)
        }
    }

    const handleFinishStory = async () => {
        if (!apiConfigured) {
            setFinishError("Kahani backend is not configured.")
            return
        }

        if (!story) {
            setFinishError("Story metadata is unavailable. Try reloading the page.")
            return
        }

        setIsFinishing(true)
        setFinishError(null)
        setDownloadError(null)

        try {
            const latestLines = await fetchStoryLines()
            setStoryLines(latestLines)

            if (!Array.isArray(latestLines) || latestLines.length === 0) {
                throw new Error("No stored story lines available to finalize. Add a few lines first.")
            }

            const lineIds = latestLines.map(line => line.id).filter((id): id is number => typeof id === "number")

            if (!lineIds.length) {
                throw new Error("Unable to determine story line identifiers.")
            }

            const canonicalResult = await canonicalizeStory({ line_ids: lineIds, title: story.title })

            if (typeof canonicalResult === "string") {
                setCanonicalStory({
                    id: Date.now(),
                    title: story.title,
                    full_text: canonicalResult,
                    original_lines_count: lineIds.length,
                    created_at: new Date().toISOString(),
                })
            } else {
                setCanonicalStory(canonicalResult)
            }
        } catch (error) {
            const message = error instanceof KahaniApiError ? error.message : (error instanceof Error ? error.message : "Unable to finish the story right now.")
            setFinishError(message)
        } finally {
            setIsFinishing(false)
        }
    }

    const handleDownloadPdf = async () => {
        if (!canonicalStory) {
            return
        }

        setIsDownloading(true)
        setDownloadError(null)

        try {
            const { jsPDF } = await import("jspdf")
            const doc = new jsPDF()
            const margin = 18
            const pageWidth = doc.internal.pageSize.getWidth()
            const pageHeight = doc.internal.pageSize.getHeight()
            const wrapWidth = pageWidth - margin * 2

            const title = canonicalStory.title || story?.title || "Kahani Story"
            const finalText = canonicalStory.full_text || ""

            doc.setFont("helvetica", "bold")
            doc.setFontSize(18)
            doc.text(title, margin, 24)

            doc.setFont("helvetica", "normal")
            doc.setFontSize(12)

            const lines = doc.splitTextToSize(finalText, wrapWidth) as string[]
            let cursorY = 40

            lines.forEach((line: string) => {
                if (cursorY > pageHeight - margin) {
                    doc.addPage()
                    cursorY = margin
                }
                doc.text(line, margin, cursorY)
                cursorY += 14
            })

            const safeTitle = title
                .toLowerCase()
                .replace(/[^a-z0-9]+/g, "-")
                .replace(/(^-|-$)/g, "") || "kahani-story"

            doc.save(`${safeTitle}.pdf`)
        } catch (error) {
            const message = error instanceof Error ? error.message : "Unable to generate PDF. Please try again."
            setDownloadError(message)
        } finally {
            setIsDownloading(false)
        }
    }

    const handleSuggestLine = async () => {
        if (!apiConfigured || !story) {
            setApiError("Kahani backend is not configured.")
            return
        }

        setIsSuggesting(true)
        setApiError(null)

        try {
            const prompt = newSentence.trim() || `Continue the story "${story.title}"`;
            const response = await suggestStoryLine({
                user_prompt: prompt,
                user_id: currentPlayer?.id || playerName || "default_user"
            })

            setSuggestedLine(response.suggestion)
            setLastProposedLine(response.suggestion)
            setSuggestionContext(response.context_used || [])

            if (!newSentence.trim()) {
                setNewSentence(response.suggestion)
            }
        } catch (error) {
            const message = error instanceof KahaniApiError ? error.message : "Unable to fetch suggestion from Kahani AI."
            setApiError(message)
        } finally {
            setIsSuggesting(false)
        }
    }

    const handleInputChange = (value: string) => {
        if (suggestedLine && value.trim() !== suggestedLine.trim()) {
            setLastProposedLine(null)
        }
        setNewSentence(value)
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
            <div className={`min-h-screen ${story.backgroundTone} flex items-center justify-center`}>
                <Card className="w-full max-w-md mx-4 bg-white border border-neutral-200 shadow-md">
                    <CardHeader className="text-center">
                        <CardTitle className={`text-2xl font-bold ${story.textAccent}`}>
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
                                className="w-full px-4 py-3 rounded-lg bg-white border border-neutral-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
                                onKeyPress={(e) => e.key === 'Enter' && handleJoinGame()}
                            />
                            <Button
                                onClick={handleJoinGame}
                                disabled={!playerName.trim()}
                                className={`w-full ${story.buttonClass} text-white py-3 transition-colors`}
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
        <div className={`min-h-screen ${story.backgroundTone} transition-all duration-1000 flex flex-col overflow-hidden`}>
            {/* Game Header */}
            <div className="sticky top-0 z-50 bg-white border-b border-neutral-200 shadow-sm">
                <div className="w-full px-6 py-3">
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
                                <h1 className={`text-lg font-bold ${story.textAccent}`}>{story.title}</h1>
                                <p className="text-xs text-gray-600">Collaborative Story</p>
                            </div>
                        </div>

                        {/* Player Info */}
                        {currentPlayer && (
                            <div className="flex items-center gap-2">
                                <Avatar className="w-8 h-8">
                                    <AvatarFallback className={`${currentPlayer.color} text-white text-xs font-semibold`}>
                                        {currentPlayer.avatar}
                                    </AvatarFallback>
                                </Avatar>
                                <div className="text-sm">
                                    <p className={`font-medium ${story.textAccent}`}>{currentPlayer.name}</p>
                                    <p className="text-xs text-gray-600">{currentPlayer.score} points</p>
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            </div>

            <main className="flex-1 overflow-hidden">
                <div className="w-full h-full px-6 py-6 lg:flex lg:items-stretch lg:gap-8 overflow-x-hidden overflow-y-auto max-w-7xl mx-auto">
                    {/* Sidebar with author + players */}
                    <aside className="lg:w-1/4 mb-6 lg:mb-0 overflow-hidden lg:h-full">
                        <Card className="bg-white border border-neutral-200 shadow-sm h-full flex flex-col">
                            <CardHeader className="pb-2">
                                <CardTitle className={`text-lg font-semibold ${story.textAccent} flex items-center gap-2`}>
                                    <BookOpen className="w-5 h-5 text-slate-600" />
                                    Story Author & Players
                                </CardTitle>
                            </CardHeader>
                            <CardContent className="flex-1 flex flex-col space-y-4 overflow-hidden">
                                <div className="space-y-3">
                                    <div className="rounded-xl overflow-hidden h-40 bg-gray-200">
                                        <img
                                            src={story.authorImage}
                                            alt={`${story.authorName} portrait`}
                                            className="h-full w-full object-cover"
                                        />
                                    </div>
                                    <div>
                                        <p className={`text-base font-semibold ${story.textAccent}`}>{story.authorName}</p>
                                        <p className="text-sm text-gray-500">{story.authorTitle}</p>
                                    </div>
                                    <p className="text-sm text-gray-600 leading-relaxed">
                                        "Every sentence is a doorway to a new adventure. Keep the magic alive with your words!"
                                    </p>
                                </div>

                                <div className="space-y-3 border-t border-neutral-200 pt-4">
                                    <div className="flex items-center justify-between text-sm font-semibold text-gray-800">
                                        <span className="flex items-center gap-2">
                                            <Users className="w-5 h-5 text-slate-600" />
                                            Online ({staticOnlinePlayers.length})
                                        </span>
                                        <Button
                                            variant="outline"
                                            size="sm"
                                            className="h-7 px-3 text-xs"
                                            onClick={() => setShowPlayers(prev => !prev)}
                                        >
                                            {showPlayers ? "Hide" : "View"}
                                        </Button>
                                    </div>
                                    {showPlayers && (
                                        <div className="space-y-3 overflow-auto pr-1">
                                            {staticOnlinePlayers.map(player => (
                                                <div key={player.name} className="flex items-start gap-3">
                                                    <Avatar className="w-9 h-9">
                                                        <AvatarFallback className="bg-slate-500 text-white text-xs font-semibold">
                                                            {player.name.slice(0, 2).toUpperCase()}
                                                        </AvatarFallback>
                                                    </Avatar>
                                                    <div>
                                                        <p className="text-sm font-medium text-gray-800">{player.name}</p>
                                                        <p className="text-xs text-gray-500">{player.role}</p>
                                                    </div>
                                                </div>
                                            ))}
                                        </div>
                                    )}
                                    <p className="text-xs text-gray-500">
                                        Players update in real time when connected to the live session.
                                    </p>
                                </div>
                            </CardContent>
                        </Card>
                    </aside>

                    {/* Story Content */}
                    <section className="lg:w-3/4 flex-1 h-full flex flex-col space-y-4 overflow-hidden">
                        {/* Story Introduction */}
                        <Card className="bg-white border border-neutral-200 shadow-sm">
                            <CardHeader>
                                <CardTitle className={`text-xl text-center font-semibold ${story.textAccent}`}>
                                    {story.title}
                                </CardTitle>
                                <p className="text-center text-gray-600 italic">{story.description}</p>
                            </CardHeader>
                        </Card>

                        {/* Story Sentences */}
                        <div className="space-y-3 flex-1 overflow-y-auto pr-1">
                            {sentences.map((sentence, index) => (
                                <Card key={sentence.id} className="bg-white border border-neutral-200 hover:bg-neutral-50 transition-colors">
                                    <CardContent className="p-4">
                                        <div className="flex items-start gap-3">
                                            <Avatar className="w-8 h-8 flex-shrink-0">
                                                <AvatarFallback className={`${sentence.color} text-white text-xs font-semibold`}>
                                                    {sentence.author.slice(0, 2).toUpperCase()}
                                                </AvatarFallback>
                                            </Avatar>
                                            <div className="flex-1">
                                                <div className="flex items-center gap-2 mb-1">
                                                    <span className="font-medium text-gray-800 text-sm">{sentence.author}</span>
                                                    <Badge variant="outline" className="text-xs border-neutral-300 text-neutral-600">
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
                        <Card className="bg-white border border-neutral-200 shadow-sm">
                            <CardContent className="p-2">
                                <div className="space-y-2">
                                    <div className="flex items-center justify-between">
                                        <h3 className="font-semibold text-gray-800">Continue the Story</h3>
                                        <div className="flex items-center gap-2 text-xs text-gray-500">
                                            {isSyncing && <span className="flex items-center gap-1 text-blue-500"><span className="h-2 w-2 rounded-full bg-blue-500 animate-ping" />Syncing</span>}
                                            {apiConfigured ? (
                                                <span>Kahani AI ready</span>
                                            ) : (
                                                <span className="text-rose-500">Offline mode</span>
                                            )}
                                        </div>
                                    </div>
                                    {apiConfigured && (
                                        <div className="flex flex-wrap items-center gap-2">
                                            <Button
                                                variant="outline"
                                                size="sm"
                                                onClick={handleSuggestLine}
                                                disabled={isSuggesting}
                                            >
                                                {isSuggesting ? "Summoning..." : "Summon Suggestion"}
                                            </Button>
                                            {suggestedLine && (
                                                <Button
                                                    variant="ghost"
                                                    size="sm"
                                                    onClick={() => handleInputChange(suggestedLine)}
                                                >
                                                    Use Suggestion
                                                </Button>
                                            )}
                                            <Button
                                                size="sm"
                                                onClick={handleFinishStory}
                                                disabled={isFinishing}
                                                className={`${story.buttonClass} text-white transition-colors`}
                                            >
                                                {isFinishing ? "Crafting finale..." : "Finish Story"}
                                            </Button>
                                        </div>
                                    )}
                                    {apiError && (
                                        <div className="rounded-lg border border-rose-200 bg-rose-50 px-3 py-2 text-xs text-rose-600">
                                            {apiError}
                                        </div>
                                    )}
                                    {finishError && (
                                        <div className="rounded-lg border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-700">
                                            {finishError}
                                        </div>
                                    )}
                                    {suggestedLine && (
                                        <div className="rounded-lg border border-purple-200 bg-purple-50 px-3 py-2 text-xs text-purple-700">
                                            <div className="flex items-center justify-between">
                                                <p className="font-semibold">Kahani suggests:</p>
                                                <button
                                                    type="button"
                                                    className="text-[11px] font-semibold uppercase tracking-wide text-purple-500"
                                                    onClick={() => {
                                                        setSuggestedLine(null)
                                                        setSuggestionContext([])
                                                    }}
                                                >
                                                    Dismiss
                                                </button>
                                            </div>
                                            <p className="mt-1 text-sm">{suggestedLine}</p>
                                            {suggestionContext.length > 0 && (
                                                <details className="mt-2 text-[11px] text-purple-500">
                                                    <summary className="cursor-pointer font-semibold">Context used ({suggestionContext.length})</summary>
                                                    <div className="mt-1 space-y-1">
                                                        {suggestionContext.slice(0, 3).map((context, index) => (
                                                            <pre key={index} className="rounded bg-purple-100/70 p-2 text-[10px] text-purple-700 overflow-x-auto">
                                                                {JSON.stringify(context, null, 2)}
                                                            </pre>
                                                        ))}
                                                    </div>
                                                </details>
                                            )}
                                        </div>
                                    )}
                                    <Textarea
                                        placeholder="Write the next sentence in the story..."
                                        value={newSentence}
                                        onChange={(e) => handleInputChange(e.target.value)}
                                        className="min-h-[44px] resize-none bg-white border border-neutral-200"
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
                                            className={`${story.buttonClass} text-white transition-colors`}
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

                        {canonicalStory && (
                            <Card className="bg-white border border-neutral-200 shadow-md">
                                <CardHeader className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                                    <div>
                                        <CardTitle className={`text-lg font-semibold ${story.textAccent}`}>
                                            Canonical Story Finale
                                        </CardTitle>
                                        <p className="text-xs text-gray-500 mt-1">
                                            {canonicalStory.created_at ? `Created ${new Date(canonicalStory.created_at).toLocaleString()}` : "Freshly forged"}
                                            {` • ${canonicalStory.original_lines_count ?? storyLines.length} lines summarized`}
                                        </p>
                                    </div>
                                    <Button
                                        variant="outline"
                                        size="sm"
                                        onClick={handleDownloadPdf}
                                        disabled={isDownloading}
                                        className="gap-1"
                                    >
                                        {isDownloading ? (
                                            "Preparing..."
                                        ) : (
                                            <>
                                                <Download className="h-4 w-4" />
                                                Download PDF
                                            </>
                                        )}
                                    </Button>
                                </CardHeader>
                                <CardContent className="space-y-4">
                                    {downloadError && (
                                        <div className="rounded-lg border border-rose-200 bg-rose-50 px-3 py-2 text-xs text-rose-600">
                                            {downloadError}
                                        </div>
                                    )}
                                    <p className="text-sm leading-relaxed text-gray-700 whitespace-pre-line">
                                        {canonicalStory.full_text}
                                    </p>
                                </CardContent>
                            </Card>
                        )}
                    </section>
                </div>
            </main>
        </div>
    )
}
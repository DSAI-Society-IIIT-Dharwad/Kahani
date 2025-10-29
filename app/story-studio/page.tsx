"use client"

import { useState, useEffect, useRef } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import Link from "next/link"
import { ArrowLeft, BookOpen, Download, Send, Users, Zap } from "lucide-react"

import supabase from "@/lib/supabaseClient"
import { Button } from "@/components/ui/button"
import { Textarea } from "@/components/ui/textarea"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { cn } from "@/lib/utils"
import {
  KahaniApiError,
  canonicalizeStory,
  editStoryLine,
  suggestStoryLine,
  type CanonicalStoryResponse,
  isKahaniApiConfigured,
} from "@/lib/kahani-api"

interface StorySentence {
  id: string
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

interface ThemePalette {
  accentColor: string
  backgroundTone: string
  buttonClass: string
  textAccent: string
  borderClass: string
  avatarBgClass: string
  softBackgroundClass: string
  focusRingClass: string
  headingTextClass: string
}

interface StoryMeta {
  id: string
  title: string
  description: string | null
  accentColor: string
  backgroundTone: string
  buttonClass: string
  textAccent: string
  borderClass: string
  avatarBgClass: string
  softBackgroundClass: string
  focusRingClass: string
  headingTextClass: string
  image?: string | null
  host?: string | null
  templateTitle?: string | null
}

const templateThemeMap: Record<string, ThemePalette> = {
  "The Magical Forest": {
    accentColor: "bg-emerald-600",
    backgroundTone: "bg-emerald-50",
    buttonClass: "bg-emerald-600 text-white hover:bg-emerald-700 focus-visible:ring-emerald-500",
    textAccent: "text-emerald-700",
    borderClass: "border-emerald-200",
    avatarBgClass: "bg-emerald-100",
    softBackgroundClass: "bg-emerald-50/40",
    focusRingClass: "focus-visible:ring-emerald-500",
    headingTextClass: "text-emerald-800",
  },
  "Ocean Dreams": {
    accentColor: "bg-sky-600",
    backgroundTone: "bg-sky-50",
    buttonClass: "bg-sky-600 text-white hover:bg-sky-700 focus-visible:ring-sky-500",
    textAccent: "text-sky-700",
    borderClass: "border-sky-200",
    avatarBgClass: "bg-sky-100",
    softBackgroundClass: "bg-sky-50/40",
    focusRingClass: "focus-visible:ring-sky-500",
    headingTextClass: "text-sky-800",
  },
  "Mountain Quest": {
    accentColor: "bg-amber-600",
    backgroundTone: "bg-amber-50",
    buttonClass: "bg-amber-600 text-white hover:bg-amber-700 focus-visible:ring-amber-500",
    textAccent: "text-amber-700",
    borderClass: "border-amber-200",
    avatarBgClass: "bg-amber-100",
    softBackgroundClass: "bg-amber-50/40",
    focusRingClass: "focus-visible:ring-amber-500",
    headingTextClass: "text-amber-800",
  },
  "Desert Mirage": {
    accentColor: "bg-orange-500",
    backgroundTone: "bg-orange-50",
    buttonClass: "bg-orange-500 text-white hover:bg-orange-600 focus-visible:ring-orange-400",
    textAccent: "text-orange-700",
    borderClass: "border-orange-200",
    avatarBgClass: "bg-orange-100",
    softBackgroundClass: "bg-orange-50/40",
    focusRingClass: "focus-visible:ring-orange-400",
    headingTextClass: "text-orange-800",
  },
  "Starlight Journey": {
    accentColor: "bg-indigo-600",
    backgroundTone: "bg-indigo-50",
    buttonClass: "bg-indigo-600 text-white hover:bg-indigo-700 focus-visible:ring-indigo-500",
    textAccent: "text-indigo-700",
    borderClass: "border-indigo-200",
    avatarBgClass: "bg-indigo-100",
    softBackgroundClass: "bg-indigo-50/40",
    focusRingClass: "focus-visible:ring-indigo-500",
    headingTextClass: "text-indigo-800",
  },
  "Garden of Joy": {
    accentColor: "bg-pink-500",
    backgroundTone: "bg-pink-50",
    buttonClass: "bg-pink-500 text-white hover:bg-pink-600 focus-visible:ring-pink-400",
    textAccent: "text-pink-700",
    borderClass: "border-pink-200",
    avatarBgClass: "bg-pink-100",
    softBackgroundClass: "bg-pink-50/40",
    focusRingClass: "focus-visible:ring-pink-400",
    headingTextClass: "text-pink-800",
  },
  "Crystal Caves": {
    accentColor: "bg-violet-600",
    backgroundTone: "bg-violet-50",
    buttonClass: "bg-violet-600 text-white hover:bg-violet-700 focus-visible:ring-violet-500",
    textAccent: "text-violet-700",
    borderClass: "border-violet-200",
    avatarBgClass: "bg-violet-100",
    softBackgroundClass: "bg-violet-50/40",
    focusRingClass: "focus-visible:ring-violet-500",
    headingTextClass: "text-violet-800",
  },
  "Sky Castle": {
    accentColor: "bg-cyan-600",
    backgroundTone: "bg-cyan-50",
    buttonClass: "bg-cyan-600 text-white hover:bg-cyan-700 focus-visible:ring-cyan-500",
    textAccent: "text-cyan-700",
    borderClass: "border-cyan-200",
    avatarBgClass: "bg-cyan-100",
    softBackgroundClass: "bg-cyan-50/40",
    focusRingClass: "focus-visible:ring-cyan-500",
    headingTextClass: "text-cyan-800",
  },
}

const defaultTheme = templateThemeMap["The Magical Forest"]

const hexThemeLookup: Record<string, ThemePalette> = {
  "#059669": templateThemeMap["The Magical Forest"],
  "#0284c7": templateThemeMap["Ocean Dreams"],
  "#f59e0b": templateThemeMap["Mountain Quest"],
  "#f97316": templateThemeMap["Desert Mirage"],
  "#6366f1": templateThemeMap["Starlight Journey"],
  "#ec4899": templateThemeMap["Garden of Joy"],
  "#a855f7": templateThemeMap["Crystal Caves"],
  "#0ea5e9": templateThemeMap["Sky Castle"],
}

const themePalette: ThemePalette[] = Object.values(templateThemeMap)

const randomAccent = () => themePalette[Math.floor(Math.random() * themePalette.length)]

const hashString = (value: string) => {
  let hash = 0
  for (let index = 0; index < value.length; index += 1) {
    hash = (hash << 5) - hash + value.charCodeAt(index)
    hash |= 0 // Force 32-bit integer
  }
  return hash
}

const deterministicTheme = (seed: string) => {
  const index = Math.abs(hashString(seed)) % themePalette.length
  return themePalette[index]
}

const UUID_PATTERN = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i

const isUuid = (value: string | null | undefined) => (typeof value === "string" ? UUID_PATTERN.test(value) : false)

const FALLBACK_TEMPLATE_TITLES = [
  "The Magical Forest",
  "Ocean Dreams",
  "Mountain Quest",
  "Desert Mirage",
  "Starlight Journey",
  "Garden of Joy",
  "Crystal Caves",
  "Sky Castle",
]

const FALLBACK_DESCRIPTIONS = [
  "Draft a brand new story using Kahani Studio.",
  "Spin the next chapter with your crew.",
  "Invite friends and craft a legendary tale.",
  "Mix magic and lore to shape a new universe.",
]

const staticOnlinePlayers = [
  { name: "StoryWeaver92", role: "Plot Architect" },
  { name: "LunaDreams", role: "Character Artist" },
  { name: "InkSlinger", role: "Scene Builder" },
  { name: "MythMaker", role: "Twist Specialist" },
]

export default function StoryStudioPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const [stories, setStories] = useState<StoryMeta[]>([])
  const [story, setStory] = useState<StoryMeta | null>(null)
  const [sentences, setSentences] = useState<StorySentence[]>([])
  const [newSentence, setNewSentence] = useState("")
  const [suggestedLine, setSuggestedLine] = useState<string | null>(null)
  const [suggestionContext, setSuggestionContext] = useState<Array<Record<string, unknown>>>([])
  const [lastProposedLine, setLastProposedLine] = useState<string | null>(null)
  const [playerName, setPlayerName] = useState("")
  const [hasJoined, setHasJoined] = useState(false)
  const [currentPlayer, setCurrentPlayer] = useState<ActivePlayer | null>(null)
  const [showPlayers, setShowPlayers] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [isSuggesting, setIsSuggesting] = useState(false)
  const [isSyncing, setIsSyncing] = useState(false)
  const [isFinishing, setIsFinishing] = useState(false)
  const [canonicalStory, setCanonicalStory] = useState<CanonicalStoryResponse | null>(null)
  const [isDownloading, setIsDownloading] = useState(false)
  const [apiError, setApiError] = useState<string | null>(null)
  const [finishError, setFinishError] = useState<string | null>(null)
  const [downloadError, setDownloadError] = useState<string | null>(null)
  const apiConfigured = isKahaniApiConfigured()

  const theme: ThemePalette = story
    ? {
        accentColor: story.accentColor,
        backgroundTone: story.backgroundTone,
        buttonClass: story.buttonClass,
        textAccent: story.textAccent,
        borderClass: story.borderClass,
        avatarBgClass: story.avatarBgClass,
        softBackgroundClass: story.softBackgroundClass,
        focusRingClass: story.focusRingClass,
        headingTextClass: story.headingTextClass,
      }
    : defaultTheme

  const gradientFromClass = theme.accentColor.startsWith("bg-")
    ? theme.accentColor.replace("bg-", "from-")
    : "from-emerald-500"
  const gradientViaClass = theme.backgroundTone.startsWith("bg-")
    ? theme.backgroundTone.replace("bg-", "via-")
    : "via-white"

  const storyIdParam = searchParams.get("id")
  const [activeStoryId, setActiveStoryId] = useState<string | null>(storyIdParam)

  const messagesEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const loadStories = async () => {
      const { data, error } = await supabase
        .from("story_projects")
        .select("id,title,summary,status,metadata,tags,host_id,template:story_templates(title)")
        .order("updated_at", { ascending: false })

      if (error) {
        console.error("Unable to load stories", error)
        return
      }

      type SupabaseStoryRow = {
        id: string
        title: string
        summary: string | null
        metadata: Record<string, string> | null
        template: { title: string }[] | { title: string } | null
      }

      const mapped = ((data as SupabaseStoryRow[] | null) ?? []).map((item) => {
        const accentFromMetadata = typeof item.metadata?.accentColor === "string"
          ? item.metadata.accentColor.trim().toLowerCase()
          : undefined

        const templateRelation = Array.isArray(item.template) ? item.template[0] : item.template

        const palette = (() => {
          if (accentFromMetadata && hexThemeLookup[accentFromMetadata]) {
            return hexThemeLookup[accentFromMetadata]
          }

          if (templateRelation?.title && templateThemeMap[templateRelation.title]) {
            return templateThemeMap[templateRelation.title]
          }

          return deterministicTheme(item.id)
        })()

        return {
          id: item.id,
          title: item.title,
          description: item.summary,
          accentColor: palette.accentColor,
          backgroundTone: palette.backgroundTone,
          buttonClass: palette.buttonClass,
          textAccent: palette.textAccent,
          borderClass: palette.borderClass,
          avatarBgClass: palette.avatarBgClass,
          softBackgroundClass: palette.softBackgroundClass,
          focusRingClass: palette.focusRingClass,
          headingTextClass: palette.headingTextClass,
          templateTitle: templateRelation?.title ?? null,
        } satisfies StoryMeta
      })

      if (mapped.length === 0) {
        const fallback = FALLBACK_TEMPLATE_TITLES.map((title, index) => {
          const palette = templateThemeMap[title] ?? deterministicTheme(title)
          const description = FALLBACK_DESCRIPTIONS[index % FALLBACK_DESCRIPTIONS.length]
          return {
            id: `static-${index}`,
            title,
            description,
            ...palette,
          }
        })
        setStories(fallback)
        setActiveStoryId(fallback[0]?.id ?? null)
        setStory(fallback[0] ?? null)
        return
      }

      setStories(mapped)
      const initial = mapped.find((item) => item.id === storyIdParam) ?? mapped[0] ?? null
      setStory(initial)
      setActiveStoryId(initial?.id ?? null)
    }

    loadStories()
  }, [storyIdParam])

  useEffect(() => {
    const loadStoryLines = async () => {
      if (!activeStoryId) return

      if (!isUuid(activeStoryId)) {
        setApiError(null)
        setSentences([])
        setIsSyncing(false)
        return
      }

      setIsSyncing(true)
      setApiError(null)

      const { data, error } = await supabase
        .from("story_lines")
        .select("id, content, author_handle, created_at, source")
        .eq("story_id", activeStoryId)
        .order("created_at", { ascending: true })

      if (error) {
        if (error.code === "PGRST116") {
          setApiError(null)
        } else {
          const message = typeof error.message === "string" && error.message.length
            ? error.message
            : "Unable to load story lines. Is Supabase configured?"
          setApiError(message)
        }
        setIsSyncing(false)
        return
      }

      type SupabaseStoryLine = {
        id: string | number | null
        content: string | null
        author_handle: string | null
        created_at: string | null
      }

      const mappedSentences: StorySentence[] = ((data as SupabaseStoryLine[] | null) ?? []).map((line, index) => ({
        id: typeof line.id === "string" ? line.id : `line-${line.id ?? index}`,
        text: line.content ?? "",
        author: line.author_handle ?? "Storyteller",
        timestamp: line.created_at ? new Date(line.created_at).toLocaleString() : "Just now",
        color: theme.accentColor,
      }))
      setSentences(mappedSentences)
      setIsSyncing(false)
    }

    loadStoryLines()
  }, [activeStoryId])

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" })
  }, [sentences])

  const handleJoinGame = () => {
    if (!playerName.trim()) return
    if (!story) return

  const palette = theme

    const newPlayer: ActivePlayer = {
      id: "current-user",
      name: playerName.trim(),
      color: palette.accentColor,
      isTyping: false,
      lastActive: "Just now",
      score: 0,
      avatar: playerName.slice(0, 2).toUpperCase(),
    }

    setCurrentPlayer(newPlayer)
    setHasJoined(true)
  }

  const handleSubmitSentence = async () => {
    if (!newSentence.trim() || !currentPlayer || !activeStoryId) return

    setIsSubmitting(true)
    setApiError(null)

    const trimmedText = newSentence.trim()
    const timestampLabel = new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })

    const newSentenceObj: StorySentence = {
      id: crypto.randomUUID(),
      text: trimmedText,
      author: currentPlayer.name,
      timestamp: timestampLabel,
      color: currentPlayer.color,
      likes: 0,
      isLiked: false,
    }

    setSentences((prev) => [...prev, newSentenceObj])
    setNewSentence("")
    setSuggestedLine(null)
    setSuggestionContext([])
    setLastProposedLine(null)
    setCanonicalStory(null)
    setDownloadError(null)

    try {
      const { error } = await supabase.from("story_lines").insert({
        story_id: activeStoryId,
        content: trimmedText,
        author_handle: currentPlayer.name,
      })

      if (error) {
        throw error
      }

      if (apiConfigured) {
        await editStoryLine({
          llm_proposed: lastProposedLine || suggestedLine || trimmedText,
          final_text: trimmedText,
          user_id: currentPlayer.id,
        })
      }

      if (isUuid(activeStoryId)) {
        const { data } = await supabase
          .from("story_lines")
          .select("id, content, author_handle, created_at")
          .eq("story_id", activeStoryId)
          .order("created_at", { ascending: true })

        if (data) {
          const mapped = data.map((line) => ({
            id: String(line.id),
            text: line.content,
            author: line.author_handle ?? "Storyteller",
            timestamp: line.created_at ? new Date(line.created_at).toLocaleString() : "Just now",
            color: theme.accentColor,
          }))
          setSentences(mapped)
        }
      }
    } catch (error) {
      console.error(error)
      const message =
        error instanceof KahaniApiError
          ? error.message
          : error instanceof Error
            ? error.message
            : "We saved the line locally but syncing with Supabase failed."
      setApiError(message)
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleSuggestLine = async () => {
    if (!apiConfigured) {
      setApiError("Kahani backend is not configured.")
      return
    }

    if (isSuggesting) return

    setIsSuggesting(true)
    setApiError(null)

    try {
      const context = sentences.slice(-5).map((sentence) => ({ text: sentence.text, author: sentence.author }))
      setSuggestionContext(context)

      const prompt = context.length
        ? context.map((entry) => `${entry.author}: ${entry.text}`).join("\n")
        : "Continue the collaborative story with a fresh line."

      const suggestionResponse = await suggestStoryLine({
        user_prompt: prompt,
      })

      setSuggestedLine(suggestionResponse.suggestion)
      setLastProposedLine(suggestionResponse.suggestion)
    } catch (error) {
      const message = error instanceof KahaniApiError ? error.message : "Unable to fetch a suggestion right now."
      setApiError(message)
    } finally {
      setIsSuggesting(false)
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
      const { data, error } = await supabase
        .from("story_lines")
        .select("id, content")
        .eq("story_id", activeStoryId)
        .order("created_at", { ascending: true })

      if (error || !data?.length) {
        throw new Error("No stored story lines available to finalize. Add a few lines first.")
      }

      const lineIds = data.map((line, index) =>
        typeof line.id === "number"
          ? line.id
          : index + 1,
      )

      const canonicalResult = await canonicalizeStory({
        line_ids: lineIds,
        title: story.title,
      })

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
      const message =
        error instanceof KahaniApiError
          ? error.message
          : error instanceof Error
            ? error.message
            : "Unable to finish the story right now."
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
      const { default: jsPDF } = await import("jspdf")

      const doc = new jsPDF({ orientation: "portrait", unit: "pt", format: "letter" })
      const margin = 48
      const pageWidth = doc.internal.pageSize.getWidth() - margin * 2

      doc.setFont("times", "bold")
      doc.setFontSize(22)
      doc.text(canonicalStory.title ?? story?.title ?? "Kahani Story", margin, margin)

      doc.setFont("times", "normal")
      doc.setFontSize(12)

      const textLines = doc.splitTextToSize(canonicalStory.full_text ?? "", pageWidth)
      doc.text(textLines, margin, margin + 36)

      doc.save(`${(canonicalStory.title ?? story?.title ?? "kahani-story").replace(/\s+/g, "-")}.pdf`)
    } catch (error) {
      console.error(error)
      setDownloadError("Unable to export PDF. Try again shortly.")
    } finally {
      setIsDownloading(false)
    }
  }

  const handleSelectStory = (storyId: string) => {
    const selected = stories.find((s) => s.id === storyId) ?? null
    setStory(selected)
    setActiveStoryId(storyId)
    setSentences([])
    setCanonicalStory(null)
    setHasJoined(false)
    setCurrentPlayer(null)
    setApiError(null)
    router.push(`/story-studio?id=${storyId}`)
  }

  return (
    <main
      className={cn(
        "min-h-screen pb-32 transition-colors duration-500 bg-gradient-to-br to-white",
        theme.backgroundTone,
        gradientFromClass,
        gradientViaClass,
      )}
    >
      <div className="mx-auto max-w-6xl px-6 py-24">
        <div className="mb-8 flex items-center justify-between">
          <Button variant="ghost" className="flex items-center gap-2" asChild>
            <Link href="/dashboard">
              <ArrowLeft className="h-4 w-4" /> Back to dashboard
            </Link>
          </Button>
          <div className="flex items-center gap-2 text-sm text-slate-500">
            <Users className="h-4 w-4" /> {staticOnlinePlayers.length + (hasJoined ? 1 : 0)} storytellers online
          </div>
        </div>

        <div className="grid gap-6 lg:grid-cols-[300px_1fr]">
          <aside className={cn("bg-white/80 backdrop-blur-sm rounded-2xl border shadow-sm p-4 flex flex-col gap-3", theme.borderClass)}>
            <h2 className="text-lg font-semibold text-slate-900">Studio Projects</h2>
            <p className="text-sm text-slate-600">Pick the story arc you want to refine.</p>
            <div className="flex flex-col gap-2 overflow-y-auto max-h-[520px] pr-1">
              {stories.map((item) => {
                const isActive = item.id === activeStoryId
                return (
                  <button
                    key={item.id}
                    onClick={() => handleSelectStory(item.id)}
                    className={cn(
                      "rounded-xl border px-3 py-3 text-left transition",
                      item.borderClass,
                      isActive ? "shadow" : "hover:shadow",
                      isActive ? item.softBackgroundClass : null,
                    )}
                  >
                    <div className="flex items-center justify-between">
                      <span className={cn("text-sm font-semibold", isActive ? item.headingTextClass : "text-slate-900")}>
                        {item.title}
                      </span>
                      {item.templateTitle && <Badge variant="outline">Template</Badge>}
                    </div>
                    <p className="text-xs text-slate-600 mt-1 line-clamp-2">{item.description ?? "Draft summary coming soon."}</p>
                  </button>
                )
              })}
            </div>
            <Button asChild className={cn("mt-2 rounded-full", theme.buttonClass)}>
              <Link href="/dashboard">Manage templates</Link>
            </Button>
          </aside>

          <section className="relative flex flex-col gap-6">
            <Card className={cn("bg-white/80 shadow-lg", theme.borderClass)}>
              <CardHeader className="flex flex-row items-start justify-between gap-4">
                <div>
                  <CardTitle className={cn("text-2xl font-semibold", theme.headingTextClass)}>{story?.title ?? "Select a project"}</CardTitle>
                  <p className="mt-2 max-w-2xl text-sm text-slate-600">
                    {story?.description ?? "Pick a story on the left to start editing in Kahani Studio."}
                  </p>
                </div>
                <div className="flex gap-2">
                  <Button
                    variant="outline"
                    className={cn("rounded-full border", theme.borderClass, theme.textAccent)}
                    onClick={handleSuggestLine}
                    disabled={!apiConfigured || isSuggesting}
                  >
                    <Zap className="mr-2 h-4 w-4" /> Suggest line
                  </Button>
                  <Button
                    className={cn("rounded-full", theme.buttonClass)}
                    onClick={handleFinishStory}
                    disabled={!apiConfigured || isFinishing}
                  >
                    <BookOpen className="mr-2 h-4 w-4" /> Canonicalize
                  </Button>
                </div>
              </CardHeader>
              <CardContent className="flex flex-col gap-6">
                <div className="flex items-center gap-4">
                  <Avatar className={cn("h-12 w-12", theme.avatarBgClass)}>
                    <AvatarFallback className={cn("font-semibold", theme.textAccent)}>
                      {(story?.title ?? "KS").slice(0, 2).toUpperCase()}
                    </AvatarFallback>
                  </Avatar>
                  <div className="flex flex-col">
                    <span className={cn("text-sm font-medium", theme.textAccent)}>
                      {story?.templateTitle ? `Template: ${story.templateTitle}` : "Custom storyline"}
                    </span>
                    <span className="text-xs text-slate-500">
                      Supabase-backed • Editable studio mode
                    </span>
                  </div>
                </div>

                {!hasJoined ? (
                  <div
                    className={cn(
                      "rounded-xl border border-dashed p-4",
                      theme.borderClass,
                      theme.softBackgroundClass,
                    )}
                  >
                    <h3 className={cn("text-lg font-semibold", theme.headingTextClass)}>Join this writers' room</h3>
                    <p className={cn("text-sm", theme.textAccent)}>
                      Share your name with collaborators and start contributing lines to the canonical story.
                    </p>
                    <div className="mt-4 flex flex-col gap-3 sm:flex-row">
                      <Textarea
                        placeholder="Enter your pen name"
                        className={cn("h-10 resize-none", theme.borderClass, theme.focusRingClass)}
                        value={playerName}
                        onChange={(event) => setPlayerName(event.target.value)}
                      />
                      <Button className={cn("rounded-full", theme.buttonClass)} onClick={handleJoinGame}>
                        Join room
                      </Button>
                    </div>
                  </div>
                ) : (
                  <div
                    className={cn(
                      "flex items-center justify-between rounded-xl border p-4",
                      theme.borderClass,
                      theme.softBackgroundClass,
                    )}
                  >
                    <div className="flex items-center gap-3">
                      <Avatar className={cn("h-10 w-10", theme.avatarBgClass)}>
                        <AvatarFallback className={cn("font-bold", theme.textAccent)}>
                          {currentPlayer?.avatar ?? "You"}
                        </AvatarFallback>
                      </Avatar>
                      <div className="flex flex-col">
                        <span className={cn("text-sm font-semibold", theme.headingTextClass)}>{currentPlayer?.name}</span>
                        <span className={cn("text-xs", theme.textAccent)}>Score: {currentPlayer?.score ?? 0}</span>
                      </div>
                    </div>
                    <Button
                      variant="outline"
                      className={cn("rounded-full border", theme.borderClass, theme.textAccent)}
                      onClick={() => setShowPlayers((prev) => !prev)}
                    >
                      <Users className="mr-2 h-4 w-4" /> {showPlayers ? "Hide" : "Show"} collaborators
                    </Button>
                  </div>
                )}

                {showPlayers && (
                  <div className="grid gap-2 rounded-xl border border-slate-200 bg-white/60 p-4 sm:grid-cols-2">
                    {staticOnlinePlayers.map((player) => (
                      <div key={player.name} className="flex items-center gap-3">
                        <div className={cn("h-8 w-8 rounded-full", theme.avatarBgClass)} />
                        <div className="flex flex-col">
                          <span className="text-sm font-medium text-slate-800">{player.name}</span>
                          <span className="text-xs text-slate-500">{player.role}</span>
                        </div>
                      </div>
                    ))}
                  </div>
                )}

                <div className="max-h-[420px] overflow-y-auto rounded-xl border border-slate-200 bg-white/70 p-4">
                  <div className="flex flex-col gap-4">
                    {sentences.map((sentence) => (
                      <div key={sentence.id} className="rounded-xl border border-slate-200 bg-white/80 p-3 shadow-sm">
                        <div className="flex items-center justify-between text-xs text-slate-500">
                          <span className="font-semibold text-slate-700">{sentence.author}</span>
                          <span>{sentence.timestamp}</span>
                        </div>
                        <p className="mt-2 text-sm text-slate-700">{sentence.text}</p>
                      </div>
                    ))}
                    <div ref={messagesEndRef} />
                  </div>
                </div>

                {apiError && <p className="text-sm text-red-500">{apiError}</p>}

                <div className="flex flex-col gap-3">
                  <Textarea
                    placeholder={suggestedLine ?? "Write the next sentence..."}
                    className={cn("min-h-[120px]", theme.borderClass, theme.focusRingClass)}
                    value={newSentence}
                    onChange={(event) => setNewSentence(event.target.value)}
                    disabled={!hasJoined}
                  />
                  <div className="flex flex-wrap items-center justify-between gap-3">
                    <div className="text-xs text-slate-500">
                      {suggestedLine ? "Suggestion inserted, feel free to edit before sending." : "Need inspiration? Try Suggest line."}
                    </div>
                    <div className="flex gap-2">
                      <Button
                        className={cn("rounded-full", theme.buttonClass)}
                        onClick={handleSubmitSentence}
                        disabled={!hasJoined || isSubmitting}
                      >
                        <Send className="mr-2 h-4 w-4" /> {isSubmitting ? "Saving..." : "Publish line"}
                      </Button>
                    </div>
                  </div>
                </div>

                {canonicalStory && (
                  <div className={cn("rounded-xl border p-4", theme.borderClass, theme.softBackgroundClass)}>
                    <h3 className={cn("text-lg font-semibold", theme.headingTextClass)}>Canonical Story Draft</h3>
                    <p className={cn("mt-2 text-sm whitespace-pre-line", theme.textAccent)}>
                      {canonicalStory.full_text ?? "Canonical text available."}
                    </p>
                    <div className="mt-4 flex gap-3">
                      <Button className={cn("rounded-full", theme.buttonClass)} onClick={handleDownloadPdf} disabled={isDownloading}>
                        <Download className="mr-2 h-4 w-4" /> {isDownloading ? "Preparing..." : "Download PDF"}
                      </Button>
                    </div>
                    {downloadError && <p className="mt-2 text-xs text-red-500">{downloadError}</p>}
                  </div>
                )}

                {finishError && <p className="text-sm text-red-500">{finishError}</p>}
              </CardContent>
            </Card>
          </section>
        </div>
      </div>
    </main>
  )
}

import { createServerComponentClient } from "@supabase/auth-helpers-nextjs"
import { cookies } from "next/headers"
import { redirect } from "next/navigation"
import Link from "next/link"
import { formatDistanceToNow } from "date-fns"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"

type StoryStatus = "draft" | "recruiting" | "in_progress" | "completed" | "archived"

type StoryProject = {
	id: string
	title: string
	summary: string | null
	status: StoryStatus
	slots_total: number | null
	slots_taken: number | null
	tags: string[] | null
	is_chain_backed: boolean | null
	chain_reference: string | null
	updated_at: string | null
	template?: {
		id: string | null
		title: string | null
		genre: string | null
	} | null
}

type StoryTemplate = {
	id: string
	title: string
	description: string | null
	genre: string | null
	difficulty: string | null
	estimated_length: string | null
	created_at: string | null
}

const STATUS_META: Record<StoryStatus, { label: string; tone: "default" | "secondary" | "outline" | "destructive" }> = {
	draft: { label: "Draft", tone: "outline" },
	recruiting: { label: "Recruiting", tone: "secondary" },
	in_progress: { label: "In Progress", tone: "default" },
	completed: { label: "Completed", tone: "outline" },
	archived: { label: "Archived", tone: "destructive" },
}

export const dynamic = "force-dynamic"

export default async function DashboardPage() {
	const cookieStore = await cookies()
		const supabase = createServerComponentClient({
			cookies: () => cookieStore as unknown as ReturnType<typeof cookies>,
		})
	const {
		data: { session },
	} = await supabase.auth.getSession()

	if (!session) {
		redirect("/auth?redirect=/dashboard")
	}

		const [
			{ data: storyData, error: storyError },
			{ data: templateData, error: templateError },
		] = await Promise.all([
		supabase
			.from("story_projects")
			.select(
				"id,title,summary,status,slots_total,slots_taken,tags,is_chain_backed,chain_reference,updated_at,template:story_templates(id,title,genre)"
			)
			.order("updated_at", { ascending: false }),
		supabase
			.from("story_templates")
			.select("id,title,description,genre,difficulty,estimated_length,created_at")
			.order("created_at", { ascending: false }),
	])

		if (storyError) {
			console.error("Error loading story projects", storyError)
		}
		if (templateError) {
			console.error("Error loading story templates", templateError)
		}

			const stories: StoryProject[] = (storyData ?? []).map((story) => ({
				...story,
				template: Array.isArray(story.template) ? story.template[0] ?? null : story.template,
			}))
		const templates: StoryTemplate[] = templateData ?? []

	const recruitingStories = stories.filter((story) => story.status === "recruiting" && !story.is_chain_backed)
	const ongoingStories = stories.filter((story) => story.status === "in_progress")
		const chainBackedStories = stories.filter((story) => Boolean(story.is_chain_backed))

	const totalSlots = stories.reduce((acc, story) => acc + (story.slots_total ?? 0), 0)
	const filledSlots = stories.reduce((acc, story) => acc + (story.slots_taken ?? 0), 0)
	const openSlots = Math.max(totalSlots - filledSlots, 0)

	const displayName = session.user.email?.split("@")[0] ?? "Storyteller"

	return (
		<main className="min-h-screen bg-gradient-to-br from-emerald-50 via-white to-cyan-100">
			<div className="mx-auto flex max-w-6xl flex-col gap-12 px-6 py-24">
				<section className="flex flex-col gap-4 rounded-3xl bg-white/70 p-8 shadow-xl ring-1 ring-emerald-200/60 backdrop-blur-lg">
					<div>
						<p className="text-sm font-medium text-emerald-600">Command Center</p>
						<h1 className="mt-2 text-4xl font-semibold text-slate-900">Welcome back, {displayName}.</h1>
						<p className="mt-2 max-w-2xl text-base text-slate-600">
							Pick up an ongoing story, join a new writers' room, or launch a fresh narrative using one of our curated
							templates. Supabase keeps everything in sync today; once the chain goes live we will seamlessly read from
							the ledger.
						</p>
					</div>
					<div className="grid gap-4 sm:grid-cols-3">
						<Card className="bg-gradient-to-br from-emerald-500/10 to-emerald-200/40">
							<CardHeader>
								<CardTitle className="text-sm font-medium text-emerald-800">Active Collaborations</CardTitle>
								<CardDescription className="text-3xl font-semibold text-emerald-900">
									{ongoingStories.length}
								</CardDescription>
							</CardHeader>
						</Card>
						<Card className="bg-gradient-to-br from-cyan-500/10 to-cyan-200/40">
							<CardHeader>
								<CardTitle className="text-sm font-medium text-cyan-800">Open Seats</CardTitle>
								<CardDescription className="text-3xl font-semibold text-cyan-900">{openSlots}</CardDescription>
							</CardHeader>
						</Card>
						<Card className="bg-gradient-to-br from-purple-500/10 to-purple-200/40">
							<CardHeader>
								<CardTitle className="text-sm font-medium text-purple-800">Chain Anchored Stories</CardTitle>
								<CardDescription className="text-3xl font-semibold text-purple-900">{chainBackedStories.length}</CardDescription>
							</CardHeader>
						</Card>
					</div>
				</section>

				<section className="flex flex-col gap-6">
					<div className="flex flex-wrap items-center justify-between gap-4">
						<div>
							<h2 className="text-2xl font-semibold text-slate-900">Active Adventures</h2>
							<p className="text-sm text-slate-600">Stories currently unfolding. Jump in and continue the tale.</p>
						</div>
						<div className="flex gap-3">
							<Button asChild variant="outline" className="rounded-full border-emerald-300 text-emerald-700">
								<Link href="/story-studio">Open Story Studio</Link>
							</Button>
							<Button asChild variant="ghost" className="rounded-full text-purple-700">
								<Link href="/browse">Browse NFTs</Link>
							</Button>
						</div>
					</div>
					<div className="grid gap-6 sm:grid-cols-2">
						{ongoingStories.length === 0 && (
							<Card className="sm:col-span-2 items-center justify-center text-center">
								<CardHeader>
									<CardTitle>No live stories yet</CardTitle>
									<CardDescription>
										When a collaboration switches to <strong>In Progress</strong>, it will appear here.
									</CardDescription>
								</CardHeader>
							</Card>
						)}
						{ongoingStories.map((story) => {
							const meta = STATUS_META[story.status]
							const lastUpdate = story.updated_at
								? formatDistanceToNow(new Date(story.updated_at), { addSuffix: true })
								: "Awaiting kickoff"
							const slotsTotal = story.slots_total ?? 0
							const slotsTaken = story.slots_taken ?? 0
							const slotsRemaining = Math.max(slotsTotal - slotsTaken, 0)
										const isOnChain = Boolean(story.is_chain_backed)

							return (
								<Card key={story.id} className="relative overflow-hidden border-emerald-100/70 bg-white/80">
									<div className="absolute inset-0 bg-gradient-to-br from-emerald-200/20 via-transparent to-cyan-200/30" />
									<CardHeader className="relative">
										<div className="flex items-center justify-between">
											<Badge variant={meta.tone}>{meta.label}</Badge>
											<span className="text-xs font-medium text-slate-500">Updated {lastUpdate}</span>
										</div>
										<CardTitle className="text-xl text-slate-900">{story.title}</CardTitle>
										<CardDescription className="text-slate-600">{story.summary ?? "No summary provided yet."}</CardDescription>
									</CardHeader>
									<CardContent className="relative flex flex-col gap-4">
										<div className="flex items-center gap-3 text-sm text-slate-600">
																<Badge variant="outline">Slots {slotsTaken} / {slotsTotal}</Badge>
																<Badge variant={isOnChain ? "secondary" : "outline"}>
																	{isOnChain ? "On-chain" : "Supabase"}
											</Badge>
											{story.template?.title && <Badge variant="outline">Template: {story.template.title}</Badge>}
										</div>
										{story.tags?.length ? (
											<div className="flex flex-wrap gap-2">
												{story.tags.map((tag) => (
													<Badge key={`${story.id}-${tag}`} variant="outline">
														#{tag}
													</Badge>
												))}
											</div>
										) : null}
									</CardContent>
									<CardFooter className="relative flex items-center justify-between">
										<div className="text-xs text-slate-500">
																{isOnChain ? (
												<span>Anchored on-chain ({story.chain_reference ?? "pending hash"})</span>
											) : (
												<span>{slotsRemaining} seats open for collaborators</span>
											)}
										</div>
										<Button
											asChild
																variant={isOnChain ? "outline" : "default"}
											className="rounded-full"
																disabled={isOnChain}
										>
											<Link href={`/story-studio?id=${story.id}`}>
																	{isOnChain ? "View ledger record" : "Enter writers' room"}
											</Link>
										</Button>
									</CardFooter>
								</Card>
							)
						})}
					</div>
				</section>

				<section className="flex flex-col gap-6">
					<div className="flex flex-wrap items-center justify-between gap-4">
						<div>
							<h2 className="text-2xl font-semibold text-slate-900">Open Collaborations</h2>
							<p className="text-sm text-slate-600">Browse rooms recruiting new writers and claim your spot.</p>
						</div>
					</div>
					<div className="grid gap-6 md:grid-cols-2">
						{recruitingStories.length === 0 && (
							<Card className="md:col-span-2 text-center">
								<CardHeader>
									<CardTitle>No recruiting stories right now</CardTitle>
									<CardDescription>Use a template below to launch a new collaboration.</CardDescription>
								</CardHeader>
							</Card>
						)}
						{recruitingStories.map((story) => {
							const slotsTotal = story.slots_total ?? 0
							const slotsTaken = story.slots_taken ?? 0
							const slotsRemaining = Math.max(slotsTotal - slotsTaken, 0)

							return (
								<Card key={story.id} className="border-cyan-100/70 bg-white/80">
									<CardHeader>
										<div className="flex items-center justify-between">
											<Badge variant="secondary">Recruiting</Badge>
											<span className="text-xs text-slate-500">
												{slotsRemaining} / {slotsTotal} seats open
											</span>
										</div>
										<CardTitle className="text-lg text-slate-900">{story.title}</CardTitle>
										<CardDescription className="text-slate-600">{story.summary ?? "Awaiting pitch."}</CardDescription>
									</CardHeader>
									<CardFooter className="flex items-center justify-between">
										<div className="flex flex-wrap gap-2 text-xs text-slate-500">
											{story.template?.title && <Badge variant="outline">Template: {story.template.title}</Badge>}
											<Badge variant="outline">Supabase</Badge>
										</div>
										<Button asChild className="rounded-full">
											<Link href={`/story-studio?id=${story.id}`}>Join story</Link>
										</Button>
									</CardFooter>
								</Card>
							)
						})}
					</div>
				</section>

				<section className="flex flex-col gap-6">
					<div className="flex flex-wrap items-center justify-between gap-4">
						<div>
							<h2 className="text-2xl font-semibold text-slate-900">Story Templates</h2>
							<p className="text-sm text-slate-600">Kick-start a brand new project using a curated blueprint.</p>
						</div>
						<Button asChild variant="secondary" className="rounded-full">
							<Link href="/story">Create from template</Link>
						</Button>
					</div>
					<div className="grid gap-6 md:grid-cols-3">
						{templates.length === 0 && (
							<Card className="md:col-span-3 text-center">
								<CardHeader>
									<CardTitle>No templates yet</CardTitle>
									<CardDescription>Head to Supabase Studio and seed the `story_templates` table.</CardDescription>
								</CardHeader>
							</Card>
						)}
						{templates.map((template) => (
							<Card key={template.id} className="border-purple-100/70 bg-white/80">
								<CardHeader>
									<div className="flex items-center justify-between">
										<Badge variant="outline">{template.genre ?? "Genre TBD"}</Badge>
										<span className="text-xs text-slate-500">
											{template.created_at
												? formatDistanceToNow(new Date(template.created_at), { addSuffix: true })
												: "Recently added"}
										</span>
									</div>
									<CardTitle className="text-lg text-slate-900">{template.title}</CardTitle>
									<CardDescription className="text-slate-600">
										{template.description ?? "Outline coming soon."}
									</CardDescription>
								</CardHeader>
								<CardContent className="flex flex-col gap-2 text-sm text-slate-600">
									<span>Difficulty: {template.difficulty ?? "Flexible"}</span>
									<span>Estimated length: {template.estimated_length ?? "Custom"}</span>
								</CardContent>
								<CardFooter>
									<Button asChild variant="outline" className="rounded-full">
										<Link href={`/story?template=${template.id}`}>Preview</Link>
									</Button>
								</CardFooter>
							</Card>
						))}
					</div>
				</section>

				
			</div>
		</main>
	)
}

import AuthForm from "@/components/auth-form"
import AnimatedBackground from "@/components/animated-background"

export default function AuthPage() {
	return (
		<div className="relative min-h-screen flex items-center justify-center bg-gradient-to-br from-emerald-50 via-emerald-100 to-teal-200 overflow-hidden">
			<AnimatedBackground activeIndex={0} colorGradient="from-emerald-400 via-emerald-500 to-teal-600" />
			<div className="relative z-10 w-full max-w-4xl flex flex-col items-center px-4 py-16 gap-6">
				<h1
					className="text-5xl font-bold tracking-tight drop-shadow-lg"
					style={{
						fontFamily: "var(--font-titan-one), sans-serif",
						WebkitTextStroke: "3px black",
						color: "transparent",
					}}
				>
					Kahani
				</h1>
				<p className="text-lg sm:text-xl text-foreground/80 text-center max-w-xl">
					Enter a world of stories crafted by dreamers and explorers. Sign in to continue your adventure or join Kahani
					today.
				</p>
				<AuthForm />
			</div>
		</div>
	)
}

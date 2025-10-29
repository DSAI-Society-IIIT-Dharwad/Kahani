"use client"

import { ChangeEvent, FormEvent, useState } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import supabase from "@/lib/supabaseClient"
import styles from "./auth-form.module.css"

export default function AuthForm() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const [status, setStatus] = useState<{
    message: string
    type: "error" | "success"
    form: "signin" | "signup"
  } | null>(null)
  const [isSubmitting, setIsSubmitting] = useState<"signin" | "signup" | null>(null)
  const [isSignUp, setIsSignUp] = useState(false)

  const handleToggle = (event: ChangeEvent<HTMLInputElement>) => {
    setIsSignUp(event.target.checked)
    setStatus(null)
  }

  const handleSignIn = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    const formData = new FormData(event.currentTarget)
    const email = (formData.get("email") || "").toString().trim()
    const password = (formData.get("password") || "").toString()

    if (!email || !password) {
      setStatus({ message: "Please provide both email and password.", type: "error", form: "signin" })
      return
    }

    setIsSubmitting("signin")
    setStatus(null)

    const { error } = await supabase.auth.signInWithPassword({ email, password })

    setIsSubmitting(null)

    if (error) {
      setStatus({ message: error.message, type: "error", form: "signin" })
      return
    }

    setStatus({
      message: "Signed in! Redirecting you to your dashboard...",
      type: "success",
      form: "signin",
    })

    const redirectTo = searchParams.get("redirect") || "/dashboard"
    setTimeout(() => router.push(redirectTo), 800)
  }

  const handleSignUp = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    const formData = new FormData(event.currentTarget)
    const name = (formData.get("name") || "").toString().trim()
    const email = (formData.get("email") || "").toString().trim()
    const password = (formData.get("password") || "").toString()

    if (!name || !email || !password) {
      setStatus({ message: "Please fill in name, email, and password.", type: "error", form: "signup" })
      return
    }

    setIsSubmitting("signup")
    setStatus(null)

    const { error } = await supabase.auth.signUp({
      email,
      password,
      options: {
        data: {
          full_name: name,
        },
      },
    })

    setIsSubmitting(null)

    if (error) {
      setStatus({ message: error.message, type: "error", form: "signup" })
      return
    }

    setStatus({
      message: "Signup successful! Check your email to confirm your account.",
      type: "success",
      form: "signup",
    })
    event.currentTarget.reset()
  }

  return (
    <div className={styles.pageWrapper}>
      <div className={styles.wrapper}>
        <div className={styles.cardSwitch}>
          <label className={styles.switch}>
            <input
              type="checkbox"
              className={styles.toggle}
              checked={isSignUp}
              onChange={handleToggle}
              aria-label={isSignUp ? "Switch to log in" : "Switch to sign up"}
            />
            <div className={styles.cardFrame}>
              <span className={styles.slider} />
              <div className={styles.flipCardInner}>
              <div className={styles.flipCardFront}>
                <div className={styles.title}>Log in</div>
                <form className={styles.flipCardForm} onSubmit={handleSignIn}>
                  <input
                    className={styles.flipCardInput}
                    name="email"
                    placeholder="Email"
                    type="email"
                    autoComplete="email"
                  />
                  <input
                    className={styles.flipCardInput}
                    name="password"
                    placeholder="Password"
                    type="password"
                    autoComplete="current-password"
                  />
                  <div
                    className={`${styles.status} ${
                      status?.form === "signin" && status.type === "success" ? styles.statusSuccess : ""
                    }`}
                  >
                    {status?.form === "signin" ? status.message : ""}
                  </div>
                  <button className={styles.flipCardButton} type="submit" disabled={isSubmitting !== null}>
                    {isSubmitting === "signin" ? "Please wait..." : "Let's go!"}
                  </button>
                </form>
              </div>
              <div className={styles.flipCardBack}>
                <div className={styles.title}>Sign up</div>
                <form className={styles.flipCardForm} onSubmit={handleSignUp}>
                  <input
                    className={styles.flipCardInput}
                    name="name"
                    placeholder="Name"
                    type="text"
                    autoComplete="name"
                  />
                  <input
                    className={styles.flipCardInput}
                    name="email"
                    placeholder="Email"
                    type="email"
                    autoComplete="email"
                  />
                  <input
                    className={styles.flipCardInput}
                    name="password"
                    placeholder="Password"
                    type="password"
                    autoComplete="new-password"
                  />
                  <div
                    className={`${styles.status} ${
                      status?.form === "signup" && status.type === "success" ? styles.statusSuccess : ""
                    }`}
                  >
                    {status?.form === "signup" ? status.message : ""}
                  </div>
                  <button className={styles.flipCardButton} type="submit" disabled={isSubmitting !== null}>
                    {isSubmitting === "signup" ? "Please wait..." : "Confirm!"}
                  </button>
                </form>
              </div>
            </div>
            </div>
          </label>
        </div>
      </div>
    </div>
  )
}

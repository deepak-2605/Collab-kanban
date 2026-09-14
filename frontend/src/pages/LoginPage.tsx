import { useState } from "react"
import { useNavigate } from "react-router-dom"
import { login } from "../api/auth"

export default function LoginPage() {
  const navigate = useNavigate()
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("")

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault() // stop the browser's default full-page form submit
    setError("")
    try {
      const { token } = await login(email, password)
      localStorage.setItem("token", token) // save the JWT
      navigate("/boards") // go to the boards page
    } catch (err) {
      setError("Invalid email or password")
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-100">
      <form onSubmit={handleSubmit} className="w-80 space-y-4 rounded-xl bg-white p-8 shadow-lg">
        <h1 className="text-xl font-bold text-slate-800">Log in</h1>

        {/* TODO: email input — value={email}, onChange={(e) => setEmail(e.target.value)} */}
        {/* TODO: password input — type="password", value={password}, onChange sets password */}
        
        <input
           type="email"
           placeholder="Email"
           value={email}
           onChange={(e) => setEmail(e.target.value)}
           className="w-full rounded-lg border border-slate-300 px-3 py-2"
        />

        <input
           type="password"
           placeholder="Password"
           value={password}
           onChange={(e) => setPassword(e.target.value)}
           className="w-full rounded-lg border border-slate-300 px-3 py-2"
        />

        {error && <p className="text-sm text-red-600">{error}</p>}

        <button type="submit" className="w-full rounded-lg bg-indigo-600 py-2 text-white hover:bg-indigo-700">
          Log in
        </button>
      </form>
    </div>
  )
}
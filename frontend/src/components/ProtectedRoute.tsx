import { Navigate } from "react-router-dom"
import type { ReactNode } from "react"

// Wraps a page that requires login. If there's no token, redirect to /login.
export default function ProtectedRoute({ children }: { children: ReactNode }) {
  const token = localStorage.getItem("token")
  if (!token) {
    return <Navigate to="/login" replace />
  }
  return <>{children}</>
}

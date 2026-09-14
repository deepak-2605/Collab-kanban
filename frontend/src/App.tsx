import { Routes, Route, Navigate } from "react-router-dom"
import LoginPage from "./pages/LoginPage"
import BoardsPage from "./pages/BoardsPage"
import BoardPage from "./pages/BoardPage"
import ProtectedRoute from "./components/ProtectedRoute"

function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route
        path="/boards"
        element={
          <ProtectedRoute>
            <BoardsPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="/boards/:boardId"
        element={
          <ProtectedRoute>
            <BoardPage />
          </ProtectedRoute>
        }
      />
      {/* default: send unknown paths to /boards (which itself redirects to /login if not authed) */}
      <Route path="*" element={<Navigate to="/boards" replace />} />
    </Routes>
  )
}

export default App

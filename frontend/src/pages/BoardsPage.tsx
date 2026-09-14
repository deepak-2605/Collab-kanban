import { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import { listBoards, createBoard } from "../api/boards"
import type { Board } from "../types"

export default function BoardsPage() {
  const navigate = useNavigate()
  const [boards, setBoards] = useState<Board[]>([])
  const [newName, setNewName] = useState("")
  const [loading, setLoading] = useState(true)

  // fetch the user's boards once, on mount
  useEffect(() => {
    listBoards()
      .then(setBoards)
      .catch((e) => console.error("failed to load boards:", e))
      .finally(() => setLoading(false))
  }, [])

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    if (!newName.trim()) return
    const board = await createBoard(newName)
    setBoards([...boards, board]) // add the new board immutably
    setNewName("")
  }

  function handleLogout() {
    localStorage.removeItem("token")
    navigate("/login")
  }

  // 👇 YOU WRITE the return / JSX:
  return (
    <div className="min-h-screen bg-slate-100 p-8">
      <div className="mx-auto max-w-3xl">
        {/* header */}
        <div className="mb-6 flex items-center justify-between">
          <h1 className="text-2xl font-bold text-slate-800">Your Boards</h1>
          <button
            onClick={handleLogout}
            className="rounded-lg border border-slate-300 bg-white px-3 py-1.5 text-sm text-slate-600 hover:bg-slate-50"
          >
            Log out
          </button>
        </div>

        {/* create form */}
        <form onSubmit={handleCreate} className="mb-6 flex gap-2">
          <input
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
            placeholder="New board name"
            className="flex-1 rounded-lg border border-slate-300 px-3 py-2"
          />
          <button
            type="submit"
            className="rounded-lg bg-indigo-600 px-4 py-2 text-white hover:bg-indigo-700"
          >
            Create
          </button>
        </form>

        {/* list */}
        {loading ? (
          <p className="text-slate-500">Loading…</p>
        ) : boards.length === 0 ? (
          <p className="text-slate-500">No boards yet — create your first one above.</p>
        ) : (
          <div className="grid gap-3 sm:grid-cols-2">
            {boards.map((b) => (
              <div
                key={b.id}
                className="cursor-pointer rounded-lg bg-white p-4 shadow hover:shadow-md"
                onClick={() => navigate(`/boards/${b.id}`)}
              >
                <h2 className="font-semibold text-slate-800">{b.name}</h2>
                <p className="mt-1 text-xs text-slate-400">
                  {b.members.length} member{b.members.length === 1 ? "" : "s"}
                </p>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )

}
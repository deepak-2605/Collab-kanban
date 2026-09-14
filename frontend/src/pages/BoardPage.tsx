import { useEffect, useState } from "react"
import { useParams, useNavigate } from "react-router-dom"
import { listColumns, createColumn } from "../api/columns"
import { listCards, createCard } from "../api/cards"
import type { Column, Card } from "../types"

export default function BoardPage() {
  const { boardId } = useParams()          // reads ":boardId" from the URL
  const navigate = useNavigate()

  const [columns, setColumns] = useState<Column[]>([])
  const [cardsByColumn, setCardsByColumn] = useState<Record<string, Card[]>>({})
  const [loading, setLoading] = useState(true)
  const [newColumnName, setNewColumnName] = useState("")
  const [cardDrafts, setCardDrafts] = useState<Record<string, string>>({})

  // fetch columns, then each column's cards, on mount / when boardId changes
  useEffect(() => {
    const id = boardId
    if (!id) return
    async function load() {
      const cols = await listColumns(id)
      setColumns(cols)
      // fetch every column's cards in parallel, build a { columnId: cards[] } map
      const entries = await Promise.all(
        cols.map(async (c) => [c.id, await listCards(c.id)] as const)
      )
      setCardsByColumn(Object.fromEntries(entries))
    }
    load()
      .catch((e) => console.error("failed to load board:", e))
      .finally(() => setLoading(false))
  }, [boardId])

  async function handleAddColumn(e: React.FormEvent) {
    e.preventDefault()
    if (!boardId || !newColumnName.trim()) return
    const col = await createColumn(boardId, newColumnName)
    setColumns([...columns, col])
    setCardsByColumn({ ...cardsByColumn, [col.id]: [] }) // give the new column an empty card list
    setNewColumnName("")
  }

  async function handleAddCard(columnId: string) {
    const title = (cardDrafts[columnId] || "").trim()
    if (!title) return
    const card = await createCard(columnId, title)
    setCardsByColumn({
      ...cardsByColumn,
      [columnId]: [...(cardsByColumn[columnId] || []), card],
    })
    setCardDrafts({ ...cardDrafts, [columnId]: "" })
  }

  if (loading) return <p className="p-8 text-slate-500">Loading board…</p>

  return (
    <div className="min-h-screen bg-slate-100 p-8">
      {/* header */}
      <button
        onClick={() => navigate("/boards")}
        className="mb-4 text-sm text-indigo-600 hover:underline"
      >
        ← Back to boards
      </button>

      {/* add-column form */}
      <form onSubmit={handleAddColumn} className="mb-6 flex gap-2">
        <input
          value={newColumnName}
          onChange={(e) => setNewColumnName(e.target.value)}
          placeholder="New column"
          className="rounded-lg border border-slate-300 px-3 py-2"
        />
        <button type="submit" className="rounded-lg bg-indigo-600 px-4 py-2 text-white hover:bg-indigo-700">
          Add column
        </button>
      </form>

      {/* the board: columns side by side */}
      <div className="flex gap-4 overflow-x-auto pb-4">
        {columns.map((col) => (
          <div key={col.id} className="w-72 flex-shrink-0 rounded-lg bg-slate-200 p-3">
            <h2 className="mb-3 font-semibold text-slate-700">{col.name}</h2>

            {/* cards */}
            <div className="space-y-2">
              {(cardsByColumn[col.id] || []).map((card) => (
                <div key={card.id} className="rounded-lg bg-white p-3 shadow-sm">
                  <p className="text-sm text-slate-800">{card.title}</p>
                </div>
              ))}
            </div>

            {/* add-card input */}
            <input
              value={cardDrafts[col.id] || ""}
              onChange={(e) => setCardDrafts({ ...cardDrafts, [col.id]: e.target.value })}
              onKeyDown={(e) => {
                if (e.key === "Enter") handleAddCard(col.id)
              }}
              placeholder="+ Add a card (Enter)"
              className="mt-3 w-full rounded-lg border border-slate-300 bg-white px-2 py-1.5 text-sm"
            />
          </div>
        ))}
      </div>
    </div>
  )
}
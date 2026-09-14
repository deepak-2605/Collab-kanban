import { useEffect, useState } from "react"
import { useParams, useNavigate } from "react-router-dom"
import { DndContext, useDraggable, useDroppable } from "@dnd-kit/core"
import type { DragEndEvent } from "@dnd-kit/core"
import { listColumns, createColumn } from "../api/columns"
import { listCards, createCard, moveCard } from "../api/cards"
import type { Column, Card } from "../types"

// ---- a single draggable card ----
function DraggableCard({ card }: { card: Card }) {
  const { attributes, listeners, setNodeRef, transform, isDragging } = useDraggable({ id: card.id })
  const style = transform
    ? { transform: `translate(${transform.x}px, ${transform.y}px)` }
    : undefined
  return (
    <div
      ref={setNodeRef}
      style={style}
      {...listeners}
      {...attributes}
      className={`cursor-grab rounded-lg bg-white p-3 shadow-sm ${isDragging ? "opacity-50" : ""}`}
    >
      <p className="text-sm text-slate-800">{card.title}</p>
    </div>
  )
}

// ---- a droppable column (highlights when a card hovers over it) ----
function DroppableColumn({ column, children }: { column: Column; children: React.ReactNode }) {
  const { setNodeRef, isOver } = useDroppable({ id: column.id })
  return (
    <div
      ref={setNodeRef}
      className={`w-72 flex-shrink-0 rounded-lg p-3 ${isOver ? "bg-slate-300" : "bg-slate-200"}`}
    >
      <h2 className="mb-3 font-semibold text-slate-700">{column.name}</h2>
      {children}
    </div>
  )
}

export default function BoardPage() {
  const { boardId } = useParams()
  const navigate = useNavigate()

  const [columns, setColumns] = useState<Column[]>([])
  const [cardsByColumn, setCardsByColumn] = useState<Record<string, Card[]>>({})
  const [loading, setLoading] = useState(true)
  const [newColumnName, setNewColumnName] = useState("")
  const [cardDrafts, setCardDrafts] = useState<Record<string, string>>({})

  useEffect(() => {
    const id = boardId
    if (!id) return
    async function load() {
      const cols = await listColumns(id)
      setColumns(cols)
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
    setCardsByColumn({ ...cardsByColumn, [col.id]: [] })
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

  // ---- the drag-and-drop handler ----
  async function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event
    if (!over) return // dropped outside any column

    const cardId = active.id as string
    const targetColumnId = over.id as string

    // which column is the card currently in?
    const sourceColumnId = Object.keys(cardsByColumn).find((colId) =>
      cardsByColumn[colId].some((c) => c.id === cardId)
    )
    if (!sourceColumnId || sourceColumnId === targetColumnId) return // no move

    const card = cardsByColumn[sourceColumnId].find((c) => c.id === cardId)!
    const newPosition = (cardsByColumn[targetColumnId] || []).length
    
    const previous = cardsByColumn
    // 1. optimistic UI update: remove from source, append to target
    setCardsByColumn((prev) => ({
      ...prev,
      [sourceColumnId]: prev[sourceColumnId].filter((c) => c.id !== cardId),
      [targetColumnId]: [...(prev[targetColumnId] || []), { ...card, columnId: targetColumnId, position: newPosition }],
    }))

    // 2. persist to the backend
    try {
      await moveCard(cardId, targetColumnId, newPosition)
    } catch (e) {
      console.error("move failed:", e) // (a real app would revert the optimistic update here)
      setCardsByColumn(previous)
    }
  }

  if (loading) return <p className="p-8 text-slate-500">Loading board…</p>

  return (
    <div className="min-h-screen bg-slate-100 p-8">
      <button onClick={() => navigate("/boards")} className="mb-4 text-sm text-indigo-600 hover:underline">
        ← Back to boards
      </button>

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

      <DndContext onDragEnd={handleDragEnd}>
        <div className="flex gap-4 overflow-x-auto pb-4">
          {columns.map((col) => (
            <DroppableColumn key={col.id} column={col}>
              <div className="space-y-2">
                {(cardsByColumn[col.id] || []).map((card) => (
                  <DraggableCard key={card.id} card={card} />
                ))}
              </div>
              <input
                value={cardDrafts[col.id] || ""}
                onChange={(e) => setCardDrafts({ ...cardDrafts, [col.id]: e.target.value })}
                onKeyDown={(e) => {
                  if (e.key === "Enter") handleAddCard(col.id)
                }}
                placeholder="+ Add a card (Enter)"
                className="mt-3 w-full rounded-lg border border-slate-300 bg-white px-2 py-1.5 text-sm"
              />
            </DroppableColumn>
          ))}
        </div>
      </DndContext>
    </div>
  )
}
// These mirror the JSON your Go API returns (json tags on the structs).

export interface User {
  id: string
  name: string
  email: string
  createdAt: string
}

export interface Board {
  id: string
  name: string
  ownerId: string
  members: string[]
  createdAt: string
  updatedAt: string
}

export interface Column {
  id: string
  boardId: string
  name: string
  position: number
  createdAt: string
  updatedAt: string
}

export interface Card {
  id: string
  boardId: string
  columnId: string
  title: string
  description: string
  position: number
  assigneeId?: string // optional — unassigned cards omit it
  createdAt: string
  updatedAt: string
}

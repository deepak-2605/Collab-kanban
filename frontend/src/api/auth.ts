import { apiFetch } from "./client";

export function login(email :string,password :string){
    return apiFetch<{token: string}>("/login",{
        method: "POST",
        body: JSON.stringify({ email, password }),
    })
}
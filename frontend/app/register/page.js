"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";

export default function RegisterPage() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const router = useRouter();

  const handleRegister = () => {
    if (!username || !password) return alert("Fill all fields");

    const users = JSON.parse(localStorage.getItem("users")) || [];
    const exists = users.find(u => u.username === username);
    if (exists) return alert("User already exists");

    users.push({ username, password });
    localStorage.setItem("users", JSON.stringify(users));
    alert("Registered successfully!");
    router.push("/login");
  };

  return (
    <div>
      <h2>Register</h2>
      <input
        placeholder="Username"
        value={username}
        onChange={e => setUsername(e.target.value)}
      /><br/>
      <input
        type="password"
        placeholder="Password"
        value={password}
        onChange={e => setPassword(e.target.value)}
      /><br/>
      <button onClick={handleRegister}>Register</button>
    </div>
  );
}

"use client";

import { useState } from "react";

export default function Home() {
  const [error, setError] = useState("");
  const [result, setResult] = useState("");

  const API = "http://localhost:8080";

  // 🔍 SEARCH FUNCTION
  const handleSearch = async () => {
    const res = await fetch(`${API}/search?error=${error}`);
    const data = await res.json();
    setResult(JSON.stringify(data));
  };

  // 💾 SAVE FUNCTION
  const handleSave = async () => {
    await fetch(`${API}/issue`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        error,
        cause: "manual entry",
        fix: "manual fix",
      }),
    });

    alert("Saved successfully");
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-800 via-purple-900 to-pink-700 text-white flex flex-col items-center justify-center p-6">

      <h1 className="text-4xl font-bold mb-6">
        DevOps Memory Assistant 🚀
      </h1>

      <input
        value={error}
        onChange={(e) => setError(e.target.value)}
        placeholder="Enter error (e.g. CrashLoopBackOff)"
        className="w-full max-w-md p-3 rounded-lg text-black mb-4"
      />

      <div className="flex gap-4">
        <button onClick={handleSearch}
          className="bg-pink-500 px-5 py-2 rounded-lg hover:bg-pink-600">
          Search
        </button>

        <button onClick={handleSave}
          className="bg-purple-500 px-5 py-2 rounded-lg hover:bg-purple-600">
          Save
        </button>
      </div>

      {result && (
        <div className="mt-6 bg-white text-black p-4 rounded-lg max-w-md">
          {result}
        </div>
      )}

    </div>
  );
}
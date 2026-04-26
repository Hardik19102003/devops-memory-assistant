"use client";

import { useState } from "react";

export default function Home() {
  const [error, setError] = useState("");
  const [results, setResults] = useState([]);

  const handleSearch = async () => {
    const res = await fetch(
      `http://localhost:8080/search?error=${error}`
    );

    const data = await res.json();
    setResults(data);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-900 via-purple-800 to-pink-700 text-white flex flex-col items-center p-6">

      <h1 className="text-4xl font-bold mb-8 text-center">
        DevOps Memory Assistant 🚀
      </h1>

      <div className="bg-white/10 backdrop-blur-lg p-6 rounded-2xl shadow-lg w-full max-w-md">

        <input
          value={error}
          onChange={(e) => setError(e.target.value)}
          placeholder="Enter error (e.g. CrashLoopBackOff)"
          className="w-full p-3 rounded-lg bg-white/90 text-black placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-pink-400 mb-4"
        />

        <button
          onClick={handleSearch}
          className="w-full bg-pink-500 py-2 rounded-lg hover:bg-pink-600 transition"
        >
          Search
        </button>

      </div>

      {/* Results Section */}
      <div className="mt-8 w-full max-w-md space-y-4">
        {results.map((item: any, index) => (
          <div
            key={index}
            className="bg-white/10 backdrop-blur-lg p-4 rounded-xl"
          >
            <p><strong>Error:</strong> {item.error}</p>
            <p><strong>Cause:</strong> {item.cause}</p>
            <p><strong>Fix:</strong> {item.fix}</p>
          </div>
        ))}
      </div>

    </div>
  );
}
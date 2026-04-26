"use client";

import { useState } from "react";
import { motion } from "framer-motion";

export default function Home() {
  const [error, setError] = useState("");
  const [cause, setCause] = useState("");
  const [fix, setFix] = useState("");
  const [results, setResults] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");

  const handleSearch = async () => {
    setLoading(true);
    setMessage("");

    try {
      const res = await fetch(
        `process.env.NEXT_PUBLIC_API_URL/search?error=${error}`
      );
      const data = await res.json();
      setResults(data);
    } catch {
      setMessage("Something went wrong ❌");
    }

    setLoading(false);
  };

  const handleSave = async () => {
    setLoading(true);
    setMessage("");

    try {
      await fetch("process.env.NEXT_PUBLIC_API_URL/issue", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ error, cause, fix }),
      });

      setMessage("Saved successfully ✅");
      setError("");
      setCause("");
      setFix("");
    } catch {
      setMessage("Failed to save ❌");
    }

    setLoading(false);
  };

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="min-h-screen bg-gradient-to-br from-purple-900 via-purple-800 to-pink-700 text-white flex flex-col items-center p-6"
    >

      <motion.h1
        initial={{ y: -20, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        className="text-4xl font-bold mb-6 text-center"
      >
        DevOps Memory Assistant 🚀
      </motion.h1>

      <motion.div
        initial={{ scale: 0.95, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        className="bg-white/10 backdrop-blur-lg p-6 rounded-2xl shadow-lg w-full max-w-md space-y-4"
      >

        <input
          value={error}
          onChange={(e) => setError(e.target.value)}
          placeholder="Error (e.g. CrashLoopBackOff)"
          className="w-full p-3 rounded-lg bg-white/90 text-black focus:ring-2 focus:ring-pink-400"
        />

        <input
          value={cause}
          onChange={(e) => setCause(e.target.value)}
          placeholder="Cause"
          className="w-full p-3 rounded-lg bg-white/90 text-black"
        />

        <input
          value={fix}
          onChange={(e) => setFix(e.target.value)}
          placeholder="Fix"
          className="w-full p-3 rounded-lg bg-white/90 text-black"
        />

        <div className="flex gap-4">
          <motion.button
            whileTap={{ scale: 0.95 }}
            onClick={handleSearch}
            disabled={loading}
            className="flex-1 bg-pink-500 py-2 rounded-lg hover:bg-pink-600 transition"
          >
            {loading ? "Searching..." : "Search"}
          </motion.button>

          <motion.button
            whileTap={{ scale: 0.95 }}
            onClick={handleSave}
            disabled={loading}
            className="flex-1 bg-purple-500 py-2 rounded-lg hover:bg-purple-600 transition"
          >
            {loading ? "Saving..." : "Save"}
          </motion.button>
        </div>

        {message && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="text-center text-sm bg-black/30 p-2 rounded"
          >
            {message}
          </motion.div>
        )}

      </motion.div>

      {/* Results */}
      <div className="mt-8 w-full max-w-md space-y-4">

        {results.length === 0 && !loading && (
          <p className="text-center text-gray-300">
            No results yet 👀
          </p>
        )}

        {results.map((item, index) => (
          <motion.div
            key={index}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: index * 0.1 }}
            className="bg-white/10 backdrop-blur-lg p-4 rounded-xl border border-white/20 hover:scale-[1.03] transition"
          >
            <p><strong>Error:</strong> {item.error}</p>
            <p><strong>Cause:</strong> {item.cause}</p>
            <p><strong>Fix:</strong> {item.fix}</p>
          </motion.div>
        ))}
      </div>

    </motion.div>
  );
}